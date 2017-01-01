package logger

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	statusCode    int
	statusWritten bool
	writer        http.ResponseWriter
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{writer: w, statusCode: 200}
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.statusWritten = true
	w.writer.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(data []byte) (n int, err error) {
	if !w.statusWritten {
		w.statusCode = 200
		w.statusWritten = true
	}
	return w.writer.Write(data)
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *loggingResponseWriter) Status() int {
	return w.statusCode
}

// New returns a http.Handler and wraps the next handler
func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		starttime := time.Now().Local()
		_, offset := starttime.Zone()
		offsetHours := int(offset / 3600)
		offsetMinutes := offset - (offsetHours * 3600)
		timestamp := fmt.Sprintf("%02v/%v/%v:%v:%v:%v +%02v%02v",
			starttime.Day(), starttime.Month().String()[0:3], starttime.Year(),
			starttime.Hour(), starttime.Minute(), starttime.Second(),
			offsetHours, offsetMinutes)
		lw := newLoggingResponseWriter(w)
		next.ServeHTTP(lw, r)
		user := "-"
		if lw.statusCode < 400 {
			u, _, ok := r.BasicAuth()
			if ok {
				user = u
			}
		}
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		logline := fmt.Sprintf("%v %v [%v] \"%v %v %v\" %v %v \"%v\" \"%v\" \"%v\"",
			ip, user, timestamp, r.Method, r.URL.String(), r.Proto, lw.statusCode, r.ContentLength,
			r.Header.Get("Referrer"), r.Header.Get("User-Agent"), r.Header.Get("Cookie"))
		fmt.Println(logline)
		return
	})
}
