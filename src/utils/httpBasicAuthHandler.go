package utils

import (
	"log"
	"net/http"
)

func HTTPBasicAuthHandler(basicUsername string, basicPassword string) func(http.HandlerFunc) http.HandlerFunc {
	authEnabled := false
	if basicUsername != "" || basicPassword != "" {
		log.Println("HTTP Basic Auth enabled")
		authEnabled = true
	} else {
		log.Println("Warning! HTTP Basic Auth disabled, your vulnerability metrics are readable by anyone!")
	}

	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if authEnabled {
				user, pass, _ := r.BasicAuth()

				if basicUsername != user || basicPassword != pass {
					w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					http.Error(w, "Unauthorized.", http.StatusUnauthorized)
					return
				}
			}

			h.ServeHTTP(w, r)
		}
	}
}
