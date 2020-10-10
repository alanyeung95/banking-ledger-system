package transactions

type operationType string

const (
	Deposit       operationType = "deposit"
	Withdraw operationType = "withdraw"
	Transfer operationType = "transfer"
)

type Transaction struct {
	ID       string `json:"id"    bson:"_id"`
	UserID string `json:"userID"    bson:"userID"`
	Operation     operationType `json:"operation"    bson:"operation"`
	Body Detail `json:"body"    bson:"body"`
}

type Detail struct {
	Amount       int `json:"amount"    bson:"amount"`
	From       string `json:"from"    bson:"from"`
	To       string `json:"to"    bson:"to"`
}
