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
	Symbol    string
	Precision int32
}

type TradesParams struct {
	Symbol string
}

type Subscribers struct {
	Clients map[*Client]bool
}

type Hub struct {
	Users       map[int64]*User
	mu          sync.RWMutex
	Subscribers map[string]map[any]*Subscribers
	logger      *slog.Logger
}

type User struct {
	ID      int64
	Clients map[*Client]bool
}

type Client struct {
	Conn   *websocket.Conn
	Topics map[string]any
	mu     sync.RWMutex
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		Users:       make(map[int64]*User),
		Subscribers: make(map[string]map[any]*Subscribers),
		logger:      logger,
		mu:          sync.RWMutex{},
	}
}

func (h *Hub) HandleWebsocket(c echo.Context) error {
	intUserID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userID",
		})
	}

	userID := int64(intUserID)

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

func (h *Hub) readPump(c *Client, userID int64) {
	defer func() {
		c.Conn.Close(websocket.StatusNormalClosure, "normal closure")

		h.mu.Lock()
		defer h.mu.Unlock()
		//delete subscriber, if no subscribers left - delete params, if no params left - delete topic
		for topic, params := range c.Topics {
			if subscribers, exists := h.Subscribers[topic][params]; exists {
				delete(subscribers.Clients, c)
				if len(subscribers.Clients) == 0 {
					delete(h.Subscribers[topic], params)
					if len(h.Subscribers[topic]) == 0 {
						delete(h.Subscribers, topic)
					}
				}
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
		if request.Params.Symbol == "" || request.Params.Precision == 0 {
			h.logger.Error("invalid parameters for orderBook subscription")
			return nil, false
		}
		return OrderBookParams{
			Symbol:    request.Params.Symbol,
			Precision: request.Params.Precision,
		}, true

	case "trades":
		if request.Params.Symbol == "" {
			h.logger.Error("invalid parameters for trades subscription")
			return nil, false
		}
		return TradesParams{
			Symbol: request.Params.Symbol,
		}, true

	default:
		h.logger.Error("attempt to subscribe to a non-existent topic")
		return nil, false
	}
}

func (h *Hub) subscribeClient(c *Client, topic string, params any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Subscribers[topic]; !ok {
		h.Subscribers[topic] = make(map[any]*Subscribers)
	}

	if _, exists := h.Subscribers[topic][params]; !exists {
		h.Subscribers[topic][params] = &Subscribers{
			Clients: make(map[*Client]bool),
		}
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
