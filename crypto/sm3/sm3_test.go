package sm3

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	println(hex.EncodeToString(Encrypt([]byte("123123"))))
}

func TestEncrypt2(t *testing.T) {
	var x = []string{"Sender", "TransactionTime", "TradeSource", "Version", "Digest", "BusinessCode", "Recver", "ReturnResult", "TransactionDate"}
	sort.Strings(x)
	fmt.Printf("%#v", x)
}

func TestEncrypto(t *testing.T) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	println(ts[len(ts)-10:])
	
}
