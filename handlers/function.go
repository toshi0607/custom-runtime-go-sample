package main

import (
	"os"

	"github.com/toshi0607/go-custom-runtime-sample/runtime"
)

var (
	runtimeClient runtime.Client
)

func init() {
	endpoint := "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API")
	runtimeClient = runtime.NewClient(endpoint)
}

func main() {
	for {
		_, ec, _ := runtimeClient.NextInvocation()
		if err := runtimeClient.InvocationResponse(ec.AwsRequestId, []byte(ec.InvokedFunctionArn)); err != nil {
			runtimeClient.InvocationError(ec.AwsRequestId, runtime.ApiError{})
		}
	}
}
