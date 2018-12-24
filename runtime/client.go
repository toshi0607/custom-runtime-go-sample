package runtime

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const (
	RUNTIME_API_VERSION = "2018-06-01"
	CONTENT_TYPE        = "application/json"
)

type Client interface {
	NextInvocation() ([]byte, *EventContext, error)
	InvocationResponse(awsRequestId string, content []byte) error
	InvocationError(awsRequestId string, err error) error
	InitializationError(err error) error
}

// https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/runtimes-api.html#runtimes-api-next
type EventContext struct {
	InvokedFunctionArn string
	AwsRequestId       string
	RuntimeTraceId     string
	DeadlineMs         int64
	ClientContext      ClientContext
	CognitoIdentity    CognitoIdentity
}

type CognitoIdentity struct {
	identityId     string
	identityPoolId string
}

type ClientContext struct {
	client      ClientApplication
	custom      map[string]string
	environment map[string]string
}

type ClientApplication struct {
	installationId string
	appTitle       string
	appVersionName string
	appVersionCode string
	appPackageName string
}

type client struct {
	endpoint string
	client   *http.Client
}

func NewClient(endpoint string) Client {
	return &client{
		endpoint: endpoint,
		client:   http.DefaultClient,
	}
}

// http://${AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime/invocation/next

// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-next
func (c *client) NextInvocation() ([]byte, *EventContext, error) {
	url := c.endpoint + "/" + RUNTIME_API_VERSION + "/runtime/invocation/next"
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, nil, errors.New("failed to get a response")
	}
	defer resp.Body.Close()
	ec, err := getEventContext(resp.Header)
	if err != nil {
		return nil, nil, errors.New("failed to get an event context")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.New("failed to get an event")
	}

	return b, ec, nil
}

// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-response
func (c *client) InvocationResponse(awsRequestId string, content []byte) error {
	// /runtime/invocation/AwsRequestId/response
	url := c.endpoint + "/" + RUNTIME_API_VERSION + "/runtime/invocation/" + awsRequestId + "/response"

	resp, err := c.client.Post(url, CONTENT_TYPE, bytes.NewBuffer(content))
	if err != nil {
		return errors.Wrap(err, "failed to post content")
	}
	defer resp.Body.Close()

	return nil
}

// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-invokeerror
func (c *client) InvocationError(awsRequestId string, aerr error) error {
	// /runtime/invocation/AwsRequestId/error
	url := c.endpoint + "/" + RUNTIME_API_VERSION + "/runtime/invocation/" + awsRequestId + "/error"
	errs, err := json.Marshal(aerr)
	if err != nil {
		return errors.Wrap(err, "failed to post content")
	}

	resp, err := c.client.Post(url, CONTENT_TYPE, bytes.NewBuffer(errs))
	if err != nil {
		return errors.Wrap(err, "failed to post content")
	}
	defer resp.Body.Close()

	return nil
}

// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html#runtimes-api-initerror
func (c *client) InitializationError(aerr error) error {
	// /runtime/init/error
	url := c.endpoint + "/" + RUNTIME_API_VERSION + "/runtime/init/error"
	errs, err := json.Marshal(aerr)
	if err != nil {
		return errors.Wrap(err, "failed to post content")
	}

	resp, err := c.client.Post(url, CONTENT_TYPE, bytes.NewBuffer(errs))
	if err != nil {
		return errors.Wrap(err, "failed to post content")
	}
	defer resp.Body.Close()

	return nil
}

func getEventContext(header http.Header) (*EventContext, error) {
	invokedFunctionArn := header.Get("Lambda-Runtime-Invoked-Function-Arn")
	awsRequestId := header.Get("Lambda-Runtime-Aws-Request-Id")
	runtimeTraceId := header.Get("Lambda-Runtime-Trace-Id")

	runtimeDeadlineMs, err := strconv.ParseInt(header.Get("Lambda-Runtime-Deadline-Ms"), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Lambda-Runtime-Deadline-Ms")
	}

	ev := &EventContext{
		InvokedFunctionArn: invokedFunctionArn,
		AwsRequestId:       awsRequestId,
		RuntimeTraceId:     runtimeTraceId,
		DeadlineMs:         runtimeDeadlineMs,
	}

	runtimeClientContext := header.Get("Lambda-Runtime-Client-Context")
	if runtimeClientContext != "" {
		var clientContext ClientContext
		if err := json.Unmarshal([]byte(runtimeClientContext), &clientContext); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal ClientContext")
		}

		ev.ClientContext = clientContext
	}

	runtimeCognitoIdentity := header.Get("Lambda-Runtime-Cognito-Identity") // 型
	if runtimeCognitoIdentity != "" {
		var cognitoIdentity CognitoIdentity
		if err := json.Unmarshal([]byte(runtimeClientContext), &cognitoIdentity); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal ClientContext")
		}

		ev.CognitoIdentity = cognitoIdentity
	}

	return ev, nil
}
