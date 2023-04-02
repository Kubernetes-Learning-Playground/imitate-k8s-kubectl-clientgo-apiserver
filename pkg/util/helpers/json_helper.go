package helpers

import (
	"k8s.io/apimachinery/pkg/util/json"
	"log"
)

func ToJson(inputStruct interface{}) []byte {
	res, err := json.Marshal(inputStruct)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
