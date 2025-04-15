package ws

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

type Hub struct {
	Users       map[int]*User
	mu          sync.RWMutex
	Subscribers map[string]map[any]*Subscribers
	logger      *slog.Logger
}

type User struct {
	ID      int
	Clients map[*Client]bool
}

type Client struct {
	Conn   *websocket.Conn
	Topics map[string]any
	mu     sync.RWMutex
}

type Subscribers struct {
	Clients map[*Client]bool
}

type OrderBookParams struct {
	Pair      string `json:"pair"`
	Precision int32  `json:"precision"`
}

type TradesParams struct {
	Pair string `json:"pair"`
}

type TickerParams struct {
	Pair string `json:"pair"`
}

type CandleStickParams struct {
	Pair      string `json:"pair"`
	Timeframe string `json:"timeframe"`
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		Users:       make(map[int]*User),
		Subscribers: make(map[string]map[any]*Subscribers),
		logger:      logger,
		mu:          sync.RWMutex{},
	}
}

func (h *Hub) HandleWebsocket(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userID",
		})
	}

	if userID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userID",
		})
	}

	conn, err := websocket.Accept(c.Response(), c.Request(), &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to upgrade connection to WebSocket",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	user, ok := h.Users[userID]
	if !ok {
		h.Users[userID] = &User{
			ID:      userID,
			Clients: make(map[*Client]bool),
		}
		user = h.Users[userID]
	}

	var client = &Client{
		Conn:   conn,
		Topics: make(map[string]any),
		mu:     sync.RWMutex{},
	}

	user.Clients[client] = true

	go h.readPump(client, userID)

	return nil
}

func (h *Hub) readPump(c *Client, userID int) {
	defer func() {
		c.Conn.Close(websocket.StatusNormalClosure, "normal closure")

		h.mu.Lock()
		defer h.mu.Unlock()
		//delete subscriber
		for topic, params := range c.Topics {
			if subscribers, exists := h.Subscribers[topic][params]; exists {
				delete(subscribers.Clients, c)
			}
		}

		//delete client, if no clients left - delete user
		if user, exists := h.Users[userID]; exists {
			delete(user.Clients, c)
			if len(h.Users[userID].Clients) == 0 {
				delete(h.Users, userID)
			}
		}
	}()

	for {
		_, msg, err := c.Conn.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway || errors.Is(err, io.EOF) {
				h.logger.Debug("client disconnected from websocket", "userID", userID)
				return
			}
			h.logger.Error("failed to read websocket connection", "error", err)
			return
		}

		var request models.SubscriptionRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			h.logger.Error("failed unmarshal SubscriptionRequest", "error", err)
			continue
		}

		valid, params, err := h.unmarshalParams(request)
		if err != nil {
			continue
		}

		if !valid {
			messageJSON, err := json.Marshal(map[string]string{
				"error": "invalid request",
			})
			if err != nil {
				h.logger.Error("failed to marshal error message", "error", err)
				continue
			}

			if err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON); err != nil {
				h.logger.Error("failed to write error message to websocket", "error", err)
				continue
			}
			continue
		}

		switch request.Action {
		case "subscribe":
			h.subscribeClient(c, request.Topic, params)

		case "unsubscribe":
			h.unsubscribeClient(c, request.Topic, params)

		default:
			h.logger.Error("attempt to perform a non-existent action")
		}
	}
}

func (h *Hub) unmarshalParams(request models.SubscriptionRequest) (bool, any, error) {
	switch request.Topic {
	case "orderBook":
		var params OrderBookParams
		if err := json.Unmarshal(request.Params, &params); err != nil {
			h.logger.Error("failed to unmarshal orderBook params", "error", err)
			return false, nil, err
		}

		if params.Pair == "" {
			h.logger.Error("invalid parameters for orderBook subscription")
			return false, nil, nil
		}
		return true, params, nil

	case "trades":
		var params TradesParams
		if err := json.Unmarshal(request.Params, &params); err != nil {
			h.logger.Error("failed to unmarshal trades params", "error", err)
			return false, nil, err
		}

		if params.Pair == "" {
			h.logger.Error("invalid parameters for trades subscription")
			return false, nil, nil
		}
		return true, params, nil

	case "ticker":
		var params TickerParams
		if err := json.Unmarshal(request.Params, &params); err != nil {
			h.logger.Error("failed to unmarshal ticker params", "error", err)
			return false, nil, err
		}

		if params.Pair == "" {
			h.logger.Error("invalid parameters for ticker subscription")
			return false, nil, nil
		}
		return true, params, nil

	case "candleStick":
		var params CandleStickParams
		if err := json.Unmarshal(request.Params, &params); err != nil {
			h.logger.Error("failed to unmarshal candleStick params", "error", err)
			return false, nil, err
		}

		if params.Pair == "" || params.Timeframe == "" {
			h.logger.Error("invalid parameters for candleStick subscription")
			return false, nil, nil
		}
		return true, params, nil

	default:
		h.logger.Error("attempt to subscribe to a non-existent topic")
		return false, nil, nil
	}
}

func (h *Hub) subscribeClient(c *Client, topic string, params any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers[topic]; !exists {
		messageJSON, err := json.Marshal(map[string]string{
			"error": "no such topic exists",
		})
		if err != nil {
			h.logger.Error("failed to marshal error message", "error", err)
			return
		}

		if err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON); err != nil {
			h.logger.Error("failed to write error message to websocket", "error", err)
			return
		}
		return
	}

	if _, exists := h.Subscribers[topic][params]; !exists {
		messageJSON, err := json.Marshal(map[string]string{
			"error": "this topic does not have such parameters",
		})
		if err != nil {
			h.logger.Error("failed to marshal error message", "error", err)
			return
		}

		if err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON); err != nil {
			h.logger.Error("failed to write error message to websocket", "error", err)
			return
		}
		return
	}

	h.Subscribers[topic][params].Clients[c] = true
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Topics[topic] = params

	subscriptionMessage := struct {
		Topic  string `json:"topic"`
		Status string `json:"status"`
	}{
		Topic:  topic,
		Status: "subscribed",
	}

	messageJSON, err := json.Marshal(subscriptionMessage)
	if err != nil {
		h.logger.Error("failed to marshal subscription message", "error", err)
		return
	}

	if err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON); err != nil {
		h.logger.Error("failed to write subscription message to websocket", "error", err)
		return
	}
}

func (h *Hub) unsubscribeClient(c *Client, topic string, params any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, exists := h.Subscribers[topic][params]; exists {
		delete(subscribers.Clients, c)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Topics, topic)

	subscriptionMessage := struct {
		Topic  string `json:"topic"`
		Status string `json:"status"`
	}{
		Topic:  topic,
		Status: "unsubscribed",
	}

	messageJSON, err := json.Marshal(subscriptionMessage)
	if err != nil {
		h.logger.Error("failed to marshal subscription message", "error", err)
		return
	}

	if err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON); err != nil {
		h.logger.Error("failed to write subscription message to websocket", "error", err)
		return
	}
}

func (h *Hub) AddOrderBookSnapshotTopic(pair string, orderBookPricePrecisions []int32) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers["orderBook"]; !exists {
		h.Subscribers["orderBook"] = make(map[any]*Subscribers)
	}

	for _, pricePrecision := range orderBookPricePrecisions {
		h.Subscribers["orderBook"][OrderBookParams{Pair: pair, Precision: pricePrecision}] = &Subscribers{Clients: map[*Client]bool{}}
	}
}

func (h *Hub) RemoveOrderBookSnapshotTopic(pair string, pricePrecisions []int32) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topic, exists := h.Subscribers["orderBook"]; exists {
		for _, pricePrecision := range pricePrecisions {
			delete(topic, OrderBookParams{Pair: pair, Precision: pricePrecision})
		}
	}
}

func (h *Hub) AddTradesTopic(pair string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers["trades"]; !exists {
		h.Subscribers["trades"] = make(map[any]*Subscribers)
	}

	h.Subscribers["trades"][TradesParams{Pair: pair}] = &Subscribers{Clients: map[*Client]bool{}}
}

func (h *Hub) RemoveTradesTopic(pair string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topic, exists := h.Subscribers["trades"]; exists {
		delete(topic, TradesParams{Pair: pair})
	}
}

func (h *Hub) AddTickerTopic(pair string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers["ticker"]; !exists {
		h.Subscribers["ticker"] = make(map[any]*Subscribers)
	}

	h.Subscribers["ticker"][TickerParams{Pair: pair}] = &Subscribers{Clients: map[*Client]bool{}}
}

func (h *Hub) RemoveTickerTopic(pair string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topic, exists := h.Subscribers["ticker"]; exists {
		delete(topic, TickerParams{Pair: pair})
	}
}

func (h *Hub) AddCandleStickTopic(pair string, timeframes []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers["candleStick"]; !exists {
		h.Subscribers["candleStick"] = make(map[any]*Subscribers)
	}

	for _, timeframe := range timeframes {
		h.Subscribers["candleStick"][CandleStickParams{Pair: pair, Timeframe: timeframe}] = &Subscribers{Clients: map[*Client]bool{}}
	}
}

func (h *Hub) RemoveCandleStickTopic(pair string, candleStickTimeframes []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topic, exists := h.Subscribers["candleStick"]; exists {
		for _, timeframe := range candleStickTimeframes {
			delete(topic, CandleStickParams{Pair: pair, Timeframe: timeframe})
		}
	}
}
