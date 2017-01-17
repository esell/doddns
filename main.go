package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type IPRecord struct {
	CurrentIP string `json:"ip"`
}

type DomainRecord struct {
	Id         int    `json:"id"`
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Priority   int    `json:"priority"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
}

type DomainRecords struct {
	ExistingDomainRecords []DomainRecord `json:"domain_records"`
	Links                 struct{}       `json:"links"`
	Meta                  struct {
		Total int64 `json:"total"`
	} `json:"meta"`
}

var (
	apiKey     = flag.String("k", "", "DO API key")
	domainName = flag.String("d", "", "base domain")
	subDomain  = flag.String("s", "", "subdomain to update")
)

func main() {
	flag.Parse()
	// Get current public IP
	if (*subDomain == "") || (*domainName == "") || (*apiKey == "") {
		log.Println("Must provide a domain, subdomain and api key. Exiting...")
		os.Exit(1)
	}
	currentIP, err := httpReq("GET", "https://api.ipify.org?format=json", nil)
	contentsIP, err := ioutil.ReadAll(currentIP.Body)
	if err != nil {
		log.Println("Error reading response body for current IP: ", err)
	}
	var currentIPJSON IPRecord
	err = json.Unmarshal(contentsIP, &currentIPJSON)
	if err != nil {
		log.Println("Error decoding JSON for current IP: ", err)
	}
	log.Println("Current IP is: ", currentIPJSON.CurrentIP)

	recordListResp, err := httpReq("GET", "https://api.digitalocean.com/v2/domains/"+*domainName+"/records", nil)
	if err != nil {
		log.Println("Error getting domain records: ", err)
	}
	defer recordListResp.Body.Close()

	contents, err := ioutil.ReadAll(recordListResp.Body)
	if err != nil {
		log.Println("Error reading response body for domain records: ", err)
	}

	var m DomainRecords
	err = json.Unmarshal(contents, &m)
	if err != nil {
		log.Println("Error decoding JSON for domain records: ", err)
	}

	subExistsResult, recordID := subExists(m.ExistingDomainRecords)
	if subExistsResult {
		// update
		log.Println("subdomain exists, updating record ID: ", strconv.Itoa(recordID))
		jsonStr := []byte(`{"data":"` + currentIPJSON.CurrentIP + `"}`)
		updateResp, err := httpReq("PUT", "https://api.digitalocean.com/v2/domains/"+*domainName+"/records/"+strconv.Itoa(recordID), jsonStr)
		if err != nil {
			log.Println("Error updating subdomain record: ", err)
		}
		defer updateResp.Body.Close()
		contents, err := ioutil.ReadAll(updateResp.Body)
		if err != nil {
			log.Println("Error reading response body for update subdomain record: ", err)
		}
		log.Println("update results: ", string(contents))
	} else {
		// add
		log.Println("subdomain does not exist, adding...")
		jsonStr := []byte(`{"type":"A","name":"` + *subDomain + `","data":"` + currentIPJSON.CurrentIP + `","priority":null,"port":null,"weight":null, "ttl":300}`)
		addResp, err := httpReq("POST", "https://api.digitalocean.com/v2/domains/"+*domainName+"/records", jsonStr)
		if err != nil {
			log.Println("Error adding subdomain record: ", err)
		}
		defer addResp.Body.Close()
		contents, err := ioutil.ReadAll(addResp.Body)
		if err != nil {
			log.Println("Error reading response body for add subdomain record: ", err)
		}
		log.Println("add results: ", string(contents))
	}
}

func httpReq(method string, URL string, postData []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(postData))
	if err != nil {
		log.Println("Error with building request for "+URL+": ", err)
		return &http.Response{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+*apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error with request for "+URL+": ", err)
		return &http.Response{}, err
	}
	return resp, nil
}

func subExists(records []DomainRecord) (bool, int) {
	for _, record := range records {
		if record.Name == *subDomain {
			log.Println("Found record: ", record.Name)
			return true, record.Id
		}
	}
	return false, 0
}
