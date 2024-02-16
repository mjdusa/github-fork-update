package httptest

import (
	"fmt"
	"io"
	"net/http"
	htst "net/http/httptest"
)

type Server struct {
	Mux        *http.ServeMux
	APIHandler *http.ServeMux
	Path       string
	Server     *htst.Server
}

func NewHTTPTestServer(path string, log io.Writer) (*Server, error) {
	srvr := Server{
		Mux:        nil,
		APIHandler: nil,
		Path:       path,
		Server:     nil,
	}

	// mux is the HTTP request multiplexer used with the test server.
	srvr.Mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	srvr.APIHandler = http.NewServeMux()
	pattern := path + "/"
	hler := http.StripPrefix(path, srvr.Mux)
	srvr.APIHandler.Handle(pattern, hler)
	srvr.APIHandler.HandleFunc("/", func(respWrtr http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(log, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(log)
		fmt.Fprintln(log, "\t"+req.URL.String())
		fmt.Fprintln(log)
		fmt.Fprintln(log, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(log, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(respWrtr, "Client.BaseURL path prefix is not preserved in the request URL.",
			http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	srvr.Server = htst.NewServer(srvr.APIHandler)

	return &srvr, nil
}

func (srvr *Server) Close() {
	if srvr.Server != nil {
		srvr.Server.Close()
	}
}
