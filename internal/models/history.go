package models

import "time"

type History struct {
	From      string    `json:"from" firestore:"from"`
	Message   string    `json:"message" firestore:"message"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}
