package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gunturaf/sukab-property/domain/property"
)

func New(importer property.Importer, lister property.Lister) *Server {
	if importer == nil {
		panic("importer service is not set")
	}
	if lister == nil {
		panic("lister service is not set")
	}

	srv := &Server{
		maxImportFileSizeBytes: 2 << 20, // accepts 2 MB of payload
		mux:                    http.NewServeMux(),
		importer:               importer,
		lister:                 lister,
	}

	srv.mux.Handle("/property/import", srv.handleImport())
	srv.mux.Handle("/property", srv.handleList())

	return srv
}

type Server struct {
	maxImportFileSizeBytes int64
	mux                    *http.ServeMux

	importer property.Importer
	lister   property.Lister
}

// Run this function will run HTTP Server
// and listen to any interrupt/termination signal
// from OS to handle graceful shutdown of the HTTP server.
func (h *Server) Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: h.mux,
		// the following attrs are as the best practice
		// for Go http server, based on this article: https://bruinsslot.jp/post/go-secure-webserver/
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	log.Printf("Serving HTTP server on %s...\n", addr)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		if err := srv.Close(); err != nil {
			log.Fatalf("HTTP close error: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}

func newHTTPErr(code int, msg string) *httpErr {
	return &httpErr{
		code: code,
		msg:  msg,
	}
}

// httpErr is a generic way to represent common HTTP err.
type httpErr struct {
	code int
	msg  string
}

func (he *httpErr) Error() string {
	return he.msg
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (srv *Server) respondSuccess(w http.ResponseWriter, data any) {
	setJSONHeader(w)
	w.WriteHeader(http.StatusOK)
	jsonBody, _ := json.Marshal(data)
	w.Write(jsonBody)
}

func (srv *Server) respondError(err error, w http.ResponseWriter) {
	var underlyingErr *httpErr
	// try to convert error interface (abstract) to concrete type:
	if !errors.As(err, &underlyingErr) {
		// assume this is a generic 500 err.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// properly write http err code + JSON body.
	setJSONHeader(w)
	w.WriteHeader(underlyingErr.code)
	jsonBody := map[string]string{"error": err.Error()}
	jsonBodyByte, _ := json.Marshal(jsonBody)
	w.Write(jsonBodyByte)
}
