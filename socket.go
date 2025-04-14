package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Socket struct {
	ID     string
	Rooms  map[string]bool
	UserID string // Optional, socket will get it after they authenticate
	WsConn *websocket.Conn
}

var (
	socketList []*Socket
	roomState  sync.Map
)

/* Socket operations */

func PrintSocketList() {
	pretty, err := json.MarshalIndent(socketList, "", "")
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	fmt.Println(string(pretty))
}

func RemoveSocketFromList(client_s *Socket) {
	newSList := []*Socket{}

	for _, s := range socketList {
		if s.ID != client_s.ID {
			newSList = append(newSList, s)
		}
	}

	socketList = newSList
}

func GetRoomState() *sync.Map {
	return &roomState
}

func AddSocketToRoom(s *Socket, roomId string) {
	// Set room ID inside socket struct for later listing
	(*s).Rooms[roomId] = true

	// Add socket to room state
	// Load socket slice from map for channel Id
	val, _ := roomState.LoadOrStore(roomId, []*Socket{})
	// Append new socket to the current list of sockets
	val = append(val.([]*Socket), s)
	// Update state map
	roomState.Store(roomId, val)
}

func RemoveSocketFromRoom(s *Socket, roomId string) {
	// Remove room Id from socket struct
	delete((*s).Rooms, roomId)

	// Remove socket from room state map
	roomState.Range(func(k, v any) bool {
		// Loop over each room in the room state and work on the one with the matching room ID
		if k.(string) == roomId {
			sList := v.([]*Socket)

			// Initialize new socket List
			newSList := []*Socket{}

			// Loop over all existing sockets in the room and add them to new list except the one which we are "removing"
			for _, _s := range sList {
				// If socket ID does not match the "removed" socket's ID, then add it to new socket list
				if _s.ID != s.ID {
					newSList = append(newSList, _s)
				}
			}

			// If the length of new socket list and old socket list are not same (i.e there have been changes), then store new in the map
			if len(sList) != len(newSList) {
				roomState.Store(roomId, newSList)
			}
		}
		return true
	})
}
