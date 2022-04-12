package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ismurov/swaggerui"
)

// Default configuration for workable example (optional).
var (
	_, self, _, _ = runtime.Caller(0)
	defaultDir    = filepath.Join(filepath.Dir(self), "..", "..", "testdata")
	defaultSpec   = "api-spec.yaml"
)

var (
	port = flag.String("port", "8888", "listening port")
	dir  = flag.String("dir", defaultDir, "path to folder with OpenAPI specifications")
	spec = flag.String("spec", defaultSpec, "relative path to specification file in --dir folder")
)

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
			Name: "API Spec",
			Path: *spec,
		}},
		os.DirFS(*dir),
	)
	if err != nil {
		return fmt.Errorf("failed to create swagger: %w", err)
	}

	http.Handle("/", http.RedirectHandler("/swagger-ui/", http.StatusFound))
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", h))

	fmt.Printf(""+
		"Run http server for serve OpenAPI Specification.\n"+
		"    Web interface will be available at link:\n"+
		"    http://127.0.0.1:%s/swagger-ui/\n", *port)
	return http.ListenAndServe(":"+*port, nil)
}
