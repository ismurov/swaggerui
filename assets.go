package swaggerui

import (
	"embed"
	"io/fs"
)

var (
	//go:embed assets/*
	assetsFolderFS embed.FS
	assetsFS, _    = fs.Sub(assetsFolderFS, "assets")
)

// AssetsFS returns file system with Swagger UI assets. It can be used with
// a configuration based on the URL query string.
//
// Simple configuration:
//
//	http://www.example.com/index.html?url=/path/to/swagger.json
//
// See the following link for more information:
//
//	https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/configuration.md
//
// Usage example:
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//	    http.Redirect(w, r, "/swagger/index.html?url=https://petstore.swagger.io/v2/swagger.json", http.StatusFound)
//	})
//	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerui.AssetsFS()))))
func AssetsFS() fs.FS {
	return assetsFS
}
