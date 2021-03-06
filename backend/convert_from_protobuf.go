package backend

import (
	"time"

	"github.com/fbarcena/grafana-plugin-sdk-go/data"
	"github.com/fbarcena/grafana-plugin-sdk-go/genproto/pluginv2"
)

type ConvertFromProtobuf struct {
}

func FromProto() ConvertFromProtobuf {
	return ConvertFromProtobuf{}
}

// User converts proto version of user to SDK version
func (f ConvertFromProtobuf) User(user *pluginv2.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func (f ConvertFromProtobuf) DataSourceConfig(proto *pluginv2.DataSourceConfig) *DataSourceConfig {
	if proto == nil {
		return nil
	}

	return &DataSourceConfig{
		ID:                      proto.Id,
		Name:                    proto.Name,
		URL:                     proto.Url,
		User:                    proto.User,
		Database:                proto.Database,
		BasicAuthEnabled:        proto.BasicAuthEnabled,
		BasicAuthUser:           proto.BasicAuthUser,
		JSONData:                proto.JsonData,
		DecryptedSecureJSONData: proto.DecryptedSecureJsonData,
		Updated:                 time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
	}
}

func (f ConvertFromProtobuf) PluginConfig(proto *pluginv2.PluginConfig) PluginConfig {
	return PluginConfig{
		OrgID:                   proto.OrgId,
		PluginID:                proto.PluginId,
		JSONData:                proto.JsonData,
		DecryptedSecureJSONData: proto.DecryptedSecureJsonData,
		Updated:                 time.Unix(0, proto.LastUpdatedMS*int64(time.Millisecond)),
		DataSourceConfig:        f.DataSourceConfig(proto.DatasourceConfig),
	}
}

func (f ConvertFromProtobuf) TimeRange(proto *pluginv2.TimeRange) TimeRange {
	return TimeRange{
		From: time.Unix(0, proto.FromEpochMS*int64(time.Millisecond)),
		To:   time.Unix(0, proto.ToEpochMS*int64(time.Millisecond)),
	}
}

func (f ConvertFromProtobuf) DataQuery(proto *pluginv2.DataQuery) *DataQuery {
	return &DataQuery{
		RefID:         proto.RefId,
		MaxDataPoints: proto.MaxDataPoints,
		TimeRange:     f.TimeRange(proto.TimeRange),
		Interval:      time.Duration(proto.IntervalMS) * time.Millisecond,
		JSON:          []byte(proto.Json),
	}
}

func (f ConvertFromProtobuf) QueryDataRequest(protoReq *pluginv2.QueryDataRequest) *QueryDataRequest {
	queries := make([]DataQuery, len(protoReq.Queries))
	for i, q := range protoReq.Queries {
		queries[i] = *f.DataQuery(q)
	}

	return &QueryDataRequest{
		PluginConfig: f.PluginConfig(protoReq.Config),
		Headers:      protoReq.Headers,
		Queries:      queries,
		User:         f.User(protoReq.User),
	}
}

func (f ConvertFromProtobuf) QueryDataResponse(protoRes *pluginv2.QueryDataResponse) (*QueryDataResponse, error) {
	qrd := &QueryDataResponse{
		Responses: make(Responses, len(protoRes.Responses)),
	}
	for refID, res := range protoRes.Responses {
		dr := DataResponse{
			JSON: res.Json,
		}
		qrd.Responses[refID] = dr
	}
	///
	frames := make([]*data.Frame, len(protoRes.Frames))
	var err error
	for i, encodedFrame := range protoRes.Frames {
		frames[i], err = data.UnmarshalArrow(encodedFrame)
		if err != nil {
			return nil, err
		}
	}
	///
	qrd.Metadata = protoRes.Metadata
	qrd.Frames = frames
	qrd.Json = protoRes.Json
	// return &QueryDataResponse{Metadata: protoRes.Metadata, Frames: frames, Json: protoRes.Json}, nil
	return qrd, nil
}

func (f ConvertFromProtobuf) CallResourceRequest(protoReq *pluginv2.CallResourceRequest) *CallResourceRequest {
	headers := map[string][]string{}
	for k, values := range protoReq.Headers {
		headers[k] = values.Values
	}

	return &CallResourceRequest{
		PluginConfig: f.PluginConfig(protoReq.Config),
		Path:         protoReq.Path,
		Method:       protoReq.Method,
		URL:          protoReq.Url,
		Headers:      headers,
		Body:         protoReq.Body,
		User:         f.User(protoReq.User),
	}
}

// HealthCheckRequest converts proto version to SDK version.
func (f ConvertFromProtobuf) HealthCheckRequest(protoReq *pluginv2.CheckHealthRequest) *CheckHealthRequest {
	return &CheckHealthRequest{
		PluginConfig: f.PluginConfig(protoReq.Config),
	}
}
