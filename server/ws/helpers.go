package ws

import (
	"context"
	"log"
)

func (w *WebSocketHandler) DispatchAndSyncSocketRooms(s *Socket) error {
	ctx := context.TODO()
	conversations, err := w.ConversationDB.GetUserInboxConversations(ctx, s.UserID)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	for _, convo := range conversations {
		AddSocketToRoom(s, convo.ConversationID)
	}

	err = w.DispatchEvent(s, OpRoomSync, conversations)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
