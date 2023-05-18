package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/travix/protoc-gen-gotf/example/providerpb"
)

var (
	version string = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	// mkdir ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
	// mkdir ~/.terraform.d/plugins/travix.com/example/terraform-provider-example/1.0.0/darwin_arm64
	opts := providerserver.ServeOpts{
		Address: "travix.com/providers/example",
		Debug:   debug,
	}
	proxy := &ProviderProxy{}
	err := providerserver.Serve(context.Background(), providerpb.New(version, proxy), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
