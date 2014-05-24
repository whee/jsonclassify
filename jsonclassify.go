package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Attributes map[string][]int
type Weights map[string]int

type Classifier struct {
	Name       string
	Categories map[string]Attributes
	Weights    Weights
}

type Category struct {
	Name       string
	Attributes map[string][]int
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

func (c *Classifier) Classify(d map[string]interface{}) string {
	return d["Data"].(string)
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
		var jsd map[string]interface{}
		if err := dec.Decode(&jsd); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if _, ok := jsd["Attributes"]; !ok {
			jsd["Attributes"] = make(map[string]interface{})
		}

		jsd["Attributes"].(map[string]interface{})[c.Name] = c.Classify(jsd)
		if err := enc.Encode(&jsd); err != nil {
			log.Fatal(err)
		}
	}
}
