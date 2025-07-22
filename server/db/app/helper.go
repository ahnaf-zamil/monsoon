package app

import (
	"context"
	"monsoon/api"
	"monsoon/db"
	"monsoon/db/tables"
)

func (u *UserDB) GetUserByID(c context.Context, userID string) (*api.UserModel, error) {
	fields := map[db.DBColumn]any{
		tables.ColUserID: userID,
	}
	user, err := u.GetUserByAnyField(c, fields)
	return user, err
}
