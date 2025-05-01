package app

import "github.com/WlayRay/order-demo/stock/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
}
