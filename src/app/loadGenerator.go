package app

import (
	"io/ioutil"
	"log"
	"net/http"
)

func pullUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
}
