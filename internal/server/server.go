package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"rinha/internal/knn"
	"rinha/internal/vectorizer"
)

type Server struct {
	store *knn.Store
	ready atomic.Bool
}

func New(store *knn.Store) *Server {
	s := &Server{store: store}
	s.ready.Store(true)
	return s
}

type reqBody struct {
	ID          string                 `json:"id"`
	Transaction vectorizer.Transaction `json:"transaction"`
	Customer    vectorizer.Customer    `json:"customer"`
	Merchant    vectorizer.Merchant    `json:"merchant"`
	Terminal    vectorizer.Terminal    `json:"terminal"`
	LastTx      *vectorizer.LastTx     `json:"last_transaction"`
}

type respBody struct {
	Approved   bool    `json:"approved"`
	FraudScore float32 `json:"fraud_score"`
}

func (s *Server) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ready",       s.handleReady)
	mux.HandleFunc("/fraud-score", s.handleFraudScore)
	return (&http.Server{Addr: addr, Handler: mux}).ListenAndServe()
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if s.ready.Load() {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(503)
}

func (s *Server) handleFraudScore(w http.ResponseWriter, r *http.Request) {
	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// FALLBACK — nunca retorna erro HTTP
		log.Printf("decode: %v", err)
		writeJSON(w, respBody{Approved: true, FraudScore: 0.4})
		return
	}

	vreq := &vectorizer.Request{
		Transaction: req.Transaction,
		Customer:    req.Customer,
		Merchant:    req.Merchant,
		Terminal:    req.Terminal,
		LastTx:      req.LastTx,
	}

	vec, key := vectorizer.Vectorize(vreq)
	score := s.store.Parts()[key].KNN(vec)
	writeJSON(w, respBody{Approved: score < 0.6, FraudScore: score})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}