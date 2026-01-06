package model

import "time"

type UrlData struct {
	Key    string     `bson:"short"`
	Url    string     `bson:"long"`
	Expiry *time.Time `bson:"expiry,omitempty"`
}
