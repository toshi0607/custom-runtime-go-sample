package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/toshi0607/go-custom-runtime-sample/runtime"
)

type MyEvent struct {
	Name string `json:"name"`
}

func MyHandler(ctx runtime.Context, event []byte) ([]byte, error) {
	log.Printf("request id: %s\n", ctx.AwsRequestId)

	var me MyEvent
	json.Unmarshal(event, &me)

	str := fmt.Sprintf("Hello %s!", me.Name)
	log.Println(str)
	return []byte(str), nil
}

func main() {
	runtime.Start(MyHandler)
}
