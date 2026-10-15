package main

import (
	"fmt"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("ICTSC_WEB_PUSH_VAPID_PUBLIC_KEY=%s\n", publicKey)
	fmt.Printf("ICTSC_WEB_PUSH_VAPID_PRIVATE_KEY=%s\n", privateKey)
}
