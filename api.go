package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type contextKey string

func (c contextKey) String() string {
	return "rates context key " + string(c)
}

var (
	contextDatabase = contextKey("database")

	queryLatest = `
		SELECT day, currency, rate
		FROM rates
		WHERE DAY = (SELECT max(day) FROM rates)
		ORDER BY currency
	`

	queryHistoricalByCurrency = `
		SELECT day, rate
		FROM rates
		WHERE currency = UCASE(?)
		ORDER BY day DESC
	`
)

func serve(c *cli.Context) error {
	db, err := sqlx.Connect("mysql", c.GlobalString("database"))
	if err != nil {
		log.Println("Failed to connect to database", c.GlobalString("database"))
		return err
	}

	r := chi.NewRouter()

	// Default setup for a REST service
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Injects database for downstream use (e.g. pulling latest and historical rates)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextDatabase, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// TODO: This could be improved a lot by caching the results in redis,
	// however as this is a technical test without redis as a requirment,
	// it's fetching data from mysql every request.
	r.Route("/api/v1/rates", func(r chi.Router) {
		r.Get("/latest", getLatest)
		r.Get("/history/{currency}", getHistoricByCurrency)
	})

	log.Println("Starting server on", c.String("listen"))
	return http.ListenAndServe(c.String("listen"), r)
}

func getLatest(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(contextDatabase).(*sqlx.DB)
	rows, err := db.Query(queryLatest)
	if err != nil {
		panic(err)
	}
	var (
		day      time.Time
		currency string
		rate     float64
	)
	out := make(map[string]float64)
	for rows.Next() {
		err := rows.Scan(&day, &currency, &rate)
		if err != nil {
			panic(err)
		}
		out[currency] = rate
	}
	rows.Close()
	render.Render(w, r, &RateLatestResponse{
		Rates:   out,
		Updated: day,
	})
}

// RateLatestResponse is a response struct for /api/v1/rates/latest
type RateLatestResponse struct {
	Rates   map[string]float64 `json:"rates"`
	Updated time.Time          `json:"updated"`
}

// Render implements render.Renderer interface for RateLatestResponse
func (rd *RateLatestResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getHistoricByCurrency(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	if currency == "" {
		render.Render(w, r, ErrNotFound)
		return
	}

	db := r.Context().Value(contextDatabase).(*sqlx.DB)
	rows, err := db.Query(queryHistoricalByCurrency, currency)
	if err != nil {
		panic(err)
	}

	var (
		day  time.Time
		rate float64
	)
	out := make(map[time.Time]float64)
	for rows.Next() {
		err := rows.Scan(&day, &rate)
		if err != nil {
			panic(err)
		}
		out[day] = rate
	}
	rows.Close()

	if len(out) == 0 {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.Render(w, r, &RateHistoricResponse{
		Rates: out,
	})
}

// RateHistoricResponse is a response struct for /api/v1/rates/latest
type RateHistoricResponse struct {
	Rates map[time.Time]float64 `json:"rates"`
}

// Render implements render.Renderer interface for RateHistoricResponse
func (rd *RateHistoricResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
