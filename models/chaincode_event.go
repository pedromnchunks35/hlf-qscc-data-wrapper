package models

type Chaincode_Event struct {
	Id          int         `json:"id"`
	Block       Block       `json:"block"`
	Chaincode   Chaincode   `json:"chaincode"`
	Transaction Transaction `json:"transaction"`
	Event_name  string      `json:"event_name"`
	Content     string      `json:"content"`
}
