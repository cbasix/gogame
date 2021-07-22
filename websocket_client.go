// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Sending ping messages to client with this interval
	pingInterval = 50 * time.Second
	// Time to wait for the next pong
	pongDeadline = 60 * time.Second
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*type InboundMessage struct {
	message []byte
	client  *Client
}*/

type Client struct {
	coordinator Coordinator

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *PlayerView

	player int
}

func (c *Client) forPlayer() int              { return c.player }
func (c *Client) channel() chan<- *PlayerView { return c.send }

func (c *Client) readTask() {
	defer func() {
		c.coordinator.unregister() <- c
		c.conn.Close()
	}()

	// In fact this first push activates the read deadline
	pushReadDeadline(c)
	c.conn.SetPongHandler(func(string) error {
		pushReadDeadline(c)
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("websocket error: %v", err)
			break
		}
		fmt.Printf("got websocket message. ignoring it")
	}
}

func pushReadDeadline(c *Client) {
	c.conn.SetReadDeadline(
		time.Now().Add(pongDeadline),
	)
}

func pushWriteDeadline(c *Client) {
	c.conn.SetWriteDeadline(
		time.Now().Add(writeWait),
	)
}

func (c *Client) writeTask() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case playerView, ok := <-c.send:
			pushWriteDeadline(c)
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			json.NewEncoder(w).Encode(playerView)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			pushWriteDeadline(c)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(coordinator Coordinator, w http.ResponseWriter, r *http.Request) {
	allowAllOrigins := func(r *http.Request) bool { return true }
	upgrader.CheckOrigin = allowAllOrigins

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{coordinator: coordinator, conn: conn, send: make(chan *PlayerView, 256)}
	client.coordinator.register() <- client

	go client.writeTask()
	go client.readTask()
}
