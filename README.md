# Layanan MultiFinance

Layanan API RESTful untuk menangani transaksi keuangan dengan batas kredit pelanggan.

## Fitur

- Manajemen pelanggan
- Manajemen batas kredit
- Pemrosesan transaksi
- Validasi input
- Dukungan transaksi database
- Penanganan error terstruktur

## Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go 1.23
- **Framework Web**: Gin
- **Basis Data**: MySQL
- **ORM**: SQLx
- **Manajemen Environment**: godotenv

## Persyaratan Sistem

- Go 1.23 atau lebih tinggi
- MySQL 5.7 atau lebih tinggi
- Git

## Instalasi

1. Clone repositori:
   ```bash
   git clone https://github.com/yourusername/multifinance.git
   cd multifinance
   ```

2. Atur variabel environment:
   ```bash
   cp .env.example .env
   ```
   Perbarui file `.env` dengan kredensial database dan pengaturan lainnya.

3. Install dependensi:
   ```bash
   go mod download
   ```

4. Inisialisasi database:
   - Buat database MySQL baru
   - Jalankan migrasi database (jika ada)

## Menjalankan Aplikasi

```bash
# Menjalankan server
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080` secara default.

## Endpoint API

### Pemeriksaan Kesehatan

- `GET /health` - Memeriksa status layanan

### Transaksi

#### Membuat Transaksi Baru

- **URL**: `POST /api/v1/transactions`
- **Request Body**:
  ```json
  {
    "customer_nik": "1234567890123456",
    "otr": 1000000,
    "admin_fee": 50000,
    "installment": 3,
    "interest": 100000,
    "asset_name": "Laptop",
    "tenor": 3
  }
  ```
- **Response Sukses**:
  ```json
  {
    "code": 201,
    "status": "Created",
    "message": "Transaksi berhasil dibuat",
    "data": {
      "contract_number": "CON-1234567890",
      "customer_nik": "1234567890123456",
      "otr": 1000000,
      "admin_fee": 50000,
      "installment": 3,
      "interest": 100000,
      "asset_name": "Laptop",
      "created_at": "2025-06-29T13:50:08+07:00"
    }
  }
  ```

## Struktur Proyek

```
.
├── cmd/
│   └── main.go           # Titik masuk aplikasi
├── config/
│   └── database.go      # Konfigurasi database
├── delivery/
│   ├── controller/      # Penangan HTTP
│   ├── dto/              # Data Transfer Objects
│   └── server.go         # Konfigurasi server dan routing
├── model/                # Model database
├── repository/           # Lapisan akses data
├── service/              # Logika bisnis
├── usecase/              # Use case aplikasi
├── .env.example          # Contoh environment variable
├── go.mod               
└── go.sum
```

## Penanganan Error

API mengembalikan respons error standar dalam format berikut:

```json
{
  "code": 400,
  "status": "Bad Request",
  "message": "Validasi gagal",
  "errors": [
    {
      "field": "customer_nik",
      "message": "customer_nik wajib diisi"
    }
  ]
}
```

## Variabel Environment

- `DB_HOST`: Host database
- `DB_PORT`: Port database
- `DB_USER`: Pengguna database
- `DB_PASSWORD`: Kata sandi database
- `DB_NAME`: Nama database
- `SERVER_PORT`: Port server (default: 8080)

# ERD
[Lihat Diagram ERD](https://dbdiagram.io/d/ERD-6861365af413ba3508615838)

# Arsitektur
[Lihat Diagram Arsitektur Aplikasi](https://www.mermaidchart.com/app/projects/83ff1754-b639-4e98-822c-c8232b169ec0/diagrams/f3ef9241-a51f-41e9-ab7c-f24047968b8e/share/invite/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb2N1bWVudElEIjoiZjNlZjkyNDEtYTUxZi00MWU5LWFiN2MtZjI0MDQ3OTY4YjhlIiwiYWNjZXNzIjoiVmlldyIsImlhdCI6MTc1MTIwMjI3MX0.CtSEpFYs3Pdrw-8WeUkbKHn-1nx8IU3hL4jo0sou7po)

# Docker
[Lihat Dockerfile](./Dockerfile)
[Lihat docker-compose.yml](./docker-compose.yml)

# Docker Hub
[Lihat Docker Hub](https://hub.docker.com/r/saefulismail01/multifinance-app)


