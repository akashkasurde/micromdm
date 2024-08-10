package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func decodeDeviceActivationLockRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req dep.ActivationLockRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeDeviceActivationUnlockRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req dep.ActivationUnlockRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeActivationLockResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp dep.ActivationLockResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func decodeActivationUnlockResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp interface{}
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeEnableActivationLockEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dep.ActivationLockRequest)
		lockReponse, err := svc.EnableActivationLock(ctx, req.Device, req.EscrowKey, req.LostMessage)
		return dep.ActivationLockResponse{SerialNumber: lockReponse.SerialNumber, Status: lockReponse.Status}, err
	}
}

func (svc *DEPService) EnableActivationLock(ctx context.Context, serial string, escrow_key string, lost_message string) (*dep.ActivationLockResponse, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}

	lockRequest := &dep.ActivationLockRequest{
		Device:      serial,
		EscrowKey:   escrow_key,
		LostMessage: lost_message,
	}
	return svc.client.ActivationLock(lockRequest)
}

func (e Endpoints) EnableActivationLock(ctx context.Context, serial string, escrow_key string, lost_message string) (*dep.ActivationLockResponse, error) {
	lockRequest := &dep.ActivationLockRequest{
		Device:      serial,
		EscrowKey:   escrow_key,
		LostMessage: lost_message,
	}
	response, err := e.EnableActivationLockEndpoint(ctx, lockRequest)
	return response.(*dep.ActivationLockResponse), err
}

func MakeEnableActivationUnlockEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dep.ActivationUnlockRequest)
		unlockResponse, err := svc.EnableActivationUnlock(ctx, req.Querystring, req.Messagebody)
		return unlockResponse, err
	}
}

func (svc *DEPService) EnableActivationUnlock(ctx context.Context, querystring string, messagebody string) (interface{}, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}

	unlockRequest := &dep.ActivationUnlockRequest{
		Querystring: querystring,
		Messagebody: messagebody,
	}
	return svc.client.ActivationUnlock(unlockRequest)
}

func (e Endpoints) EnableActivationUnlock(ctx context.Context, querystring string, messagebody string) (interface{}, error) {
	unlockRequest := &dep.ActivationUnlockRequest{
		Querystring: querystring,
		Messagebody: messagebody,
	}
	response, err := e.EnableActivationUnlockEndpoint(ctx, unlockRequest)
	return response, err
}
