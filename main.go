package main

import (
	"errors"
	"fmt"
	"github.com/sufu777/JumpProxy/encrypt"
	"net/http"
)

func main() {
	err := encrypt.SetupKeys()
	if err != nil {
		if errors.Is(err, encrypt.PlatformPrivateKeyNotExisted) || errors.Is(err, encrypt.BankPublicNotExisted) {
			fmt.Printf("%s", err.Error())
			return
		}
	}
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

}
