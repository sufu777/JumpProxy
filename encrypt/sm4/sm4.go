package sm4

import (
	"bytes"
	"errors"
	"github.com/tjfoc/gmsm/sm4"
	"strconv"
)

// Decrypt pkcs5填充 cbc模式的sm4
func Decrypt(key []byte, iv []byte, data []byte) (out []byte, err error) {
	if len(key) != sm4.BlockSize {
		return nil, errors.New("SM4: invalid key size " + strconv.Itoa(len(key)))
	}
	out = make([]byte, len(data))
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(data)/16; i++ {
		inTmp := data[i*16 : i*16+16]
		outTmp := make([]byte, 16)
		cipher.Decrypt(outTmp, inTmp)
		outTmp = xor(outTmp, iv)
		copy(out[i*16:i*16+16], outTmp)
		iv = inTmp
	}
	return pkcs5UnPadding(out), nil
}

func Encrypt(key []byte, iv []byte, data []byte) (out []byte, err error) {
	if len(key) != sm4.BlockSize {
		return nil, errors.New("SM4: invalid key size " + strconv.Itoa(len(key)))
	}
	inData := pkcs5Padding(data, sm4.BlockSize)
	out = make([]byte, len(inData))
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(inData)/16; i++ {
		inTmp := xor(inData[i*16:i*16+16], iv)
		outTmp := make([]byte, 16)
		cipher.Encrypt(outTmp, inTmp)
		copy(out[i*16:i*16+16], outTmp)
		iv = outTmp
	}
	return out, nil
}

func xor(in, iv []byte) (out []byte) {
	if len(in) != len(iv) {
		return nil
	}

	out = make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[i] ^ iv[i]
	}
	return
}

// pkcs5Padding pkcs5Padding
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// pkcs5UnPadding pkcs5UnPadding
func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unPadding := int(src[length-1])
	return src[0 : length-unPadding]
}
