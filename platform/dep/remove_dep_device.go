package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func decodeDisownDeviceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req dep.DeviceDisownRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeDisownDeviceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp interface{}
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeDisownDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dep.DeviceDisownRequest)
		disownReponse, err := svc.DisownDevices(ctx, req.Devices)
		return disownReponse, err
	}
}

func (svc *DEPService) DisownDevices(ctx context.Context, devices []string) (interface{}, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}

	disownRequest := &dep.DeviceDisownRequest{
		Devices: devices,
	}
	return svc.client.DisownDevice(disownRequest)
}

func (e Endpoints) DisownDevices(ctx context.Context, devices []string) (interface{}, error) {
	disownRequest := &dep.DeviceDisownRequest{
		Devices: devices,
	}
	response, err := e.DisownDeviceEndPoint(ctx, disownRequest)
	return response, err
}
