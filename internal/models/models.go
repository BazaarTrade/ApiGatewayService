package models

import "time"

type Order struct {
	ID         int64  `json:"orderID"`
	UserID     int64  `json:"userID"`
	IsBid      bool   `json:"isBid"`
	Symbol     string `json:"symbol"`
	Price      string `json:"price"`
	Qty        string `json:"qty"`
	SizeFilled string `json:"sizeFilled"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	CreatedAt  string `json:"createdAt"`
	ClosedAt   string `json:"closedAt"`
}

type PlaceOrderReq struct {
	UserID int64  `json:"userID"`
	IsBid  bool   `json:"isBid"`
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Qty    string `json:"qty"`
	Type   string `json:"type"` // Market или Limit
}

type SubscriptionRequest struct {
	Action string `json:"action"`
	Topic  string `json:"topic"`
	Params struct {
		Symbol    string `json:"symbol"`
		Precision int32  `json:"precision"`
	} `json:"params"`
}

type OrderBookSnapshot struct {
	Symbol  string  `json:"symbol"`
	Bids    []Limit `json:"bids"`
	Asks    []Limit `json:"asks"`
	BidsQty string  `json:"bidsQty"`
	AsksQty string  `json:"asksQty"`
}

type Limit struct {
	Price string `json:"price"`
	Qty   string `json:"qty"`
}

type Trades struct {
	Symbol string  `json:"symbol"`
	Trades []Trade `json:"trades"`
}

type Trade struct {
	IsBid bool      `json:"isBid"`
	Price string    `json:"price"`
	Qty   string    `json:"qty"`
	Time  time.Time `json:"time"`
}
