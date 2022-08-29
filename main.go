package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//struct for storing API responses
type UpiResponse struct {
	IsUpiRegistered bool   `json: "isUpiRegistered, bool"`
	Name            string `json: "name, string"`
	Message         string `json: "message, string"`
}

func makeAPIRequest(number string, suffix string) string {
	//prepare request
	vpa := number + "@" + suffix

	postBody, _ := json.Marshal(map[string]string{
		"upi": vpa,
	})
	reqBody := bytes.NewBuffer(postBody)

	//make request
	resp, err := http.Post("https://upibankvalidator.com/api/upiValidation?upi="+vpa, "application/json", reqBody)
	if err != nil {
		log.Println("Error occurred!")
		log.Fatalln(err)
	}
	//process response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occurred!")
		log.Fatalln(err)
	}
	sb := string(respBody)
	return sb
}

func getNameIfExists(number string, suffix string) string {
	var processedResp UpiResponse
	rawResp := makeAPIRequest(number, suffix)
	//unmarshal response into struct
	err := json.Unmarshal([]byte(rawResp), &processedResp)
	if err != nil {
		log.Println("Error occurred!")
		if (rawResp == "error code: 1015") {
			log.Fatalln("Too many requests!")
		}
		log.Fatalln(err)
	}
	if !processedResp.IsUpiRegistered {
		return ""
	}
	return processedResp.Name
}

func sendToChannel(ch chan map[string]string, number string, suffix string) {
	name := getNameIfExists(number, suffix)
	if (name == "") {
		return;
	}
	var tempMap = make(map[string]string)
	tempMap[number] = name
	ch <- tempMap
}

func performBulkLookup(numbers []string, lookedUpNames map[string]string) {

	ch := make(chan map[string]string, len(numbers))

	var suffices = []string{"paytm"}
	for _, suffix := range suffices {
		for _, number := range numbers {
			go sendToChannel(ch, number, suffix)
		}
	}

	for i:=0; i < len(numbers); i++ {
		var tempMap = make(map[string]string)
		tempMap = <-ch
		for number, _ := range tempMap {
			lookedUpNames[number] = tempMap[number]
		}
	}

}

func main() {

	argLength := len(os.Args[1:])
	if (argLength != 1) {
		log.Fatalln("USAGE: ./main /path/to/list/of/phone/nums.txt")
	}

	numbersFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	numbers := strings.Split(string(numbersFile[:]), "\n")
	var lookedUpNames = make(map[string]string)
	performBulkLookup(numbers, lookedUpNames)

	for number, name := range lookedUpNames {
		fmt.Println(number, ": ", name)
	}

}
