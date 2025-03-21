// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package ports

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /customer/{customerID}/orders)
	PostCustomerCustomerIDOrders(c *gin.Context, customerID string)

	// (GET /customer/{customerID}/orders/{orderID})
	GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// PostCustomerCustomerIDOrders operation middleware
func (siw *ServerInterfaceWrapper) PostCustomerCustomerIDOrders(c *gin.Context) {

	var err error

	// ------------- Path parameter "customerID" -------------
	var customerID string

	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostCustomerCustomerIDOrders(c, customerID)
}

// GetCustomerCustomerIDOrdersOrderID operation middleware
func (siw *ServerInterfaceWrapper) GetCustomerCustomerIDOrdersOrderID(c *gin.Context) {

	var err error

	// ------------- Path parameter "customerID" -------------
	var customerID string

	err = runtime.BindStyledParameterWithOptions("simple", "customerID", c.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter customerID: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "orderID" -------------
	var orderID string

	err = runtime.BindStyledParameterWithOptions("simple", "orderID", c.Param("orderID"), &orderID, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter orderID: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCustomerCustomerIDOrdersOrderID(c, customerID, orderID)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/customer/:customerID/orders", wrapper.PostCustomerCustomerIDOrders)
	router.GET(options.BaseURL+"/customer/:customerID/orders/:orderID", wrapper.GetCustomerCustomerIDOrdersOrderID)
}
