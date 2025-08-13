package openapi

import (
	"bytes"
	"io/fs"
	"mime"
	"path"
	"path/filepath"
	"slices"
	"strings"

	openapi "github.com/nobid-lsp-latvia/go-openapi/openapi"

	"azugo.io/azugo"
	"azugo.io/core"
)

// OpenAPI instance
type OpenAPI struct {
	definition []byte
	static     fs.FS
}

// NewOpenAPIHandler creates new OpenAPI handler instance
func NewOpenAPIHandler(static fs.FS, definition []byte) *OpenAPI {
	return &OpenAPI{
		definition: definition,
		static:     static,
	}
}

// NewDefaultOpenAPIHandler creates new OpenAPI handler instance
//
// If environment is not provided, swagger will be available in Development and Staging environments by default
func NewDefaultOpenAPIHandler(definition []byte, a *azugo.App, environment ...core.Environment) *OpenAPI {
	// Check if swagger should be allowed in current environment
	if (len(environment) > 0 && !slices.Contains(environment, a.Env())) ||
		(len(environment) == 0 && !a.Env().IsDevelopment() && !a.Env().IsStaging()) {
		return nil
	}

	// Swagger documentation
	static, err := fs.Sub(openapi.Docs, "public")
	if err != nil {
		panic(err)
	}

	oa := &OpenAPI{
		definition: definition,
		static:     static,
	}

	a.Get("/swagger/swagger.json", oa.SwaggerJSON)
	a.Get("/swagger/{filepath?:*}", oa.Docs)
	a.Get("/docs/{filepath?:*}", oa.Docs)

	return oa
}

func replaceBaseURL(data []byte, serverURL []byte) []byte {
	return bytes.ReplaceAll(data, []byte("{{SERVER_URL}}"), serverURL)
}

// SwaggerJSON writes swagger definition to response
func (o *OpenAPI) SwaggerJSON(ctx *azugo.Context) {
	ctx.Header.Set("Content-Type", "application/json")
	if _, err := ctx.Context().Write(replaceBaseURL(o.definition, []byte(ctx.BaseURL()))); err != nil {
		ctx.Log().Sugar().Errorf("Error serving swagger.json - %v", err)
	}
}

// Docs writes Swagger UI or Redoc static files to response
func (o *OpenAPI) Docs(ctx *azugo.Context) {
	fsPath := strings.TrimPrefix(ctx.Path(), ctx.BasePath())
	fsPath = strings.Trim(fsPath, "/")
	s, err := fs.Stat(o.static, fsPath)
	if err != nil {
		ctx.Log().Sugar().Errorf("Error serving static documentation data - %v; fsPath=%s", err, fsPath)
		ctx.NotFound()
		return
	}
	if s.IsDir() {
		fsPath = path.Join(fsPath, "index.html")
	}
	data, err := fs.ReadFile(o.static, fsPath)
	if err != nil {
		ctx.Log().Sugar().Errorf("Error reading file - %v; fsPath=%s", err, fsPath)
		ctx.NotFound()
		return
	}
	ext := filepath.Ext(fsPath)
	ctx.Header.Set("Content-Type", mime.TypeByExtension(ext))
	if ext == ".html" {
		data = replaceBaseURL(data, []byte(ctx.BaseURL()))
	}
	if _, err := ctx.Context().Write(data); err != nil {
		ctx.Log().Sugar().Errorf("Error serving static documentation data - %v", err)
	}
}
