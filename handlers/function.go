package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

func main() {
	runtimeApiEndpointPrefix := "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
	// next API のエンドポイント
	nextEndpoint := runtimeApiEndpointPrefix + "next"

	for {
		resp, _ := http.Get(nextEndpoint)
		defer resp.Body.Close()

		rId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
		log.Printf("実行中のリクエストID" + rId)
		http.Post(runtimeApiEndpointPrefix + rId + "/response", "application/json", bytes.NewBuffer([]byte(rId)))
	}
}
