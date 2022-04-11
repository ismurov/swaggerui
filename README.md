# Swagger UI

Package swaggerui implements http handler for serve Swagger UI interface.<br/>
Swagger UI bundle is embedded in package.

### Use cases
* Using http.Handler for serve Swagger UI interface with custom specification files from fs.FS.
* Server Swagger UI bundle with URL query configuration.

### Runnable examples
* [embed](./examples/embed/main.go)
* [filesystem](./examples/filesystem/main.go)
* [static](./examples/static/main.go)

### Example Usage
```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ismurov/swaggerui"
)

func main() {
	h, err := swaggerui.New(
		[]swaggerui.SpecFile{
			{
				Name: "Sample API",
				Path: "api-spec.yaml",
			},
			{
				Name: "Petstore API",
				Path: "petstore/swagger.json",
			},
		},
		os.DirFS("/path/to/specs"),
	)
	if err != nil {
		log.Fatalf("swaggerui.New: %v", err)
	}

	http.Handle("/swagger/", http.StripPrefix("/swagger/", h))

	http.ListenAndServe(":8080", nil)
}
```
