package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sufu777/JumpProxy/crypto"
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
	start := time.Now().UnixMilli()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error read request body from %s", r.RemoteAddr)
		return
	}
	reqBodyStr := string(reqBody)
	startIdx := strings.Index(reqBodyStr, "<arg0 xmlns=\"\">")
	endIdx := strings.Index(reqBodyStr, "</arg0>")
	reqJson := reqBodyStr[startIdx+14 : endIdx]
	println(reqJson)
	message, err := crypto.DecryptBankMessage(reqJson)
	if err != nil {
		fmt.Printf("error decrypt message: %s", err.Error())
		return
	}
	var rspMap = map[string]string{}
	rspMap["CODE"] = "00000000"
	rspMap["MESSAGE"] = "success"
	rspMap["AAC058"] = "01"
	rspMap["AAC147"] = message.MessageBody["AAC147"]
	rspMap["AAC003"] = message.MessageBody["AAC003"]
	rspMap["AAC002"] = ""
	rspMap["AAB301"] = "301122"
	rspMap["AIC674"] = "05"
	rspMap["AIC500"] = message.MessageBody["AAC147"] + "777P"
	rspMap["AIC501"] = ""
	rspMap["AIC657"] = ""
	rspMap["AIC539"] = ""
	rspMap["AIC509"] = "12000"
	rspMap["AAC341"] = message.MessageBody["AAZ341"]
	rspMap["AAZ345"] = NextAAZ345()

	dst, err := crypto.EncodeMessage(*message, rspMap)
	if err != nil {
		fmt.Printf("error encode message: " + err.Error())
		return
	}
	bytes, err := json.Marshal(dst)
	if err != nil {
		fmt.Printf("error marshal message to json: " + err.Error())
		return
	}
	end := time.Now().UnixMilli()
	if end-start > delay {
		time.Sleep(time.Duration(delay-end+start) * time.Millisecond)
	}
	w.Write(bytes)
}

var cc = "0000000"

// 9开头的18位流水号
func NextAAZ345() string {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	suffix := strconv.FormatInt(int64(id.NextId()), 10)
	return "9" + instanceId + ts[len(ts)-10:] + cc[:6-len(suffix)] + suffix
}
