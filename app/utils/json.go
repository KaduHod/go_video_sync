package utils

import (
	"encoding/json"
	"fmt"
)

func JsonEncode(data interface{}) string {
    jsonBytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Erro converting struct to json", err)
        fmt.Println("Dado",  data)
        panic(err)
    }
    return string(jsonBytes)
}
func JsonDecode[T any](data string, target T) {
	err := json.Unmarshal([]byte(data), target)
	if err != nil {
		fmt.Println("Erro convertendo JSON para struct:", err)
        fmt.Println("Dado", data)
		panic(err)
	}
}
func JasonParse[T any](data string, target T) error {
	err := json.Unmarshal([]byte(data), target)
	if err != nil {
		fmt.Println("Erro convertendo JSON para struct:", err)
        fmt.Println("Dado", data)
        return err
	}
    return nil
}
func JsonStringify(data interface{}) (string, error) {
    jsonBytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Erro converting struct to json", err)
        return "", err
    }
    return string(jsonBytes), nil
}
