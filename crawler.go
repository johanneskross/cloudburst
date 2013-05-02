package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const URL = "http://www.example.com"

var domains = []string{"domain1", "domain2", "domain3", "domain4", "domain5"}

func pullUrl(domain string, url string, c chan bool) {
	fmt.Printf("started %s\n", domain)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("finished %s\n", domain)
	c <- true
}

func startDomains() {
	// start domains and pull url on each domain
	c := make(chan bool, len(domains))
	for _, v := range domains {
		go pullUrl(v, URL, c)
	}

	// wait until all go routines are ended
	for i := 0; i < cap(c); i++ {
		<-c
	}
}

func main() {
	startDomains()
}
