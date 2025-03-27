package app

import "github.com/WlayRay/order-demo/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePaymentLink command.CreatePaymentHandler
}
