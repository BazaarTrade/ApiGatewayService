package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

type OrderBookParams struct {
	Pair      string
	Precision int32
}

type TradesParams struct {
	Pair string
}

type TickerParams struct {
	Pair string
}

type Subscribers struct {
	Clients map[*Client]bool
}

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
	h.mu.Unlock()

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
			h.logger.Error("failed to read websocket connection", "error", err)
			return
		}

		var request models.SubscriptionRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			h.logger.Error("failed unmarshal SubscriptionRequest", "error", err)
			continue
		}

		params, valid := h.validateParams(request)
		if !valid {
			messageJSON, err := json.Marshal(map[string]string{
				"error": "no such topic exists",
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

func (h *Hub) validateParams(request models.SubscriptionRequest) (any, bool) {
	switch request.Topic {
	case "orderBook":
		if request.Params.Pair == "" || request.Params.Precision == 0 {
			h.logger.Error("invalid parameters for orderBook subscription")
			return nil, false
		}
		return OrderBookParams{
			Pair:      request.Params.Pair,
			Precision: request.Params.Precision,
		}, true

	case "trades":
		if request.Params.Pair == "" {
			h.logger.Error("invalid parameters for trades subscription")
			return nil, false
		}
		return TradesParams{
			Pair: request.Params.Pair,
		}, true

	case "ticker":
		if request.Params.Pair == "" {
			h.logger.Error("invalid parameters for ticker subscription")
			return nil, false
		}
		return TickerParams{
			Pair: request.Params.Pair,
		}, true
	default:
		h.logger.Error("attempt to subscribe to a non-existent topic")
		return nil, false
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

	if _, exists := h.Subscribers[topic][params]; exists {
		delete(h.Subscribers[topic][params].Clients, c)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Topics, topic)
}

func (h *Hub) AddOrderBookUpdateTopic(pair string, pricePrecisions []int32) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Subscribers["orderBook"]; !exists {
		h.Subscribers["orderBook"] = make(map[any]*Subscribers)
	}

	for _, pricePrecision := range pricePrecisions {
		h.Subscribers["orderBook"][OrderBookParams{Pair: pair, Precision: pricePrecision}] = &Subscribers{Clients: map[*Client]bool{}}
	}
}

func (h *Hub) RemoveOrderBookUpdateTopic(pair string, pricePrecisions []int32) {
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
