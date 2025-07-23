package msg

import (
	"context"
	"fmt"
	"log"
	"monsoon/db"
	"monsoon/db/tables"
	"time"

	"github.com/jackc/pgx/v5"
)

type MessageDB struct {
	MsgDB *db.MsgDB
}

type IMessageDB interface {
	CreateMessage(ctx context.Context, messageID int64, conversationID string, authorID string, content any) error
}

func GetMessageDB() IMessageDB {
	msg_db := &MessageDB{MsgDB: db.GetMsgDB()}
	return msg_db
}

func (m *MessageDB) CreateMessage(ctx context.Context, messageID int64, conversationID string, authorID string, content any) error {
	tx, err := m.MsgDB.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Println("TX rollback error:", err)
		}
	}()

	createdAt := time.Now().Unix()

	err = insertMessage(tx, ctx, messageID, conversationID, authorID, content, createdAt, false)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("TX commit error:", err)
	}

	return nil
}

func insertMessage(tx pgx.Tx, ctx context.Context, messageID int64, conversationID string, authorID string, content any, createdAt int64, deleted bool) error {
	query := fmt.Sprintf("INSERT INTO %s (id, conversation_id, author_id, content, created_at, deleted) VALUES ($1, $2, $3, $4, $5, $6)", tables.TableMessages)
	_, err := tx.Exec(ctx, query, messageID, conversationID, authorID, content, createdAt, deleted)
	return err
}
