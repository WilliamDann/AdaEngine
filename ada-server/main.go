package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-search"
)

const (
	startFEN    = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	maxDepth    = 64
	maxMovetime = 5 * time.Second
)

type moveRequest struct {
	FEN        string `json:"fen"`
	MovetimeMs int    `json:"movetimeMs"`
}

type moveResponse struct {
	Move  string `json:"move"`
	Score int    `json:"score"`
	Depth int    `json:"depth"`
	Nodes uint64 `json:"nodes"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// searchSemaphore caps concurrent engine searches. Each search spawns
// runtime.NumCPU() goroutines, so allowing N concurrent requests would
// oversubscribe the CPU N-fold. Cloud Run containers typically have 1-2
// vCPUs; one search at a time keeps response times predictable.
var searchSemaphore = make(chan struct{}, 1)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{"method not allowed"})
		return
	}

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid request body"})
		return
	}
	if req.FEN == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{"missing fen"})
		return
	}

	pos, err := fen.Parse(req.FEN)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid fen: " + err.Error()})
		return
	}
	pos.Zobrist = pos.ComputeZobrist()

	moves := movegen.LegalMoves(pos)
	if moves.Count() == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{"no legal moves"})
		return
	}

	movetime := time.Duration(req.MovetimeMs) * time.Millisecond
	if movetime <= 0 || movetime > maxMovetime {
		movetime = maxMovetime
	}

	searchSemaphore <- struct{}{}
	defer func() { <-searchSemaphore }()

	res := search.Search(pos, maxDepth, 0, movetime)

	writeJSON(w, http.StatusOK, moveResponse{
		Move:  res.Move.String(),
		Score: res.Score,
		Depth: res.Depth,
		Nodes: res.Nodes,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// corsMiddleware allows the frontend (deployed separately) to call this
// API cross-origin. Set CORS_ORIGIN to a specific origin in production
// instead of the permissive "*".
func corsMiddleware(next http.Handler) http.Handler {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "*"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	addr := ":" + envOr("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/move", handleMove)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"fen": startFEN})
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("ada-server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
