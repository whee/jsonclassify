// Copyright (c) 2014 Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC
// license which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Attributes map[string][]float64
type Weights map[string]int
type Data map[string]interface{}

type Classifier struct {
	Name       string
	ScoreName  string
	Categories map[string]Attributes
	Weights    Weights
}

type Category struct {
	Name  string
	Score int
}

func NewClassifier(file string) (*Classifier, error) {
	configData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c Classifier
	err = json.Unmarshal(configData, &c)
	c.ScoreName = c.Name + " Score"
	return &c, err
}

func (c *Classifier) Classify(d Data) Category {
	var cat Category
	for name, attrs := range c.Categories {
		if s := attrs.Score(d, c.Weights); s > cat.Score {
			cat = Category{name, s}
		}
	}
	return cat
}

func (a Attributes) Score(d Data, w Weights) int {
	var score int
	for attr, i := range a {
		if d[attr] == nil {
			continue
		}
		low, high := i[0], i[1]
		val, ok := d[attr].(float64)
		if !ok {
			var err error
			val, err = strconv.ParseFloat(d[attr].(string), 64)
			if err != nil {
				log.Fatal(err)
			}
		}
		if val >= low && val < high {
			score += w[attr]
		}
	}
	return score
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: jsonclassify [configuration file]\n")
	}
	configFile := os.Args[1]
	c, err := NewClassifier(configFile)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	decoded := make(chan Data, 10)
	scored := make(chan Data, 10)

	go func() {
		defer close(decoded)
		for {
			var jsd Data
			if err := dec.Decode(&jsd); err != nil {
				if err == io.EOF {
					return
				}
				log.Fatal(err)
			}
			decoded <- jsd
		}
	}()

	go func() {
		defer close(scored)
		for d := range decoded {
			category := c.Classify(d)
			d[c.Name] = category.Name
			d[c.ScoreName] = category.Score
			scored <- d
		}
	}()

	for s := range scored {
		if err := enc.Encode(&s); err != nil {
			log.Fatal(err)
		}
	}
}
