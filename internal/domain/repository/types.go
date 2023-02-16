package repository

import "time"

type sqlDate []byte

func (d sqlDate) Time() time.Time {
	return time.Now().UTC()
}