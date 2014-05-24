package main

import (
	"encoding/json"
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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: jsonclassify [configuration file]\n")
	}
	configFile := os.Args[1]
	c, err := NewClassifier(configFile)
	if err != nil {
		log.Fatal(err)
	}

}
