package api

import (
	"errors"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"net/http"
	"os"
	"path/filepath"
)

var mux *http.ServeMux

const Prefix = "/api"

func GetHttpMux() *http.ServeMux {
	if mux == nil {
		cmn.HandleError(errors.New("Server is not yet initalized"),
			cmn.ErrorLevels.Critical)
	}
	return mux
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "pong")
	// io.WriteString(w, "Hello, HTTP!\n")
	return
}

func pingWithCors(w http.ResponseWriter, r *http.Request) {
	var allowedMethods = []string{http.MethodGet, http.MethodOptions}
	allowedHeaders := []string{"Content-Type"}
	configureHttpAccessControlHeaders(w, allowedMethods, allowedHeaders)
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "pong")
		// io.WriteString(w, "Hello, HTTP!\n")
		return
	}
}

// serves the openapi.json file
func openapiFile(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	path := os.Getenv("OPENAPI_FILE")
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			HandleRequestError(w, cmn.ErrorCodes.ReadFailed, "")
		}
		path = filepath.Join(path, "cmd", "api", "openapi.json")
	}
	http.ServeFile(w, r, path)
	return
}

func registerSystemHandlers(mux *http.ServeMux) {
	mux.HandleFunc(SystemReadRoutes.Ping, ping)
	mux.HandleFunc(SystemReadRoutes.Ping2, ping)
	mux.HandleFunc(SystemReadRoutes.PingCors, pingWithCors)
	mux.HandleFunc(SystemReadRoutes.OpenAPI, openapiFile)
	mux.HandleFunc(SystemReadRoutes.OpenAPI2, openapiFile)
}

// must be the last handler called
func RegisterInvalidUrlHandler() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		cmn.Log(fmt.Sprintf("Invalid route called %s\n", r.URL), cmn.LogLevels.Warning)
		http.Error(w, "Invalid route", http.StatusBadRequest)
	}
	mux.HandleFunc("/", handler)
}

func InitApi() {
	cmn.Log("Init Api", cmn.LogLevels.Info)
	mux = http.NewServeMux()
	registerSystemHandlers(mux)
}

func logReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cmn.Log(fmt.Sprintf("http %s request on %s", r.Method, r.URL), cmn.LogLevels.Info)
			next.ServeHTTP(w, r)
		})
}

func StartApi() {
	port := os.Getenv("BACKEND_PORT")
	cmn.Log("Starting BAM backend on port %s", cmn.LogLevels.Info, port)
	wrappedMux := logReqMiddleware(mux)
	err := http.ListenAndServe(":"+port, wrappedMux)
	if errors.Is(err, http.ErrServerClosed) {
		cmn.Log("server closed", cmn.ErrorLevels.Info)
	} else if err != nil {
		// error starting the server
		cmn.HandleError(err, cmn.ErrorLevels.Critical)
	}
}
