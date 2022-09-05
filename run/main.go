package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	_ "github.com/fin-fet/coredns-blocklist"
)

func init() {
	dnsserver.Directives = append(
		[]string{"log", "blocklist"},
		dnsserver.Directives...,
	)
}

func main() {
	coremain.Run()
}
