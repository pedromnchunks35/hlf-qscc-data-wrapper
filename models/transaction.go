package models

type Type_transaction struct {
	Id               int    `json:"id"`
	Description_type string `json:"description_type"`
}
type Tx_validation_type struct {
	Id               int    `json:"id"`
	Description_type string `json:"description_type"`
}
type Creator struct {
	Id       int    `json:"id"`
	Msp_id   string `json:"msp_id"`
	Id_bytes string `json:"id_bytes"`
}
type Endorsement struct {
	Id      int      `json:"id"`
	Creator *Creator `json:"creator"`
}
type Operation_arg struct {
	Id       int    `json:"id"`
	Argument string `json:"argument"`
}
type Chaincode struct {
	Id           int    `json:"id"`
	Chaincode_id string `json:"chaincode_id"`
}
type Operation struct {
	Id             int              `json:"id"`
	Chaincode      *Chaincode       `json:"chaincode"`
	Chaincode_type string           `json:"chaincode_type"`
	Operation_args []*Operation_arg `json:"operation_args"`
}
type Transaction struct {
	Tx_id              string              `json:"tx_id"`
	Block_number       int                 `json:"block_number"`
	Timestamp          int                 `json:"timestamp"`
	Type_transaction   *Type_transaction   `json:"type_transaction"`
	Tx_validation_type *Tx_validation_type `json:"tx_validation_type"`
	Creator            *Creator            `json:"creator"`
	Operation          *Operation          `json:"operation"`
	Payload            string              `json:"payload"`
	Endorsements       []*Endorsement      `json:"endorsements"`
}
