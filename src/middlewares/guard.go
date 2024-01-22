package middlewares

import (
	"net/http"

	"github.com/atanurdemir/gatekeeper/src/store"
	"github.com/atanurdemir/gatekeeper/src/utils"
)

func GuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := utils.GetIP(r)

		if store.IsIPBanned(ip) {
			http.Error(w, "Forbidden - IP Banned", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
