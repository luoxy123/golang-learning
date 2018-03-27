package feiniubus

import (
	"net/http"
)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return NewConfig().
		WithHTTPClient(http.DefaultClient).
		WithDisableSSL(true)
}

// DefaultHandlers returns the default request handlers
func DefaultHandlers() Handlers {
	var handlers Handlers

	handlers.Validate.PushBackNamed(ValidateEndpointHandler)
	handlers.Validate.AfterEachFn = HandlerListStopOnError
	handlers.Build.PushBackNamed(BuildHandler)
	handlers.Build.AfterEachFn = HandlerListStopOnError
	handlers.Build.PushBackNamed(SDKVersionUserAgentHandler)
	handlers.Build.PushBackNamed(BuildContentLengthHandler)
	handlers.Send.PushBackNamed(SendHandler)
	handlers.ValidateResponse.PushBackNamed(ValidateResponseHandler)
	handlers.UnmarshalError.PushBackNamed(UnmarshalErrorHandler)
	handlers.Unmarshal.PushBackNamed(UnmarshalHandler)

	return handlers
}
