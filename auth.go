package main

import "net/http"

type Auth struct {
	s *Storage
}

func (a *Auth) checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, ok := a.s.GetUserByUserName(username)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user.Password != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
