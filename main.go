package main

import (
	"flag"

	"github.com/automato-io/terraform-provider-binocs/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New()}
	opts.Debug = debugMode

	plugin.Serve(opts)
}
