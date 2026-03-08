package model

import "time"

type URL struct {
	Addr          string
	RedirectCount int64
	CreationDate  time.Time
}
