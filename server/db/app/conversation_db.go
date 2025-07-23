package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"monsoon/api"
	"monsoon/db"
	"monsoon/db/tables"
	"time"

	"github.com/jackc/pgx/v5"
)

type ConversationDB struct {
	AppDB *db.AppDB
}

type IConversationDB interface {
	CreateUserDM(ctx context.Context, conversationID int64, user1ID string, user2ID string) error
	GetExistingDM(ctx context.Context, user1ID string, user2ID string) (*api.DirectConversationModel, error)
	GetConvesationParticipantsByFields(ctx context.Context, fields map[db.DBColumn]any) ([]api.ConversationParticipant, error)
	GetUserInboxConversations(ctx context.Context, userID string) ([]api.InboxConversation, error)
}

func GetConversationDB() *ConversationDB {
	convo_db := &ConversationDB{AppDB: db.GetAppDB()}
	return convo_db
}

func (cv *ConversationDB) CreateUserDM(ctx context.Context, conversationID int64, user1ID string, user2ID string) error {
	/* Service function to create a user
	 */

	tx, err := cv.AppDB.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Println("TX rollback error:", err)
		}
	}()

	createdAt := time.Now().Unix()
	updatedAt := createdAt
	joinedAt := createdAt

	// Create conversation entry
	err = insertConversation(tx, ctx, conversationID, db.CONVERSATION_DM, createdAt, updatedAt)
	if err != nil {
		return err
	}

	// Add participants
	err = insertConversationParticipant(tx, ctx, conversationID, user1ID, joinedAt, db.CONVERSATION_ROLE_MEMBER)
	if err != nil {
		return err
	}
	err = insertConversationParticipant(tx, ctx, conversationID, user2ID, joinedAt, db.CONVERSATION_ROLE_MEMBER)
	if err != nil {
		return err
	}

	// Create DM conversation entry with the two participants
	err = insertDMConversation(tx, ctx, conversationID, user1ID, user2ID, createdAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("TX commit error:", err)
	}

	return nil
}

func (cv *ConversationDB) GetExistingDM(ctx context.Context, user1ID string, user2ID string) (*api.DirectConversationModel, error) {
	cols := []db.DBColumn{tables.ColDMID, tables.ColDMUser1, tables.ColDMUser2, tables.ColDMCreatedAt}
	selected_columns := db.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE user1 = LEAST($1::BIGINT, $2::BIGINT) AND user2 = GREATEST($1::BIGINT, $2::BIGINT)", selected_columns, tables.TableDirectConversations)

	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	row := cv.AppDB.DBPool.QueryRow(ctx, query, user1ID, user2ID)

	var dm api.DirectConversationModel
	err := row.Scan(&dm.ConversationID, &dm.User1ID, &dm.User2ID, &dm.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		// Return nil if no rows
		return nil, nil
	} else if err != nil {
		// If there's error, just return error
		return nil, err
	} else {
		// Else, just return model
		return &dm, nil
	}
}


func (cv *ConversationDB) GetConvesationParticipantsByFields(ctx context.Context, fields map[db.DBColumn]any) ([]api.ConversationParticipant, error) {
	field_arr := []db.DBColumn{}
	value_arr := []any{}
	for k, v := range fields {
		field_arr = append(field_arr, k)
		value_arr = append(value_arr, v)
	}

	// Generate the "OR" conditions using given fields
	or_fields := db.GenerateDBOrFields(field_arr)

	// Generate the SELECT columns
	cols := []db.DBColumn{tables.ColConvoPartConversationID, tables.ColConvoPartUserID, tables.ColConvoPartJoinedAt, tables.ColConvoPartRole}
	selected_columns := db.GenerateDBQueryFields(cols)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		selected_columns,
		tables.TableConversationParticipants,
		or_fields)

	// The value_arr maintains same sequence of parameters as the columns, which is why we separated the map into two slices
	rows, err := cv.AppDB.DBPool.Query(ctx, query, value_arr...)
	if err != nil {
		return nil, err
	}

	conversations, err := pgx.CollectRows(rows, pgx.RowToStructByName[api.ConversationParticipant])
	if err != nil {
		return nil, err
	}

	return conversations, err
}

func (cv *ConversationDB) GetUserInboxConversations(ctx context.Context, userID string) ([]api.InboxConversation, error) {
	query := `
(
select
	c.id as conversation_id,
	c.type as type,
	c.group_name as name,
	c.updated_at as updated_at,
	null as user_id
from
	conversations c
join conversation_participants cp 
on
	cp.conversation_id = c.id
where
	c."type" = 'GROUP'
	and cp.user_id = $1)
union 
(
select 
	dc.conversation_id as conversation_id, 
	c.type as type, 
	u.display_name as name,
	c.updated_at as updated_at,
	u.id as user_id
from
	direct_conversations dc
join conversations c on
	dc.conversation_id = c.id
join users u on
(
	case
		when dc.user1 = $1 then dc.user2
	else
		dc.user1
end
) = u.id
where
	dc.user1 = $1
	or dc.user2 = $1
)
order by updated_at desc;`
	rows, err := cv.AppDB.DBPool.Query(ctx, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	conversations, err := pgx.CollectRows(rows, pgx.RowToStructByName[api.InboxConversation])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conversations, err
}


func insertConversation(tx pgx.Tx, ctx context.Context, conversationID int64, convo_type db.ConversationType, createdAt int64, updatedAt int64) error {
	query := fmt.Sprintf("INSERT INTO %s (id, type, created_at, updated_at) VALUES ($1, $2, $3, $4)", tables.TableConversations)
	_, err := tx.Exec(ctx, query, conversationID, convo_type, createdAt, updatedAt)
	return err
}

func insertConversationParticipant(tx pgx.Tx, ctx context.Context, conversationID int64, userID string, joinedAt int64, role db.ConversationRole) error {
	query := fmt.Sprintf("INSERT INTO %s (conversation_id, user_id, joined_at, role) VALUES ($1, $2, $3, $4)", tables.TableConversationParticipants)
	_, err := tx.Exec(ctx, query, conversationID, userID, joinedAt, role)
	return err
}

func insertDMConversation(tx pgx.Tx, ctx context.Context, conversationID int64, user1ID string, user2ID string, createdAt int64) error {
	query := fmt.Sprintf("INSERT INTO %s (conversation_id, user1, user2, created_at) VALUES ($1, LEAST($2::BIGINT, $3::BIGINT), GREATEST($2::BIGINT, $3::BIGINT), $4)", tables.TableDirectConversations)
	_, err := tx.Exec(ctx, query, conversationID, user1ID, user2ID, createdAt)
	return err
}
