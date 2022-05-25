package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const solrConfigsDir = "/Users/d066419/go/src/github.com/abbi-gaurav/go-projects/ultimate-go-programming/json-files"
const suggestComponent = "suggest-component.json"
const suggestRequestHandler = "suggest-request-handler.json"
const commandAddSearchComponent = "add-searchcomponent"
const commandUpdateSearchComponent = "update-searchcomponent"
const commandAddRequestHandler = "add-requesthandler"
const commandUpdateRequestHandler = "update-requesthandler"

func main() {
	addSearchComponent := make(map[string]map[string]interface{})
	searchComponentConfig := readConfigJson(solrConfigsDir + "/" + suggestComponent)
	addSearchComponent["add-searchcomponent"] = searchComponentConfig
	bs2, err := json.Marshal(addSearchComponent)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(bs2))
}

func readConfigJson(path string) map[string]interface{} {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(bs, &m)
	if err != nil {
		panic(err)
	}
	return m
}
