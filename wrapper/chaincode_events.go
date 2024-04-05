package wrapper

import (
	"hlf-block-explorer/db_connector/models"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func Wrapper_chaincode_events(_chaincode_events *client.ChaincodeEvent) *models.Chaincode_Event {
	new_chaincode_event := &models.Chaincode_Event{}
	new_chaincode_event.Block.Block_number = int(_chaincode_events.BlockNumber)
	new_chaincode_event.Chaincode.Chaincode_id = _chaincode_events.ChaincodeName
	new_chaincode_event.Content = string(_chaincode_events.Payload)
	new_chaincode_event.Event_name = _chaincode_events.EventName
	new_chaincode_event.Transaction.Tx_id = _chaincode_events.TransactionID
	return new_chaincode_event
}
