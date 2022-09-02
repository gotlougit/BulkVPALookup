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
	"time"
)

//struct for storing API responses
type UpiResponse struct {
	IsUpiRegistered bool   `json: "isUpiRegistered, bool"`
	Name            string `json: "name, string"`
	Message         string `json: "message, string"`
}

func makeAPIRequest(number string, suffix string) string {

	useragent := "BulkVPALookup/1.0"
	baseurl := "https://upibankvalidator.com/api/upiValidation?upi="

	//prepare request
	vpa := number + "@" + suffix

	postBody, _ := json.Marshal(map[string]string{
		"upi": vpa,
	})
	reqBody := bytes.NewBuffer(postBody)

	client := &http.Client{}

	//make request
	req, err := http.NewRequest("POST",baseurl+vpa, reqBody)
	if err != nil {
		log.Println("Error occurred!")
		log.Fatalln(err)
	}        
	req.Header.Set("User-Agent", useragent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

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
		if rawResp == "error code: 1015" {
			log.Fatalln("Too many requests!")
		}
		log.Println(rawResp)
		log.Fatalln(err)
	}
	if !processedResp.IsUpiRegistered {
		return ""
	}
	return processedResp.Name
}

func sendToChannel(number string, suffix string, mappings map[string]string) {
	name := getNameIfExists(number, suffix)
	//name := "Dummy McDumbface" //dummy request response
	if name == "" {
		return
	}
	mappings[number] = name
	fmt.Println(number, ":", name)
}

func performBulkLookup(numbers []string, lookedUpNames map[string]string) {

	var mapMutex = make(chan int, 1)

	var suffices = []string{"paytm"}

	for _, suffix := range suffices {
		for _, number := range numbers {
			if len(number) != 10 || lookedUpNames[number] != "" {
				continue
			}
			mapMutex <- 1
			go func() {
				sendToChannel(number, suffix, lookedUpNames)
				time.Sleep(500 * time.Millisecond)
				<-mapMutex
			}()
		}
	}
}

func main() {

	argLength := len(os.Args[1:])
	if argLength != 1 {
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
		fmt.Println(number, ":", name)
	}

}
