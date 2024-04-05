package wrapper

import (
	"hlf-block-explorer/db_connector/models"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
)

func Wrap_Block_Event(_block *common.Block, _network string, _channel string) *models.Block {
	new_block := &models.Block{}
	new_block.Network = &models.Network{Ip: _network}
	new_block.Channel = &models.Channel{Channel_name: _channel}
	new_block.Block_number = int(_block.Header.Number)
	new_block.Next_hash = string(_block.Header.DataHash)
	new_block.Prev_hash = string(_block.Header.PreviousHash)
	return new_block
}
