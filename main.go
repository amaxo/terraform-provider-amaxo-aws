package main

import (
	"github.com/amaxo/terraform-provider-amaxo-aws/amaxoaws"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: amaxoaws.Provider,
	})
}
