package DataStructures

import "time"

type Tuple struct {
	// element
	// TTL
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func NewTuple(value interface{}, ttl time.Duration) *Tuple {
	return &Tuple{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}
