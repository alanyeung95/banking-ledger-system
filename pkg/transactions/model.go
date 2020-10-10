package transactions

import "time"

type operationType string

const (
	Deposit       operationType = "deposit"
	Withdraw operationType = "withdraw"
	Transfer operationType = "transfer"
	Undo operationType = "undo"
)

type Transaction struct {
	ID       string `json:"id"    bson:"_id"`
	Time time.Time     `json:"time"    bson:"time"`
	Operation     operationType `json:"operation"    bson:"operation"`
	Body Detail `json:"body"    bson:"body"`
	Notes string `json:"notes"    bson:"notes"`
}

type Detail struct {
	Amount       int `json:"amount"    bson:"amount"`
	From       string `json:"from"    bson:"from"`
	To       string `json:"to"    bson:"to"`
}
