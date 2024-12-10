package models

import "time"

type Order struct {
	ID         int    `json:"orderID"`
	UserID     int    `json:"userID"`
	IsBid      bool   `json:"isBid"`
	Pair       string `json:"pair"`
	Price      string `json:"price"`
	Qty        string `json:"qty"`
	SizeFilled string `json:"sizeFilled"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	CreatedAt  string `json:"createdAt"`
	ClosedAt   string `json:"closedAt"`
}

type PlaceOrderReq struct {
	UserID int    `json:"userID"`
	IsBid  bool   `json:"isBid"`
	Pair   string `json:"pair"`
	Price  string `json:"price"`
	Qty    string `json:"qty"`
	Type   string `json:"type"` // Market или Limit
}

type SubscriptionRequest struct {
	Action string `json:"action"`
	Topic  string `json:"topic"`
	Params struct {
		Pair      string `json:"pair"`
		Precision int32  `json:"precision"`
	} `json:"params"`
}

type OrderBookSnapshot struct {
	Pair    string  `json:"pair"`
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
	Pair   string  `json:"pair"`
	Trades []Trade `json:"trades"`
}

type Trade struct {
	IsBid bool      `json:"isBid"`
	Price string    `json:"price"`
	Qty   string    `json:"qty"`
	Time  time.Time `json:"time"`
}

type PairParams struct {
	Pair            string  `json:"pair"`
	PricePrecisions []int32 `json:"pricePrecisions"`
	QtyPrecision    int32   `json:"qtyPrecision"`
}

type Ticker struct {
	Price     string `json:"price"`
	Change    string `json:"change"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	Turnover  string `json:"turnover"`
}
