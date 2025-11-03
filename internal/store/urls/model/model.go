package model

import "time"

type UrlData struct {
	Key    string     `bson:"key_value"`
	Url    string     `bson:"url"`
	Expiry *time.Time `bson:"expiry,omitempty"`
}
