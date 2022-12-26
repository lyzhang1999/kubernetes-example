package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	http.Handle("/ping", WithLogging(pingHandler()))
	http.Handle("/http", WithLogging(httpHandler()))

	addr := "0.0.0.0:8080"
	logrus.WithField("addr", addr).Info("starting server")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		logrus.WithField("event", "start server").Fatal(err)
	}
}

func pingHandler() http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(rw, "pong")
	}
	return http.HandlerFunc(fn)
}

func httpHandler() http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().Unix())
		var resCodeArray = [...]int{200, 400, 500}
		resCode := resCodeArray[rand.Intn(len(resCodeArray))]
		rw.WriteHeader(resCode)
		_, _ = fmt.Fprintf(rw, "OK!")
	}
	return http.HandlerFunc(fn)
}

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func WithLogging(h http.Handler) http.Handler {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: rw, // compose original http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lrw, req) // inject our implementation of http.ResponseWriter

		duration := time.Since(start)

		logrus.SetFormatter(&logrus.TextFormatter{})

		rand.Seed(time.Now().Unix())
		var resCodeArray = [...]int{200, 400, 500}
		resCode := resCodeArray[rand.Intn(len(resCodeArray))]

		logrus.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"status":   resCode,
			"duration": duration,
			"size":     responseData.size,
		}).Info("request completed")

		logrus.SetFormatter(&logrus.JSONFormatter{})

		logrus.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"status":   resCode,
			"duration": duration,
			"size":     responseData.size,
		}).Info("request completed")

	}
	return http.HandlerFunc(loggingFn)
}
