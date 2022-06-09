package entity

import (
	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `table:"posts"`
	Id            int    `json:"id"      bun:"id"`
	UserId        int    `json:"user_id" bun:"user_id"`
	Title         string `json:"title"   bun:"title"`
	Body          string `json:"body"    bun:"body"`
}
