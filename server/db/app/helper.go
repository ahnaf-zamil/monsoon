package app

import (
	"context"
	"monsoon/api"
	"monsoon/db"
)

func (u *UserDB) GetUserByID(c context.Context, userID string) (*api.UserModel, error) {
	fields := map[db.UserColumn]any{
		db.ColUserID: userID,
	}
	user, err := u.GetUserByAnyField(c, fields)
	return user, err
}
