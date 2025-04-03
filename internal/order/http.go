package main

import (
	"fmt"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/order/app"
	"github.com/WlayRay/order-demo/order/app/command"
	"github.com/WlayRay/order-demo/order/app/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customerID":  req.CustomerID,
		"orderID":     r.OrderID,
		"redirectURL": fmt.Sprintf("http://localhost:9000/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    o,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"Order": o,
		},
	})
}
