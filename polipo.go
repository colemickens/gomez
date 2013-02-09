//+build ignore

// This should probably just be removed...

package main

import (
	"net/http"
	"os/exec"

	"code.google.com/p/go.net/proxy"
	"log"
	"time"
)

var (
	socksDialer  proxy.Dialer
	polipoClient *http.Client
)

func GetPolipoHttpClient(hostport string) *http.Client {
	if polipoClient == nil {
		log.Println("starting polipo client")
		exec.Command("polipo", "-c", "polipo.conf")
	}
	log.Println("sleeping for polipo startup")
	time.Sleep(10 * time.Second)

	if socksDialer == nil {
		var err error
		socksDialer, err = proxy.FromURL(hostport)
		if err != nil {
			panic(err)
		}

		return &http.Client{
			Transport: &http.Transport{Dial: socksDialer},
		}
	}
}
