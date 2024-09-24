package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

type Hub struct {
	Users map[int64]*User
	mu    sync.RWMutex
}

type User struct {
	ID      int64
	Clients map[*Client]bool
}

type Client struct {
	Conn   *websocket.Conn
	Topics map[string]bool
	mu     sync.Mutex
}

func NewHub(topics []string) *Hub {
	h := &Hub{
		Users: make(map[int64]*User),
	}

	return h
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

	//TODO
	//check if user exists

	conn, err := websocket.Accept(c.Response(), c.Request(), nil)
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
		Topics: make(map[string]bool),
		mu:     sync.Mutex{},
	}

	user.Clients[client] = true
	h.mu.Unlock()

	go h.readPump(client, userID)

	return nil
}

func (h *Hub) readPump(c *Client, userID int64) {
	defer func() {
		c.Conn.Close(websocket.StatusNormalClosure, "normal closure")
		//delete client and user if no clients left
		h.mu.Lock()
		delete(h.Users[userID].Clients, c)
		if len(h.Users[userID].Clients) == 0 {
			delete(h.Users, userID)
		}
		h.mu.Unlock()
	}()

	for {
		_, msg, err := c.Conn.Read(context.Background())
		if err != nil {
			return
		}

		var request models.SubscriptionRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			continue
		}

		switch request.Action {
		case "subscribe":
			switch request.Topic {
			case "orderUpdate":
				c.mu.Lock()
				c.Topics["orderUpdate"] = true
				c.mu.Unlock()
			}

		case "unsubscribe":
			switch request.Topic {
			case "orderUpdate":
				c.mu.Lock()
				delete(c.Topics, "orderUpdate")
				c.mu.Unlock()
			}
		}
	}
}

func (h *Hub) BroadcastUpdatedOrder(order models.Order) {
	h.mu.RLock()
	user, ok := h.Users[order.UserID]
	h.mu.RUnlock()

	if ok {
		for client := range user.Clients {
			client.mu.Lock()
			_, ok := client.Topics["orderUpdate"]
			client.mu.Unlock()

			if ok {
				go func(c *Client) {
					orderJSON, err := json.Marshal(order)
					if err != nil {
						return
					}

					c.mu.Lock()
					defer c.mu.Unlock()

					err = c.Conn.Write(context.Background(), websocket.MessageText, orderJSON)
					if err != nil {
						return
					}

				}(client)
			}
		}
	}
}
