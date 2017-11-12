//jasonxu
package handler

import (
	"../log"
	"net/http"
	"runtime/debug"
)

func NewHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Error: %+v\n %s", err, debug.Stack())
			}
		}()
		logrespwriter := &LogRespWriter{responseWriter: resp}
		handler.ServeHTTP(logrespwriter, req)
		log.Info("%v %v %v (%v)", logrespwriter.Status(), req.Method, req.URL.Path, req.RemoteAddr)
	})
}
