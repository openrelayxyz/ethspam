package main

import (
	"fmt"
	"encoding/json"
	"math/rand"
	"time"
	"io/ioutil"
	"os"

)

func jsonUnmarshal(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening %s for json unmarshalling, %v", file, err.Error())
	}
	raw, _ := ioutil.ReadAll(f)
	var values []string
	json.Unmarshal(raw, &values)
	return values, nil
}

var emptySlice []string

var types = map[int]interface{}{0: nil, 1: emptySlice, 2: ""}

func getRandom(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	value := rand.Intn(max-min+1) + min
	return value
}

func fuzzAddress() (string, error) {

	addresses, err := jsonUnmarshal("testdata/address.json")
	if err != nil {
		fmt.Printf("Error ummarshalling %s data, %v", "address", err.Error())
	}

	return addresses[getRandom(0, (len(addresses) -1))], nil

}

func fuzzTopics() ([]string, error) {

	topicZero, err := jsonUnmarshal("testdata/topic0.json")
	if err != nil {
		fmt.Printf("Error ummarshalling topic %s data, %v", "zero", err.Error())
	}

	topicOne, err := jsonUnmarshal("testdata/topic1.json")
	if err != nil {
		fmt.Printf("Error ummarshalling topic %s data, %v", "one", err.Error())
	}

	topicTwo, err := jsonUnmarshal("testdata/topic2.json")
	if err != nil {
		fmt.Printf("Error ummarshalling topic %s data, %v", "two", err.Error())
	}

	topicThree, err := jsonUnmarshal("testdata/topic3.json")
	if err != nil {
		fmt.Printf("Error ummarshalling topic %s data, %v", "three", err.Error())
	}

	topics := map[int][]string{0: topicZero, 1: topicOne, 2: topicTwo, 3: topicThree}

	length := getRandom(0, 4)

	tParams := make([]interface{}, length)

	for i, _ := range tParams {
		tParams[i] = types[getRandom(0, 2)]
	}

	if len(tParams) == 0 || len(tParams) == 1 && tParams[0] == nil {
		// p, err := json.Marshal("")
		// if err != nil {
		// 	fmt.Printf("Error marshalling empty return value from fuzz %v", err.Error())
		// 	return "", err
		// }
		return []string{}, nil
	}

	for i, item := range tParams {

		switch item.(type) {
		case string:
			tParams[i] = topics[i][getRandom(0, (len(topics[i])-1))]
		case []string:
			inner := make([]string, getRandom(1, 3))
			for j, _ := range inner {
				inner[j] = topics[i][getRandom(0, (len(topics[i])-1))]
			}
			tParams[i] = inner
		default:
			tParams[i] = item
		}
	}

	p, err := json.Marshal(tParams)
	if err != nil {
		fmt.Printf("Error marshalling non emptpy return values from fuzz %v", err.Error())
	}

	return []string{string(p)}, nil
}
