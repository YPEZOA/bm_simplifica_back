# Etapa de build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar go.mod y go.sum primero (para cachear dependencias)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código y compilar
COPY . .
RUN go build -o main ./cmd

# Etapa final, más liviana
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
