package ws

import (
	"context"
	"log"
)

func (w *WebSocketHandler) DispatchAndSyncSocketRooms(s *Socket) error {
	ctx := context.TODO()
	conversations, err := w.ConversationDB.GetUserInboxConversations(ctx, s.UserID)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, convo := range conversations {
		AddSocketToRoom(s, convo.ConversationID)
	}

	w.DispatchEvent(s, OpRoomSync, conversations)

	return nil
}
