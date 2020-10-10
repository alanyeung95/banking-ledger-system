package transaction

type Transaction struct {
	ID       string `json:"id"    bson:"_id"`
	Operation     string `json:"operation"    bson:"operation"`
	Body Detail `json:"body"    bson:"body"`

}

type Detail struct {
	Amount       string `json:"amount"    bson:"amount"`
	From       string `json:"from"    bson:"from"`
	To       string `json:"to"    bson:"to"`
}
