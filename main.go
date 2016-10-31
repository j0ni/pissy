package main

import (
	"fmt"

	"mig.ninja/mig/pgp/pinentry"
)

func acquirePassphrase() []byte {
	request := &pinentry.Request{
		Desc:   "Passphrase Dialog for Pissy",
		Prompt: "Enter passphrase",
	}
	passphrase, err := request.GetPIN()
	maybeBail(err)
	return []byte(passphrase)
}

func maybeBail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Hello World")
}