package server

import (
	"db_tools/cmd/database"
	"db_tools/cmd/seeder"
	"errors"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"net/http"
	"os"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "pong")
	return
}
func migrateUp(w http.ResponseWriter, r *http.Request) {
	database.UpMigration()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "OK")
	return
}
func migrateDown(w http.ResponseWriter, r *http.Request) {
	database.DownMigration()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "OK")
	return
}
func seed(w http.ResponseWriter, r *http.Request) {
	seeder.SeedDatabase()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "OK")
	return
}
func fresh(w http.ResponseWriter, r *http.Request) {
	database.DownMigration()
	database.UpMigration()
	seeder.SeedDatabase()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "OK")
	return
}
func exit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	fmt.Fprintln(w, "OK")
	os.Exit(0)
	return
}

func logReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cmn.Log(fmt.Sprintf("http %s request on %s", r.Method, r.URL), cmn.LogLevels.Info)
			next.ServeHTTP(w, r)
		})
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/up", migrateUp)
	mux.HandleFunc("/down", migrateDown)
	mux.HandleFunc("/seed", seed)
	mux.HandleFunc("/fresh", fresh)
	mux.HandleFunc("/exit", exit)
	wrappedMux := logReqMiddleware(mux)

	port := os.Getenv("DB_TOOLS_PORT")
	cmn.Log("Database tool connected on port %s", cmn.LogLevels.Info, port)
	cmn.Log("Seedeing database", cmn.LogLevels.Info)
	err := http.ListenAndServe(":"+port, wrappedMux)
	if errors.Is(err, http.ErrServerClosed) {
		cmn.Log("server closed", cmn.ErrorLevels.Info)
	} else if err != nil {
		// error starting the server
		cmn.HandleError(err, cmn.ErrorLevels.Critical)
	}
}
