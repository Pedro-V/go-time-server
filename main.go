package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	h := LogIpJson(CurrentTimeHandler{})

	mux := http.NewServeMux()
	mux.Handle("GET /now", h)
	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Listening on http://localhost:8080")
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

// Logging Middleware.
func LogIpJson(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := &slog.HandlerOptions{Level: slog.LevelInfo}
		handler := slog.NewJSONHandler(os.Stdout, options)
		mySlog := slog.New(handler)
		ip, _, _ := strings.Cut(r.RemoteAddr, ":")
		mySlog.Info("incoming ip", "ip_address", ip)
		h.ServeHTTP(w, r)
	})
}

// Handler.
type CurrentTimeHandler struct{}

func (h CurrentTimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var out []byte
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Add("Content-Type", "application/json")
		out = JsonOutput(now)
	} else {
		out = []byte(now.Format(time.RFC3339))
	}
	// Append newline for better curl output.
	out = append(out, '\n')
	w.Write(out)
}

// Response may be in JSON.
type DateTime struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      string `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
}

func JsonOutput(t time.Time) []byte {
	dt := DateTime{
		DayOfWeek:  t.Weekday().String(),
		DayOfMonth: t.Day(),
		Month:      t.Month().String(),
		Year:       t.Year(),
		Hour:       t.Hour(),
		Minute:     t.Minute(),
		Second:     t.Second(),
	}
	out, _ := json.Marshal(dt)
	return out
}
