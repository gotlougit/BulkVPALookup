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

//perform the actual network request with a suitable user agent
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
	req, err := http.NewRequest("POST", baseurl+vpa, reqBody)
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

//returns name if exists, empty string otherwise
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

//bridge between network requests and going through all numbers given to us
func sendToChannel(number string, suffix string, mappings map[string]string) {
	if mappings[number] != "" {
		return
	}
	name := getNameIfExists(number, suffix)
	//name := "Dummy McDumbface" //dummy request response for testing
	if name == "" {
		return
	}
	mappings[number] = name
	fmt.Println(number, ":", name)
}

func performBulkLookup(numbers []string, lookedUpNames map[string]string) {

	var suffices = []string{"paytm", "ybl"}

	for _, suffix := range suffices {
		for _, number := range numbers {
			if len(number) != 10 || lookedUpNames[number] != "" {
				continue
			}
			sendToChannel(number, suffix, lookedUpNames)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

//somewhat redundant function created as a result of bad network conditions halting the program before it could create the VCF
func getBulkLookupResults(filename string) map[string]string {
	rawcontent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	m := make(map[string]string)
	content := string(rawcontent)
	lines := strings.Split(content, "\n")
	for line := range lines {
		pair := strings.Split(lines[line], ":")
		if m[pair[0]] == "" {
			if len(pair) > 1 {
				m[pair[0]] = pair[1]
			}
		}
	}
	return m
}

//help export our results to VCF for easy importing
func writeResultsToVCF(lookedUpNames map[string]string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	for number, name := range lookedUpNames {
		file.Write([]byte("BEGIN:VCARD\nVERSION:3.0\n"))
		file.Write([]byte("FN:" + name + "\n"))
		splitrep := strings.Split(name, " ")
		newrep := ""
		for word := range splitrep {
			if word == 0 {
				newrep = splitrep[word]
			} else {
				newrep = splitrep[word] + ";" + newrep
			}
		}
		newrep += ";;"
		file.Write([]byte("N:" + newrep + "\n"))
		file.Write([]byte("TEL;TYPE=cell:+91 " + number + "\n"))
		file.Write([]byte("END:VCARD\n\n"))
	}

}

func main() {

	argLength := len(os.Args[1:])
	if argLength != 2 {
		log.Fatalln("USAGE: ./main /path/to/list/of/phone/nums.txt /path/to/vcf/file.vcf")
	}

	numbersFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	numbers := strings.Split(string(numbersFile[:]), "\n")
	var lookedUpNames = make(map[string]string)
	performBulkLookup(numbers, lookedUpNames)

	writeResultsToVCF(lookedUpNames, os.Args[2])

}
