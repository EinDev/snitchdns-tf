// Package main provides the Terraform provider for SnitchDNS.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"snitchdns-tf/internal/provider"
)

// version is set by the goreleaser at build time
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/EinDev/snitchdns",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, nil), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
