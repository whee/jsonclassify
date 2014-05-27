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
	Categories map[string]Attributes
	Weights    Weights
}

func NewClassifier(file string) (*Classifier, error) {
	configData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c Classifier
	err = json.Unmarshal(configData, &c)
	return &c, err
}

func (c *Classifier) Classify(d Data) string {
	var high int
	var category string
	for cat, attrs := range c.Categories {
		if s := attrs.Score(d, c.Weights); s > high {
			high = s
			category = cat
		}
	}
	return category
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

	for {
		var jsd Data
		if err := dec.Decode(&jsd); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		jsd[c.Name] = c.Classify(jsd)
		if err := enc.Encode(&jsd); err != nil {
			log.Fatal(err)
		}
	}
}
