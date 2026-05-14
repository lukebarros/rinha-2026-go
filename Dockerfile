# Stage 1: pré-processamento (roda no build, não em runtime)
FROM golang:1.22-alpine AS preprocessor
WORKDIR /build
COPY go.mod .
COPY cmd/preprocess ./cmd/preprocess
COPY internal ./internal
COPY resources ./resources
RUN go run ./cmd/preprocess/main.go \
    -in  resources/references.json.gz \
    -out resources/references.bin

# Stage 2: build da API
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod .
COPY cmd/api ./cmd/api
COPY internal ./internal
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /api ./cmd/api

# Stage 3: imagem final mínima
FROM alpine:3.19
WORKDIR /app
COPY --from=builder      /api                            ./api
COPY --from=preprocessor /build/resources/references.bin ./resources/references.bin
ENV REFS_PATH=/app/resources/references.bin
EXPOSE 9999
CMD ["./api"]