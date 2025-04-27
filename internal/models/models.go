package models

import "encoding/json"

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
	Action string          `json:"action"`
	Topic  string          `json:"topic"`
	Params json.RawMessage `json:"params"`
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

type Trade struct {
	Pair  string `json:"pair"`
	IsBid bool   `json:"isBid"`
	Price string `json:"price"`
	Qty   string `json:"qty"`
	Time  string `json:"time"`
}

type PairParams struct {
	Pair                     string   `json:"pair"`
	OrderBookPricePrecisions []int32  `json:"orderBookPricePrecisions"`
	QtyPrecision             int32    `json:"qtyPrecision"`
	CandleStickTimeframes    []string `json:"candleStickTimeframes"`
}

type Ticker struct {
	Pair      string `json:"pair"`
	LastPrice string `json:"lastPrice"`
	Change    string `json:"change"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	Volume    string `json:"volume"`
	Turnover  string `json:"turnover"`
}

type CandleStick struct {
	ID         int    `json:"ID"`
	Pair       string `json:"pair"`
	Timeframe  string `json:"timeframe"`
	OpenTime   string `json:"openTime"`
	CloseTime  string `json:"closeTime"`
	OpenPrice  string `json:"openPrice"`
	ClosePrice string `json:"closePrice"`
	HighPrice  string `json:"highPrice"`
	LowPrice   string `json:"lowPrice"`
	Volume     string `json:"volume"`
	Turnover   string `json:"turnover"`
	IsClosed   bool   `json:"isClosed"`
}

type CandleStickHistoryRequest struct {
	Pair      string `json:"pair"`
	Timeframe string `json:"timeframe"`
	CandleID  int    `json:"candleID"`
	Limit     int    `json:"limit"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
