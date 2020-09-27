package backend

import (
	"github.com/fbarcena/grafana-plugin-sdk-go/data"
	"github.com/fbarcena/grafana-plugin-sdk-go/genproto/pluginv2"
)

//"github.com/grafana/grafana-plugin-sdk-go/data"

// QueryDataResponse contains the results from a QueryDataRequest.
// It is the return type of a QueryData call.
type QueryDataResponse struct {
	// Responses is a map of RefIDs (Unique Query ID) to *DataResponse.
	Responses Responses
	Frames    []*data.Frame
	Metadata  map[string]string
	Json      string
}

// NewQueryDataResponse returns a QueryDataResponse with the Responses property initialized.
func NewQueryDataResponse() *QueryDataResponse {
	return &QueryDataResponse{
		Responses: make(Responses),
	}
}

// Responses is a map of RefIDs (Unique Query ID) to DataResponses.
// The QueryData method the QueryDataHandler method will set the RefId
// property on the DataRespones' frames based on these RefIDs.
type Responses map[string]DataResponse

// DataResponse contains the results from a DataQuery.
// A map of RefIDs (unique query identifers) to this type makes up the Responses property of a QueryDataResponse.
// The Error property is used to allow for partial success responses from the containing QueryDataResponse.
type DataResponse struct {
	// The data returned from the Query. Each Frame repeats the RefID.
	Frames data.Frames

	// Error is a property to be set if the the corresponding DataQuery has an error.
	Error error

	// Add Tables and Timeseries for compatibility with 6.7
	Series []*pluginv2.TimeSeries
	Tables []*pluginv2.Table

	JSON string
}
