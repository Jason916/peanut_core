//jasonxu
package handler

import "net/http"

type LogRespWriter struct {
	responseWriter http.ResponseWriter
	status         int
}

func (logrespwriter *LogRespWriter) Status() int {
	return logrespwriter.status
}

func (logrespwriter *LogRespWriter) ResponseWriter() http.ResponseWriter {
	return logrespwriter.responseWriter
}

func (logrespwriter *LogRespWriter) Write(data []byte) (int, error) {
	return logrespwriter.responseWriter.Write(data)
}

func (logrespwriter *LogRespWriter) Header() http.Header {
	return logrespwriter.responseWriter.Header()
}

func (logrespwriter *LogRespWriter) WriteHeader(status int) {
	logrespwriter.status = status
	logrespwriter.responseWriter.WriteHeader(status)
}


