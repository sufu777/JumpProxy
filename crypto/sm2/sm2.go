package sm2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

func init() {
	SetupKeys()
}

var bankPublicKey *sm2.PublicKey
var platformPrivateKey *sm2.PrivateKey

var PlatformPrivateKeyNotExisted = errors.New("platform public key not existed: ./.platform_private")
var BankPublicNotExisted = errors.New("platform public key not existed: ./.bank_public")

// SetupKeys 公钥加密 私钥解密
func SetupKeys() error {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pltPubPath := filepath.Join(dir, ".platform_public")
	platformPrivateKeyBytes, err := os.ReadFile(pltPubPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return PlatformPrivateKeyNotExisted
		}
		return err
	}
	platformPrivateKey, err = x509.ReadPrivateKeyFromHex(hex.EncodeToString(platformPrivateKeyBytes))
	if err != nil {
		return errors.New("read platform private key error: " + err.Error())
	}
	bkPriPath := filepath.Join(dir, ".bank_private")
	bankPublicKeyBytes, err := os.ReadFile(bkPriPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return BankPublicNotExisted
		}
		return err
	}
	bankPublicKey, err = x509.ReadPublicKeyFromHex(hex.EncodeToString(bankPublicKeyBytes))
	if err != nil {
		return errors.New("read platform public key error: " + err.Error())
	}
	return nil
}

// 使用账户行公钥加密数据并Base64编码
func EncryptBase64(data []byte) (string, error) {
	encryptedBytes, err := sm2.Encrypt(bankPublicKey, data, rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// 使用平台私钥解密Base64编码的数据
func Base64Decrypt(data string) ([]byte, error) {
	dataBytes := []byte(data)
	var dstBytes []byte = make([]byte, len(dataBytes))
	len, err := base64.StdEncoding.Decode(dstBytes, dataBytes)
	if err != nil {
		return nil, err
	}
	return sm2.Decrypt(platformPrivateKey, dstBytes[0:len], sm2.C1C3C2)
}

func SignBase64(data []byte) (string,error){
	dst,err  := platformPrivateKey.Sign(rand.Reader, data, nil)
	if err != nil {
		return "",err
	}
	return base64.StdEncoding.EncodeToString(dst),nil
}