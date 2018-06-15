package utils

import "DATx/utils/common"

type TransactionHeader struct {
	//transcation expiration ,the seconds from 1970.
	Expiration uint64

	//reference the latest block number
	RefBlockNum uint16

	//reference block prefix
	RefBlockPerfix uint32
}

type Transaction struct {
	//transaction header, inherited from TransactionHeader
	TransactionHeader

	//action list
	Actions []*Action

	//transaction hash
	TransactionHash common.Hash
}

//transaction constructor
func NewTrx(time uint64) *Transaction {
	return &Transaction{
		TransactionHeader: TransactionHeader{
			Expiration: time,
		},
	}
}