package backend

import (
	"github.com/fbarcena/grafana-plugin-sdk-go/backend/plugin"
)

const defaultServerMaxReceiveMessageSize = 1024 * 1024 * 16

// GRPCSettings settings for gRPC.
type GRPCSettings struct {
	// MaxReceiveMsgSize the max gRPC message size in bytes the plugin can receive.
	// If this is <= 0, gRPC uses the default 16MB.
	MaxReceiveMsgSize int

	// MaxSendMsgSize the max gRPC message size in bytes the plugin can send.
	// If this is <= 0, gRPC uses the default `math.MaxInt32`.
	MaxSendMsgSize int
}

//ServeOpts options for serving plugins.
type ServeOpts struct {
	// CheckHealthHandler handler for health checks.
	CheckHealthHandler CheckHealthHandler

	// CallResourceHandler handler for resource calls.
	// Optional to implement.
	CallResourceHandler CallResourceHandler

	// QueryDataHandler handler for data queries.
	// Required to implement if data source.
	QueryDataHandler QueryDataHandler

	// TransformDataHandler handler for data transformations.
	// Very experimental and shouldn't be implemented in most cases.
	// Optional to implement.
	TransformDataHandler TransformDataHandler

	// GRPCSettings settings for gPRC.
	GRPCSettings GRPCSettings
}

// Serve starts serving the plugin over gRPC.
func Serve(opts ServeOpts) error {
	sdkAdapter := &sdkAdapter{
		CheckHealthHandler:   opts.CheckHealthHandler,
		CallResourceHandler:  opts.CallResourceHandler,
		QueryDataHandler:     opts.QueryDataHandler,
		TransformDataHandler: opts.TransformDataHandler,
	}

	pluginOpts := plugin.ServeOpts{
		DiagnosticsServer: sdkAdapter,
	}

	if opts.CallResourceHandler != nil {
		pluginOpts.ResourceServer = sdkAdapter
	}

	if opts.QueryDataHandler != nil {
		pluginOpts.DataServer = sdkAdapter
	}

	if opts.TransformDataHandler != nil {
		pluginOpts.TransformServer = sdkAdapter
	}

	return plugin.Serve(pluginOpts)
}
