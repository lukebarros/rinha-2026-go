package main

import (
	"log"
	"os"
	"runtime"

	"rinha/internal/knn"
	"rinha/internal/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	refsPath := getenv("REFS_PATH", "resources/references.bin")

	log.Println("Carregando dataset...")
	store, err := knn.Load(refsPath)
	if err != nil {
		log.Fatalf("Erro ao carregar dataset: %v", err)
	}
	log.Printf("Dataset carregado: %d vetores", store.Total())

	srv := server.New(store)
	log.Println("Servidor iniciado em :9999")
	if err := srv.ListenAndServe(":9999"); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}