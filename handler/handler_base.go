//jasonxu
package handler

import (
	"github.com/Jason916/peanut_core/log"
	"net/http"
	"runtime/debug"
)

func NewHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Error: ", err, debug.Stack())
			}
		}()
		logrespwriter := &LogRespWriter{responseWriter: resp}
		handler.ServeHTTP(logrespwriter, req)
		log.Info("%v %v %v (%v)", logrespwriter.Status(), req.Method, req.URL.Path, req.RemoteAddr)
	})
}

