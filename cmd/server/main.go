package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JaniL/tulospalvelu"
	"github.com/JaniL/tulospalvelu/internal/api"
	"github.com/JaniL/tulospalvelu/internal/geofilter"
	"github.com/JaniL/tulospalvelu/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		log.Printf("%s %s %d", r.Method, r.URL.Path, sw.status)
	})
}

func main() {
	port := envOr("PORT", "8080")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL ympäristömuuttuja puuttuu")
	}
	seed := os.Getenv("SEED_TEST_DATA") == "true"
	geoRestrict := os.Getenv("GEO_RESTRICT") == "true"
	trustXFF := os.Getenv("TRUST_XFF") == "true"

	var filter *geofilter.Filter
	if geoRestrict {
		f, err := geofilter.New()
		if err != nil {
			log.Fatalf("IP-listan jäsentäminen epäonnistui: %v", err)
		}
		filter = f
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("virheellinen DATABASE_URL: %v", err)
	}
	cfg.MaxConns = 8

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("tietokantayhteys epäonnistui: %v", err)
	}
	defer pool.Close()

	if err := store.Migrate(ctx, pool, seed); err != nil {
		log.Fatalf("tietokannan alustus epäonnistui: %v", err)
	}

	handler := withLogging(api.New(store.NewPG(pool), tulospalvelu.StaticFS, api.Options{
		GeoRestrict: geoRestrict,
		TrustXFF:    trustXFF,
		Filter:      filter,
	}).Handler())

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("sammutus: %v", err)
		}
	}()

	log.Printf("tulospalvelu kuuntelee portissa %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
