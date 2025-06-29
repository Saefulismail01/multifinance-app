# Tahap build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Salin go mod dan sum files
COPY go.mod go.sum ./

# Unduh semua dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Buat aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Tahap final
FROM alpine:latest

WORKDIR /app

# Salin file biner yang telah dibuat sebelumnya dari tahap sebelumnya
COPY --from=builder /app/main .

# Salin file .env
COPY .env .

# Install CA certificates (diperlukan untuk koneksi HTTPS)
RUN apk --no-cache add ca-certificates

# Buka port 8080
EXPOSE 8080

# Perintah untuk menjalankan executable
CMD ["./main"]
