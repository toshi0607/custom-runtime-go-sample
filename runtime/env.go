package runtime

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type ConfigProvider interface {
	GetFunctionSettings() (*FunctionSettings, error)
	GetRuntimeApiEndpoint() string
}

type FunctionSettings struct {
	functionName string
	memorySize   int32
	version      string
	logStream    string
	logGroup     string
}

type configProvider struct{}

func NewConfigProvider() ConfigProvider {
	return &configProvider{}
}

func (c configProvider) GetFunctionSettings() (*FunctionSettings, error) {
	functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	version := os.Getenv("AWS_LAMBDA_FUNCTION_VERSION")
	logStream := os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME")
	logGroup := os.Getenv("AWS_LAMBDA_LOG_GROUP_NAME")
	memoryStr := os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE")
	memorySize, err := strconv.Atoi(memoryStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse memoryStr: %s to int", memoryStr)
	}

	return &FunctionSettings{
		functionName: functionName,
		memorySize:   int32(memorySize),
		version:      version,
		logStream:    logStream,
		logGroup:     logGroup,
	}, nil
}

func (c configProvider) GetRuntimeApiEndpoint() string {
	return "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API")
}
