package main

import (
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		if req.RequestURI != "/api/health" {
			log.Printf("%s %s %s %s", req.Method, req.RequestURI, req.Proto, time.Since(start))
		}
	})
}

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
				return
			}
		}()
		next.ServeHTTP(w, req)
	})
}

var disallowWithoutAuthList []*regexp.Regexp = []*regexp.Regexp{
	//regexp.MustCompile(".*/api/v1/health.*"),
	regexp.MustCompile(".*/api/v1/user.*"),
	regexp.MustCompile(".*/api/v1/chat.*"),
	regexp.MustCompile(".*/api/v1/card.*"),
	regexp.MustCompile(".*/api/v1/[0-9]+/msg.*"),
	regexp.MustCompile(".*/api/v1/file.*"),
	regexp.MustCompile(".*/api/v1/transaction.*"),
	regexp.MustCompile(".*/api/v1/[0-9]+/post.*"),
}

func AuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for i := 0; i < len(disallowWithoutAuthList); i++ {
			if disallowWithoutAuthList[i].MatchString(req.RequestURI) {
				c, err := req.Cookie("token")
				if err != nil {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
				tknStr := c.Value
				_, ok := CookieToUserMap[tknStr]
				if ok {
					next.ServeHTTP(w, req)
					return
				} else {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			}
		}
		next.ServeHTTP(w, req)
	})

}
