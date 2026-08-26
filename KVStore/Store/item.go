package store

import "time"

type item struct {
	value any
	expiration int64
}

func (i item) isExpired() bool {
	if i.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.expiration
}

