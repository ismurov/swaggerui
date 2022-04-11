package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ismurov/swaggerui"
)

const (
	specName = "Petstore API"
	specFile = "swagger.json"
)

//go:embed swagger.json
var specFS embed.FS

var port = flag.String("port", "8888", "listening port")

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	flag.Parse()

	h, err := swaggerui.New(
		[]swaggerui.SpecFile{{
			Name: specName,
			Path: specFile,
		}},
		specFS,
	)
	if err != nil {
		return fmt.Errorf("failed to create swagger: %w", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
	})
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", h))

	fmt.Printf(""+
		"Run http server for serve OpenAPI Specification.\n"+
		"    Web interface will be available at link:\n"+
		"    http://127.0.0.1:%s/swagger-ui/\n", *port)
	return http.ListenAndServe(":"+*port, nil)
}
