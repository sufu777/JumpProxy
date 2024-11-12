package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/sufu777/JumpProxy/id"
)

var instanceId string
var delay int64

func main() {
	flag.StringVar(&instanceId, "i", strconv.FormatInt(int64(rand.Int())%10, 10), "指定当前实例id")
	flag.Int64Var(&delay, "d", 150, "指定延迟")
	flag.Parse()
	println("instance id " + instanceId)
	http.HandleFunc("/dszz/services/ThirdPillarService", handler)
	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		fmt.Printf("error start server on 8091,%s", err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message = Message{}
	err := decoder.Decode(&message)
	if err != nil {
		return
	}
	reqHeader := message.MessageHeader
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
	reqMap := message.MessageBody
	var rspMap = map[string]string{}
	rspMap["CODE"] = "00000000"
	rspMap["MESSAGE"] = "success"
	rspMap["AAC058"] = "01"
	rspMap["AAC147"] = reqMap["AAC147"]
	rspMap["AAC003"] = reqMap["AAC003"]
	rspMap["AAC002"] = reqMap["AAC002"]
	rspMap["AAB301"] = "301122"
	rspMap["AIC674"] = "05"
	rspMap["AIC500"] = reqMap["AAC147"] + "777P"
	rspMap["AIC501"] = ""
	rspMap["AIC657"] = ""
	rspMap["AIC539"] = ""
	rspMap["AIC509"] = "12000"
	rspMap["AAC341"] = reqMap["AAZ341"]
	rspMap["AAZ345"] = NextAAZ345()
	time.Sleep(time.Duration(delay) * time.Millisecond)
	var response = map[string]any{}
	response["MessageHeader"] = header
	response["MessageBody"] = rspMap
	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(marshal)
}

var cc = "0000000"

// NextAAZ345 9开头的18位流水号
func NextAAZ345() string {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	suffix := strconv.FormatInt(id.NextId(), 10)
	return "9" + instanceId + ts[len(ts)-10:] + cc[:6-len(suffix)] + suffix
}

type Message struct {
	MessageBody   map[string]string `json:"MessageBody"`
	MessageHeader MessageHeader     `json:"MessageHeader"`
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
