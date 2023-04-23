package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	errors2 "github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
)

type Transport struct {
	httpClient *http.Client
	BaseURL    string
	BasePath   string
}

func New(baseURL, basePath string) transport.Transport {
	return &Transport{
		httpClient: &http.Client{},
		BaseURL:    baseURL,
		BasePath:   basePath,
	}
}

func (t *Transport) Create(context context.Context, req *utils.CreateAccountRequest) (*utils.CreateAccountResponse, error) {
	url := t.BaseURL + t.BasePath

	body, err := utils.EncodeJSONRequest(req)
	if err != nil {
		return nil, &errors2.ErrBadRequest{Detail: "failed to marshal request body"}
	}

	httpReq, err := http.NewRequestWithContext(context, http.MethodPost, url, body)
	if err != nil {
		return nil, &errors2.ErrBadRequest{Detail: "failed to create HTTP request"}
	}

	httpReq.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors2.ErrPermanentFailure{Detail: fmt.Sprintf("failed to send HTTP request: %v", err)}
	}
	defer resp.Body.Close()

	if err := errors2.HandleHTTPError(resp); err != nil {
		return nil, err
	}

	var createResp utils.CreateAccountResponse
	if err := utils.DecodeJSONResponse(resp, &createResp); err != nil {
		return nil, &errors2.ErrPermanentFailure{Detail: "failed to unmarshal response body"}
	}

	return &createResp, nil
}

func (t *Transport) Fetch(context context.Context, accountID string) (*utils.FetchAccountResponse, error) {
	url := t.BaseURL + t.BasePath + accountID

	httpReq, err := http.NewRequestWithContext(context, http.MethodGet, url, nil)
	if err != nil {
		return nil, &errors2.ErrBadRequest{Detail: "failed to create HTTP request"}
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors2.ErrPermanentFailure{Detail: fmt.Sprintf("failed to send HTTP request: %v", err)}
	}
	defer resp.Body.Close()

	if err := errors2.HandleHTTPError(resp); err != nil {
		return nil, err
	}

	var fetchResp utils.FetchAccountResponse
	if err := utils.DecodeJSONResponse(resp, &fetchResp); err != nil {
		return nil, &errors2.ErrPermanentFailure{Detail: "failed to unmarshal response body"}
	}

	return &fetchResp, nil
}

func (t *Transport) Delete(context context.Context, req *utils.DeleteAccountRequest) error {
	url := t.BaseURL + t.BasePath + req.ID + "?version=" + strconv.FormatInt(req.Version, 10)

	httpReq, err := http.NewRequestWithContext(context, http.MethodDelete, url, nil)
	if err != nil {
		return &errors2.ErrBadRequest{Detail: "failed to create HTTP request"}
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return &errors2.ErrPermanentFailure{Detail: fmt.Sprintf("failed to send HTTP request: %v", err)}
	}
	defer resp.Body.Close()

	if err := errors2.HandleHTTPError(resp); err != nil {
		return err
	}

	return nil
}
