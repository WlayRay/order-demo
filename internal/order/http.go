package main

import (
	"github.com/WlayRay/order-demo/order/app"
	"github.com/WlayRay/order-demo/order/app/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	//TODO implement me
	panic("implement me")
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerID: "fake-customer-ID",
		OrderID:    "fake-ID",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    o,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    o,
	})
}
