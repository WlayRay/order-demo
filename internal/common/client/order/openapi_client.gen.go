// Package order provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetCustomerCustomerIDOrdersOrderID request
	GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostOrdersWithBody request with any body
	PostOrdersWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostOrders(ctx context.Context, body PostOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetCustomerCustomerIDOrdersOrderID(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCustomerCustomerIDOrdersOrderIDRequest(c.Server, customerID, orderID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostOrdersWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostOrdersRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostOrders(ctx context.Context, body PostOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostOrdersRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetCustomerCustomerIDOrdersOrderIDRequest generates requests for GetCustomerCustomerIDOrdersOrderID
func NewGetCustomerCustomerIDOrdersOrderIDRequest(server string, customerID string, orderID string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "customerID", runtime.ParamLocationPath, customerID)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/customer/%s/orders/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostOrdersRequest calls the generic PostOrders builder with application/json body
func NewPostOrdersRequest(server string, body PostOrdersJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostOrdersRequestWithBody(server, "application/json", bodyReader)
}

// NewPostOrdersRequestWithBody generates requests for PostOrders with any type of body
func NewPostOrdersRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/orders")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetCustomerCustomerIDOrdersOrderIDWithResponse request
	GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error)

	// PostOrdersWithBodyWithResponse request with any body
	PostOrdersWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostOrdersResponse, error)

	PostOrdersWithResponse(ctx context.Context, body PostOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostOrdersResponse, error)
}

type GetCustomerCustomerIDOrdersOrderIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Order
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetCustomerCustomerIDOrdersOrderIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCustomerCustomerIDOrdersOrderIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostOrdersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Order
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r PostOrdersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostOrdersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetCustomerCustomerIDOrdersOrderIDWithResponse request returning *GetCustomerCustomerIDOrdersOrderIDResponse
func (c *ClientWithResponses) GetCustomerCustomerIDOrdersOrderIDWithResponse(ctx context.Context, customerID string, orderID string, reqEditors ...RequestEditorFn) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
	rsp, err := c.GetCustomerCustomerIDOrdersOrderID(ctx, customerID, orderID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp)
}

// PostOrdersWithBodyWithResponse request with arbitrary body returning *PostOrdersResponse
func (c *ClientWithResponses) PostOrdersWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostOrdersResponse, error) {
	rsp, err := c.PostOrdersWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostOrdersResponse(rsp)
}

func (c *ClientWithResponses) PostOrdersWithResponse(ctx context.Context, body PostOrdersJSONRequestBody, reqEditors ...RequestEditorFn) (*PostOrdersResponse, error) {
	rsp, err := c.PostOrders(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostOrdersResponse(rsp)
}

// ParseGetCustomerCustomerIDOrdersOrderIDResponse parses an HTTP response from a GetCustomerCustomerIDOrdersOrderIDWithResponse call
func ParseGetCustomerCustomerIDOrdersOrderIDResponse(rsp *http.Response) (*GetCustomerCustomerIDOrdersOrderIDResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCustomerCustomerIDOrdersOrderIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Order
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParsePostOrdersResponse parses an HTTP response from a PostOrdersWithResponse call
func ParsePostOrdersResponse(rsp *http.Response) (*PostOrdersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostOrdersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Order
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
