package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	c := newTClient()
	err := c.get()
	if err != nil {
		log.Print(err)
	}
}

type tclient struct {
	url    string
	client *http.Client
}

func newTClient() *tclient {
	tr := &http.Transport{

		TLSHandshakeTimeout: 1 * time.Second,
		MaxIdleConns:        10,
		IdleConnTimeout:     1 * time.Second,
	}
	// 0 -- timeout on a client level
	// this timeout does not get sent to server perse, but it still
	// kills the request on a server via cancel() call on a function exit.
	client := &http.Client{Timeout: 10 * time.Second, Transport: tr}
	return &tclient{url: "http://localhost:8888", client: client}
}

// Verify takes names-strings and options and returns verification result.
func (tc *tclient) get() error {
	// 1 -- timeout via context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		"GET", tc.url+"/ping", nil)

	resp, err := tc.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("SUCCESS: '%s'", string(body))
	return nil
}
