# go-openapi

Azugo OpenAPI definition generation support

## Usage

1. Add `openapi\generator\generate.go` file to project

	```go
	package main

	import (
		"log"

		"github.com/lafriks-fork/goas"
	)

	func main() {
		p, err := goas.NewParser("../", "openapi.go", "", false)
		if err != nil {
			log.Fatalf("Can not initialize goas generator: %v", err)
		}
		err = p.CreateOASFile("openapi.json")
		if err != nil {
			log.Fatalf("Error while generating openapi.json: %v", err)
		}
	}
	```

2. Add `openapi/openapi.go` file to project

	```go
	//go:generate go run generator/generate.go

	// @version 1.0
	// @title Test API
	// @description TEST REST API
	// @contactName SIA ZZ Dats
	// @contactEmail zzdats@zzdats.lv
	// @contactURL https://www.zzdats.lv/
	// @server {{SERVER_URL}}
	// @security APIKeyAuth
	// @securityScheme APIKeyAuth apiKey header X-API-Key X-API-Key header using API key scheme. Example: "X-API-Key: {api-key}"
	package openapi

	import (
		_ "embed"
	)

	var (
		//go:embed openapi.json
		OpenAPIDefinition []byte
	)

	```

3. Add openapi handler to router:

	```go
	package routes

	import (
		oa "github.com/nobid-lsp-latvia/go-openapi"
	)

	func Init(a *template.App) error {
		r := &router{
			App:     a,
			openapi: oa.NewDefaultOpenAPIHandler(openapi.OpenAPIDefinition, a.App),
		}
		...
	}
	```

4. Call `go generate ./...` from terminal