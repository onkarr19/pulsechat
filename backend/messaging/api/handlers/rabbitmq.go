package handlers

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var (
	rabbitMQURL   = "amqp://guest:guest@localhost:5672/"
	exchangeName  = "chat_rooms"
	roomMutex     sync.Mutex
	RoomProcesses = make(map[string]*RoomProcess)
)

// RoomProcess struct represents a goroutine managing a room
type RoomProcess struct {
	room        string
	connections map[*websocket.Conn]bool
	stopCh      chan struct{}
}

func NewRoomProcess(room string) *RoomProcess {
	return &RoomProcess{
		room:        room,
		connections: make(map[*websocket.Conn]bool),
		stopCh:      make(chan struct{}),
	}
}

func (rp *RoomProcess) Start() {
	go func() {
		conn, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ for room %s: %s", rp.room, err)
			close(rp.stopCh)
			return
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel for room %s: %s", rp.room, err)
			close(rp.stopCh)
			return
		}
		defer ch.Close()

		err = ch.ExchangeDeclare(
			exchangeName,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to declare exchange for room %s: %s", rp.room, err)
			close(rp.stopCh)
			return
		}

		for {
			select {
			case <-rp.stopCh:
				return
			default:
				if err := rp.ConsumeFromRabbitMQ(ch); err != nil {
					log.Printf("Error consuming message from RabbitMQ for room %s: %s", rp.room, err)
				}
			}
		}
	}()
}

func (rp *RoomProcess) ConsumeFromRabbitMQ(ch *amqp.Channel) error {
	msgs, err := ch.Consume(
		rp.room,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages for room %s: %s", rp.room, err)
	}

	for msg := range msgs {
		rp.BroadcastMessage(msg.Body)
	}

	return nil
}

func (rp *RoomProcess) BroadcastMessage(message []byte) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	for conn := range rp.connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to connection in room %s: %s", rp.room, err)
		}
	}
}

func (rp *RoomProcess) AddConnection(conn *websocket.Conn) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	rp.connections[conn] = true
}

func (rp *RoomProcess) RemoveConnection(conn *websocket.Conn) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	delete(rp.connections, conn)
}

func PublishToRabbitMQ(room string, message []byte) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println("Failed to declare exchange:", err)
		return
	}

	err = ch.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		log.Println("Failed to publish message to RabbitMQ:", err)
	}
	fmt.Println("Published message to RabbitMQ:", string(message))
}
