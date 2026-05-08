package chat

import "fmt"

type Hub struct {
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

type Message struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.RegisterNewClient(client)
		case client := <-h.unregister:
			h.RemoveClient(client)
		case message := <-h.broadcast:
			h.HandleMessage(message)
		}
	}
}

func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.ID] = connections
	}
	h.clients[client.ID][client] = true

	fmt.Println("Size of clients: ", len(h.clients[client.ID]))
}

func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients[client.ID], client)
		close(client.Send)
		fmt.Println("Removed client")
	}
}

func (h *Hub) HandleMessage(message Message) {
	if message.Type == "message" {
		clients := h.clients[message.ID]
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients[message.ID], client)
			}
		}
	}

	if message.Type == "notification" {
		fmt.Printf("Notification: %s\n", message.Content)
		clients := h.clients[message.Recipient]
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients[message.Recipient], client)
			}
		}
	}
}
