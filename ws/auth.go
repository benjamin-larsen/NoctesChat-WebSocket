package ws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/benjamin-larsen/NoctesChat-WebSocket/database"
	"github.com/benjamin-larsen/NoctesChat-WebSocket/models"
)

var ErrLoggedOut = errors.New("Login: token is logged out")

func (s *Socket) ProcessLogin(msg models.LoginInbound) error {
	token, err := models.UserTokenFromString(msg.Token)

	if err != nil {
		s.Close(1008, "Invalid token")
		return err
	}

	tx, err := database.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	hasToken, err := database.HasUserToken(token, tx)

	if err != nil {
		s.RaiseException(err)
		return err
	}

	if hasToken != true {
		s.conn.WriteJSON(models.AuthErrLoggedOut)
		s.Close(1008, "You've been logged out. Please log in and try again.")
		return ErrLoggedOut
	}

	s.userId = token.UserId
	s.userToken = token

	// set subs here

	s.hasAuth = true

	rows, err := tx.Query(`SELECT
    cm.channel_id AS id,
    cm.last_accessed,
    c.name,
    c.created_at,
    o.id AS owner_id,
    o.username AS owner_username,
    o.created_at AS owner_created_at
FROM channel_members cm
JOIN channels c ON cm.channel_id = c.id
LEFT JOIN users o ON c.owner = o.id
WHERE cm.user_id = ?
FOR SHARE OF cm;`, s.userId)

	if err != nil {
		s.RaiseException(err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		channel, err := models.ChannelFromRows(rows)

		if err != nil {
			s.RaiseException(err)
			return err
		}

		json, _ := json.MarshalIndent(channel, "", "\t")

		fmt.Println(string(json))
	}
	
	return nil
}