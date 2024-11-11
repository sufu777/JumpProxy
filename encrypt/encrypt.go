package encrypt

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/sufu777/JumpProxy/encrypt/sm4"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"os"
	"path/filepath"
)

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

type Message struct {
	Iv            string `json:"IV"`
	SymmetricKey  string `json:"SymmetricKey"`
	MessageBody   string `json:"MessageBody"`
	MessageHeader string `json:"MessageHeader"`
	Sign          string `json:"Sign"`
}

func DecryptBankMessage(messageStr string) (map[string]any, error) {
	var message Message
	err := json.Unmarshal([]byte(messageStr), &message)
	if err != nil {
		return nil, errors.New("error parse Json")
	}
	var ivDecoded []byte
	_, err = base64.StdEncoding.Decode(ivDecoded, []byte(message.Iv))
	if err != nil {
		return nil, errors.New("error decode iv")
	}
	decryptedIv, err := sm2.Decrypt(platformPrivateKey, ivDecoded, sm2.C1C3C2)
	if err != nil {
		return nil, errors.New("error decrypt iv: " + err.Error())
	}
	var symmetricKeyDecoded []byte
	_, err = base64.StdEncoding.Decode(symmetricKeyDecoded, []byte(message.SymmetricKey))
	if err != nil {
		return nil, errors.New("error decode symmetric key")
	}
	decryptedSymmetricKey, err := sm2.Decrypt(platformPrivateKey, symmetricKeyDecoded, sm2.C1C3C2)
	if err != nil {
		return nil, errors.New("error decrypt symmetric key: " + err.Error())
	}
	var requestBodyDecoded []byte
	_, err = base64.StdEncoding.Decode(requestBodyDecoded, []byte(message.MessageBody))
	if err != nil {
		return nil, errors.New("error decode message body")
	}
	out, err := sm4.Decrypt(decryptedSymmetricKey, decryptedIv, requestBodyDecoded)
	if err != nil {
		return nil, errors.New("error decrypt message: " + err.Error())
	}
	var req = make(map[string]any)
	err = json.Unmarshal(out, &req)
	if err != nil {
		return nil, errors.New("error unmarshal json: " + err.Error())
	}
	return req, nil
}

func EncodeMessage(rspBody map[string]any) ([]byte, error) {
	return nil, nil
}
