package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	runtimeApiEndpointPrefix string
	nextEndpoint string
)

func init() {
	runtimeApiEndpointPrefix = "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
	nextEndpoint = runtimeApiEndpointPrefix + "next"
}

func main() {
	log.Println("handler started")
	fmt.Println("handler started2")

	for {
		func() {
			resp, _ := http.Get(nextEndpoint)
			defer func() {
				resp.Body.Close()
			}()

			rId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
			log.Printf("実行中のリクエストID" + rId)
			http.Post(respEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
		}()
	}
}

func respEndpoint(requestId string) string {
	return runtimeApiEndpointPrefix + requestId + "/response"
}
