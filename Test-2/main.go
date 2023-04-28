package main

import (
	"log"
	"net/http"
	"os"
	"io"
	"mime"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
)

func JSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt == "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func (h http.Handler) http.Handler  {
		return handlers.LoggingHandler(dst, h)		
	}
}


func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}


func main() {

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND,0664)
	if err != nil {
		log.Fatal(err)
	}

	loggingHandler := newLoggingHandler(logFile)

	authHandler := httpauth.SimpleBasicAuth("alice", "pa$$word")
	
	mux := http.NewServeMux()
	
	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", loggingHandler(authHandler(JSONHandler(finalHandler))))

	log.Println("Listening on :3000.....")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}