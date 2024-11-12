package crypto

import (
	"encoding/json"
	"errors"

	"github.com/sufu777/JumpProxy/crypto/sm2"
	"github.com/sufu777/JumpProxy/crypto/sm3"
	"github.com/sufu777/JumpProxy/crypto/sm4"
)

type Message struct {
	Iv            string        `json:"IV"`
	SymmetricKey  string        `json:"SymmetricKey"`
	MessageBody   string        `json:"MessageBody"`
	MessageHeader MessageHeader `json:"MessageHeader"`
	Sign          string        `json:"Sign"`
}

type MessageHeader struct {
	Sender       string `json:"Sender"`
	TxnTime      string `json:"TransactionTime"`
	TradeSource  string `json:"TradeSource"`
	Version      string `json:"Version"`
	Digest       string `json:"Digest"`
	BusinessCode string `json:"BusinessCode"`
	Recver       string `json:"Recver"`
	ReturnResult string `json:"ReturnResult"`
	TxnDate      string `json:"TransactionDate"`
}

type DecodedMessage struct {
	Iv            []byte
	SymmetricKey  []byte
	MessageHeader MessageHeader
	Sign          []byte
	MessageBody   map[string]string
}

func DecryptBankMessage(jsonMessage string) (*DecodedMessage, error) {
	var message Message
	err := json.Unmarshal([]byte(jsonMessage), &message)
	if err != nil {
		return nil, errors.New("error parse Json")
	}
	iv, err := sm2.Base64Decrypt(message.Iv)
	if err != nil {
		return nil, errors.New("error decrypt IV: " + err.Error())
	}
	symmetricKey, err := sm2.Base64Decrypt(message.SymmetricKey)
	if err != nil {
		return nil, errors.New("error decrypt symmetric key: " + err.Error())
	}
	bodyJson, err := sm4.Base64Decrypt(symmetricKey, iv, message.MessageBody)
	if err != nil {
		return nil, errors.New("error decrypt message: " + err.Error())
	}
	var req = make(map[string]string)
	err = json.Unmarshal(bodyJson, &req)
	if err != nil {
		return nil, errors.New("error unmarshal json: " + err.Error())
	}
	return &DecodedMessage{
		Iv:            iv,
		SymmetricKey:  symmetricKey,
		MessageHeader: message.MessageHeader,
		MessageBody:   req,
	}, nil
}

// / 使用请求报文相同的 key iv变量加密数据并返回
func EncodeMessage(req DecodedMessage, rspMap map[string]string) (Message, error) {
	// 1. rsp Json -> string
	bodyJson, err := json.Marshal(rspMap)
	if err != nil {
		return Message{}, errors.New("error marshal josn: " + err.Error())
	}
	iv := req.Iv
	symmetricKey := req.SymmetricKey
	rspJson, err := sm4.EncryptBase64(symmetricKey, iv, bodyJson)
	if err != nil {
		return Message{}, errors.New("error encrypt rsp body: " + err.Error())
	}
	ivString, err := sm2.EncryptBase64(iv)
	if err != nil {
		return Message{}, errors.New("error encrypt Iv: " + err.Error())
	}
	symmetricKeyString, err := sm2.EncryptBase64(symmetricKey)
	if err != nil {
		return Message{}, errors.New("error encrypt symmetricKey: " + err.Error())
	}
	reqHeader := req.MessageHeader
	var header = MessageHeader{
		Sender:       reqHeader.Recver,
		Recver:       reqHeader.Sender,
		TxnTime:      reqHeader.TxnTime,
		TxnDate:      reqHeader.TxnDate,
		TradeSource:  "0",
		Version:      reqHeader.Version,
		Digest:       reqHeader.Digest,
		BusinessCode: reqHeader.BusinessCode,
		ReturnResult: "0000",
	}
	var res string = ""
	// 计算签名
	res = res + "BusinessCode=" + header.BusinessCode + "&"
	res = res + "Digest=" + header.Digest + "&"
	res = res + "Recver=" + header.Recver + "&"
	res = res + "ReturnResult=" + header.ReturnResult + "&"
	res = res + "Sender=" + header.Sender + "&"
	res = res + "TradeSource=" + header.TradeSource + "&"
	res = res + "TransactionDate=" + header.TxnDate + "&"
	res = res + "TransactionTime=" + header.TxnTime + "&"
	res = res + "Version=" + header.Version + "&"
	res = res + "Data=" + rspJson + "&"
	res = res + "SymmetricKey=" + symmetricKeyString + "&"
	res = res + "IV=" + ivString
	hashStr := sm3.EncryptHexUpperCase([]byte(res))
	sign, err := sm2.SignBase64([]byte(hashStr))
	if err != nil {
		return Message{}, errors.New("error in sign rsp message: " + err.Error())
	}

	var rspMessage = Message{
		MessageBody:   rspJson,
		Iv:            ivString,
		SymmetricKey:  symmetricKeyString,
		MessageHeader: header,
		Sign:          sign,
	}

	return rspMessage, nil
}
