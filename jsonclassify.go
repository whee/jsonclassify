package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: jsonclassify [configuration file]\n")
	}

	configFile := os.Args[1]
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c Classifier
	if err := json.Unmarshal(configData, &c); err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
