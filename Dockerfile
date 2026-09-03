# ==============================================================================
# STAGE 1: Build Backend Go (Golang 1.22 Native Engine)
# ==============================================================================
FROM golang:1.22-alpine AS go-builder

WORKDIR /build

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/edgego-voice ./cmd/server

# ==============================================================================
# STAGE 2: Imagem Final de Produção (Alpine Minimal ~18MB)
# ==============================================================================
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata ffmpeg

# Copiar binário compilado Go
COPY --from=go-builder /app/edgego-voice /app/edgego-voice

# Copiar bundle estático gerado pelo Next.js (pasta ./web)
COPY web/ /app/web/

# Copiar arquivos de configuração inicial
COPY voices.json /app/voices.json
RUN mkdir -p /app/sfx

EXPOSE 5050

ENTRYPOINT ["/app/edgego-voice"]
