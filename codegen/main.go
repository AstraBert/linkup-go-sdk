package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	// Dynamically fetch OpenAPI spec
	resp, err := http.Get("https://api.linkup.so/v1/openapi.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert JSON to YAML
	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		panic(err)
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("openapi.yaml", yamlData, 0644); err != nil {
		panic(err)
	}
}
