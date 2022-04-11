// Package swaggerui implements http handler for serve Swagger UI interface.
// Note: Swagger UI bundle is embedded in package.
package swaggerui

import (
	_ "embed" // Imported for embed index.html template.
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	//go:embed templates/index.html
	indexTmplData string

	// indexTmpl is a private variable. To get access, you should use
	// getIndexTemplate function.
	indexTmpl     *template.Template
	indexTmplOnce sync.Once
)

// getIndexTemplate parses and returns swagger ui index page.
func getIndexTemplate() *template.Template {
	indexTmplOnce.Do(func() {
		indexTmpl = template.Must(template.New("index.html").Parse(indexTmplData))
	})
	return indexTmpl
}

type handler struct {
	tmpl        *template.Template
	specFiles   []SpecFile
	specPaths   map[string]struct{}
	assetServer http.Handler
	specsServer http.Handler
}

// SpecFile represents information about an OpenAPI specification file.
// The name field contains the displayed file information for UI.
// The path field must contains relative path to file in specification
// file system (specFS is the second argument to the New function).
type SpecFile struct {
	Name string
	Path string
}

// New creates new handler for serve Swagger UI interface with custom
// specification files. The handler strictly checks access to passed
// files in order to avoid unauthorized access to other files from
// file system.
//
// In case the handler will be mounted on a path other than the root path,
// use http.StripPrefix to skip the prefix.
func New(spec []SpecFile, specFS fs.FS) (http.Handler, error) {
	if specFS == nil {
		return nil, fmt.Errorf("specFS is nil")
	}

	specFiles := make([]SpecFile, 0, len(spec))
	specPaths := make(map[string]struct{}, len(spec))
	for _, f := range spec {
		pth := filepath.ToSlash(filepath.Clean("/" + f.Path))
		specFiles = append(specFiles, SpecFile{
			Name: f.Name,
			Path: path.Join("./specs", pth),
		})
		specPaths[path.Join("/specs", pth)] = struct{}{}
	}

	return &handler{
		tmpl:        getIndexTemplate(),
		specFiles:   specFiles,
		specPaths:   specPaths,
		assetServer: http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))),
		specsServer: http.StripPrefix("/specs/", http.FileServer(http.FS(specFS))),
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	switch {
	case upath == "/" || upath == "/index.html":
		if err := h.tmpl.ExecuteTemplate(w, "index.html", map[string][]SpecFile{
			"SpecFiles": h.specFiles,
		}); err != nil {
			// Internal package error, possible only during development.
			log.Printf("[swaggerui] http.Handler: failed to execute template: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return

	case strings.HasPrefix(upath, "/assets/"):
		h.assetServer.ServeHTTP(w, r)
		return

	case h.hasAccessToSpec(upath):
		h.specsServer.ServeHTTP(w, r)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (h *handler) hasAccessToSpec(pth string) bool {
	_, exists := h.specPaths[pth]
	return exists
}
