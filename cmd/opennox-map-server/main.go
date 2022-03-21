package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/noxworld-dev/opennox-lib/common"
	"github.com/noxworld-dev/opennox-lib/maps"
)

var (
	fPath = flag.String("data", ".", "path to Nox game data")
	fHost = flag.String("host", fmt.Sprintf(":%d", common.GameHTTPPort), "host ot listen on")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	srv := maps.NewServer(*fPath)
	return http.ListenAndServe(*fHost, srv)
}
