package domain

import "time"

type AuthKey struct {
	Id      string    `json:"id"`
	Ttl     int       `json:"ttl"`
	Created time.Time `json:"created"`
	UserId  int       `json:"userId"`
	CId     int64     `json:"cId"`
	MId     int       `json:"mId"`
}
