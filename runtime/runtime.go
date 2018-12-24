package runtime

import (
	"encoding/json"
	"log"
)

type Handler interface {
	Run(ctx Context, event []byte) ([]byte, error)
}

type handler func(ctx Context, event []byte) ([]byte, error)

func (h handler) Run(ctx Context, event []byte) ([]byte, error) {
	response, err := h(ctx, event)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

func Start(h handler) {
	startWithConfig(h, NewConfigProvider())
}

func startWithConfig(h handler, config ConfigProvider) {
	endpoint := config.GetRuntimeApiEndpoint()
	settings, err := config.GetFunctionSettings()
	if endpoint == "" || err != nil {
		log.Fatal("failed to init lambda")
	}

	startWithRuntimeClient(h, *settings, NewClient(endpoint))

}

func startWithRuntimeClient(h handler, s FunctionSettings, c Client) {
	NewRuntime(c, h, s).Start()
}

type Runtime struct {
	client   Client
	handler  Handler
	settings FunctionSettings
}

type Context struct {
	MemorySize         int32
	FunctionName       string
	FunctionVersion    string
	InvokedFunctionArn string
	AwsRequestId       string
	XrayTraceId        string
	LogStreamName      string
	LogGroupName       string
	ClientContext      ClientContext
	Identity           CognitoIdentity
	Deadline           int64
}

func NewRuntime(client Client, handler Handler, settings FunctionSettings) *Runtime {
	return &Runtime{
		client:   client,
		handler:  handler,
		settings: settings,
	}
}

func (r *Runtime) Start() {
	for {
		ctx, ev := r.getNextEvent()
		requestId := ctx.AwsRequestId
		result, err := r.handler.Run(ctx, ev)
		if err != nil {
			err := r.client.InvocationError(requestId, err)
			if err != nil {
				log.Fatal("failed to invoke an error")
			}
		}
		err = r.client.InvocationResponse(requestId, result)
		if err != nil {
			err := r.client.InvocationError(requestId, err)
			if err != nil {
				log.Fatal("failed to invoke an error")
			}
		}
	}
}

func (r *Runtime) getNextEvent() (Context, []byte) {
	ev, ctx, err := r.client.NextInvocation()
	if err != nil {
		r.client.InitializationError(err)
	}
	return Context{
		MemorySize:         r.settings.memorySize,
		FunctionName:       r.settings.functionName,
		FunctionVersion:    r.settings.version,
		LogStreamName:      r.settings.logStream,
		LogGroupName:       r.settings.logGroup,
		InvokedFunctionArn: ctx.InvokedFunctionArn,
		AwsRequestId:       ctx.AwsRequestId,
		XrayTraceId:        ctx.RuntimeTraceId,
		ClientContext:      ctx.ClientContext,
		Identity:           ctx.CognitoIdentity,
		Deadline:           ctx.DeadlineMs,
	}, ev
}
