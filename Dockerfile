# Usamos Go completo para desarrollo
FROM golang:1.25-alpine

# Instalamos herramientas útiles
RUN apk add --no-cache git bash

# Instalar air para hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copiar mod y sum primero (para cachear dependencias)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código
COPY . .

# Exponer puerto
EXPOSE 8080

# Comando para desarrollo con hot reload
CMD ["air", "-c", ".air.toml"]
