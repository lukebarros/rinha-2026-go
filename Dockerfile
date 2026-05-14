FROM --platform=linux/amd64 golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN mkdir -p /out/resources

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-s -w' -o /out/build-index ./cmd/preprocess
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-s -w' -o /out/api ./cmd/api
RUN chmod +x /out/build-index /out/api

COPY resources/references.json.gz /tmp/references.json.gz
RUN /out/build-index -in /tmp/references.json.gz -out /out/resources/references.bin && \
    ls -lah /out/resources && \
    test -f /out/resources/references.bin && \
    rm -f /tmp/references.json.gz

FROM --platform=linux/amd64 alpine:3.20 AS runtime
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/api /app/api
COPY --from=builder /out/resources/references.bin /app/resources/references.bin
COPY --from=builder /src/resources/normalization.json /app/resources/normalization.json
COPY --from=builder /src/resources/mcc_risk.json /app/resources/mcc_risk.json

RUN chmod +x /app/api

EXPOSE 9999
CMD ["/app/api"]