package models

type BlockChainError struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg"`
}

var NoError = BlockChainError{Code: 00000, Msg: "no error"}

var NotFoundError = BlockChainError{Code: 20401, Msg: "data not found"}

var BadRequestError = BlockChainError{Code: 40001, Msg: "Bad request"}

var InternalServerError = BlockChainError{Code: 50001, Msg: "Internal server error"}
