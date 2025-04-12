package main

import (
	"fmt"
	"github.com/WlayRay/order-demo/common"
	client "github.com/WlayRay/order-demo/common/client/order"
	"github.com/WlayRay/order-demo/order/app"
	"github.com/WlayRay/order-demo/order/app/command"
	"github.com/WlayRay/order-demo/order/app/dto"
	"github.com/WlayRay/order-demo/order/app/query"
	"github.com/WlayRay/order-demo/order/convertor"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
	common.BaseResponse
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		H.Response(c, err, nil)
		return
	}

	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.GetItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		H.Response(c, err, nil)
		return
	}

	resp.CustomerID = req.CustomerID
	resp.OrderID = r.OrderID
	resp.RedirectURL = fmt.Sprintf("http://localhost:9000/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)
	H.Response(c, err, resp)
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp struct {
			Order *domain.Order `json:"order"`
		}
	)

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		H.Response(c, err, nil)
		return
	}

	resp.Order = o
	H.Response(c, err, resp)
}
