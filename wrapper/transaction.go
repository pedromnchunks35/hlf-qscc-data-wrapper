package wrapper

import (
	"fmt"
	"hlf-block-explorer/db_connector/models"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/msp"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

func Wrapper_transaction(_gateway *client.Gateway, _block_number int, _txs_ids []string, _channel string) ([]*models.Transaction, error) {
	channel := _gateway.GetNetwork(_channel)
	contract := channel.GetContract("qscc")
	new_transactions_object := []*models.Transaction{}
	for i := 0; i < len(_txs_ids); i++ {
		new_collect_transaction := &models.Transaction{}
		new_collect_transaction.Block_number = _block_number
		//? Evaluate the txId received
		res_bytes, err := contract.EvaluateTransaction("GetTransactionByID", _channel, _txs_ids[i])
		if err != nil {
			return nil, fmt.Errorf("Error evaluating the txId %v with error %v \n", _txs_ids[i], err)
		}
		//? Getting the main body
		main_body := &peer.ProcessedTransaction{}
		err = proto.Unmarshal(res_bytes, main_body)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshling the protobuf format object for the main_body for the txId %v with error %v \n", _txs_ids[i], err)
		}
		//? Getting the payload
		payload := &common.Payload{}
		err = proto.Unmarshal(main_body.TransactionEnvelope.Payload, payload)
		if err != nil {
			return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload for the txId %v with error %v \n", _txs_ids[i], err)
		}
		//! GETTING THE CREATOR
		//? Getting the Signature header
		signature_header := &common.SignatureHeader{}
		err = proto.Unmarshal(payload.Header.SignatureHeader, signature_header)
		if err != nil {
			return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Header->SignatureHeader for the txId %v with error %v \n", _txs_ids[i], err)
		}
		//? Getting the creator and adding the creator to the model
		creator := &msp.SerializedIdentity{}
		err = proto.Unmarshal(signature_header.Creator, creator)
		if err != nil {
			return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Header->SignatureHeader->Creator for the txId %v with error %v \n", _txs_ids[i], err)
		}
		new_collect_transaction.Creator = &models.Creator{Msp_id: creator.Mspid, Id_bytes: string(creator.IdBytes)}
		//! GETTING THE ENDORSEMENTS
		//? Getting payload data
		payload_data := &peer.Transaction{}
		err = proto.Unmarshal(payload.Data, payload_data)
		if err != nil {
			return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Data for the txId %v with error %v\n", _txs_ids[i], err)
		}
		//? Getting the actions on the transaction
		for j := 0; j < len(payload_data.Actions); j++ {
			single_action := &peer.ChaincodeActionPayload{}
			err = proto.Unmarshal(payload_data.Actions[j].Payload, single_action)
			if err != nil {
				return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Data->payload->actions[%v]->action from txid %v with error %v\n", j, _txs_ids[i], err)
			}
			//? Getting the endorsements
			new_collect_transaction.Endorsements = []*models.Endorsement{}
			for k := 0; k < len(single_action.Action.Endorsements); k++ {
				endorsement := &msp.SerializedIdentity{}
				err = proto.Unmarshal(single_action.Action.Endorsements[k].Endorser, endorsement)
				if err != nil {
					return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Data->payload->actions[%v]->action->Endorsements[%v] from txid %v with error %v \n", j, k, _txs_ids[i], err)
				}
				new_collect_transaction.Endorsements = append(new_collect_transaction.Endorsements, &models.Endorsement{
					Creator: &models.Creator{
						Msp_id:   endorsement.Mspid,
						Id_bytes: string(endorsement.IdBytes),
					},
				})
			}
			//! GETTING THE OPERATION
			//? Chaincode proposal payload
			single_action_chaincode_proposal_payload := &peer.ChaincodeProposalPayload{}
			err = proto.Unmarshal(single_action.ChaincodeProposalPayload, single_action_chaincode_proposal_payload)
			if err != nil {
				return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Data->payload->actions[%v]->action->chaincode_proposal_payload from txid %v with error %v\n", j, _txs_ids[i], err)
			}
			//? inputs for the chaincode proposal payload
			inputs_chaincode_proposal_payload := &peer.ChaincodeInvocationSpec{}
			err = proto.Unmarshal(single_action_chaincode_proposal_payload.Input, inputs_chaincode_proposal_payload)
			if err != nil {
				return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->Payload->Data->payload->actions[%v]->action->chaincode_proposal_payload->Input from txid %v with error %v \n", j, _txs_ids[i], err)
			}
			args := []*models.Operation_arg{}
			for l := 0; l < len(inputs_chaincode_proposal_payload.ChaincodeSpec.Input.Args); l++ {
				args = append(args, &models.Operation_arg{
					Argument: string(inputs_chaincode_proposal_payload.ChaincodeSpec.Input.Args[l]),
				})
			}
			new_collect_transaction.Operation = &models.Operation{
				Chaincode: &models.Chaincode{
					Chaincode_id: inputs_chaincode_proposal_payload.ChaincodeSpec.ChaincodeId.Name,
				},
				Chaincode_type: inputs_chaincode_proposal_payload.ChaincodeSpec.Type.String(),
				Operation_args: args,
			}
		}
		//! GETTING THE PAYLOAD
		new_collect_transaction.Payload = string(main_body.TransactionEnvelope.Payload)
		//! GETTING THE TIMESTAMP
		channel_header := &common.ChannelHeader{}
		err = proto.Unmarshal(payload.Header.ChannelHeader, channel_header)
		if err != nil {
			return nil, fmt.Errorf("Error getting the ProcessedTransaction->TransactionEnvelope->payload->Header->ChannelHeader from txId %v with error %v", _txs_ids[i], err)
		}
		new_collect_transaction.Timestamp = int(channel_header.Timestamp.Seconds)
		//! GETTING THE TX_ID
		new_collect_transaction.Tx_id = _txs_ids[i]
		//! GETTING THE TX_VALIDATION_TYPE
		new_collect_transaction.Tx_validation_type = &models.Tx_validation_type{
			Id: int(main_body.ValidationCode),
		}
		//! GETTING THE TYPE_TRANSACTION
		new_collect_transaction.Type_transaction = &models.Type_transaction{
			Id: int(channel_header.Type),
		}
		//? Add the object to the array of transactions
		new_transactions_object = append(new_transactions_object, new_collect_transaction)
	}
	return new_transactions_object, nil
}
