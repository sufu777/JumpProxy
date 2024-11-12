package sm3

import (
	"encoding/hex"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

var salt = "e44c1712cd8742aaa3cb"
var digestCount = 2

// sm3 带盐值计算两轮hash
func Encrypt(data []byte) []byte {
	innerSm3 := sm3.New()
	innerSm3.Write([]byte(salt))
	term := innerSm3.Sum(data)
	innerSm3.Reset()
	return innerSm3.Sum(term)
}

func EncryptHexUpperCase(data []byte) string {
	sum := Encrypt(data)
	sumHex := hex.EncodeToString(sum)
	return strings.ToUpper(sumHex)
}
