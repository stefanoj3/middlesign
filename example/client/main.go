package main

import (
	"fmt"
	"github.com/stefanoj3/middlesign"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8000/"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := middlesign.DefaultConfig("my_super_secret")

	q := req.URL.Query()
	q.Add(c.TimestampKey, time.Now().Format(c.TimestampFormat))

	signature := middlesign.SignString(q.Encode(), c.Secret)
	q.Add(c.SignatureKey, signature)

	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response", res)
}
