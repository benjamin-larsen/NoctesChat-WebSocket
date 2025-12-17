package models

import "database/sql"

type Channel struct {
	ID        uint64 `json:"id,string"`
	Name      string `json:"name"`
	Owner     User   `json:"owner"`
	Members   []User `json:"members"`
	CreatedAt int64  `json:"created_at"`
}

func ChannelFromRows(rows *sql.Rows) (Channel, error) {
	channel := Channel{}

	var skipCol interface{};

	err := rows.Scan(
		&channel.ID,
		&skipCol,
		&channel.Name,
		&channel.CreatedAt,
		&channel.Owner.ID,
		&channel.Owner.Username,
		&channel.Owner.CreatedAt,
	)

	return channel, err
}