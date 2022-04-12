package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ismurov/swaggerui"
)

const defaultSpecURL = "https://petstore.swagger.io/v2/swagger.json"

var (
	port    = flag.String("port", "8888", "listening port")
	specURL = flag.String("url", defaultSpecURL, "OpenAPI specification url")
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	flag.Parse()

	http.Handle("/", http.RedirectHandler("/swagger/?url="+*specURL, http.StatusFound))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerui.AssetsFS()))))

	fmt.Printf(""+
		"Run http server for serve OpenAPI Specification by URL (%s).\n"+
		"    Web interface will be available at link:\n"+
		"    http://127.0.0.1:%s/\n", *specURL, *port)
	return http.ListenAndServe(":"+*port, nil)
}
