# Arsitektur Aplikasi Multifinance

Dokumen ini menjelaskan arsitektur dan alur kerja aplikasi Multifinance yang dibangun dengan Go menggunakan arsitektur clean architecture.

## Struktur Direktori

```
.
├── cmd/                 # Entry point aplikasi
├── config/             # Konfigurasi aplikasi
├── database/           # File-file database (DDL dan DML)
│   ├── DDL.sql
│   └── DML.sql
├── delivery/           # Layer delivery (HTTP handlers)
│   ├── controller/    # HTTP controllers
│   └── dto/           # Data Transfer Objects
├── middleware/         # HTTP middlewares
├── model/              # Entity/domain models
│   └── dto/           # DTOs untuk model
├── repository/         # Layer repository (database access)
├── service/           # Business logic
├── usecase/           # Use case implementations
│   └── transaction/   # Use case khusus transaksi
└── utils/             # Utility functions
```

## Layer dan Alur Kerja

### 1. Model Layer

**Lokasi**: `/model`

Berisi struktur data yang merepresentasikan entitas bisnis:

- `customer.go`: Model untuk data nasabah
- `transaction.go`: Model untuk data transaksi
- `limit.go`: Model untuk limit kredit nasabah

Contoh model:
```go
type Transaction struct {
    ContractNumber string    `db:"contract_number"`
    CustomerNIK    string    `db:"customer_nik"`
    OTR            int64     `db:"otr"`
    AdminFee       int64     `db:"admin_fee"`
    Installment    int64     `db:"installment"`
    Interest       int64     `db:"interest"`
    AssetName      string    `db:"asset_name"`
    CreatedAt      time.Time `db:"created_at"`
}
```

### 2. Repository Layer

**Lokasi**: `/repository`

Bertanggung jawab untuk berinteraksi dengan database. Menggunakan `sqlx` untuk eksekusi query.

File penting:
- `customer_repo.go`: Operasi CRUD untuk nasabah
- `transaction_repo.go`: Operasi untuk transaksi
- `limit_repo.go`: Manajemen limit kredit
- `types.go`: Interface untuk database operations

Contoh method repository:
```go
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx repository.DBTx, t *model.Transaction) error {
    query := `INSERT INTO transactions (...) VALUES (?, ?, ...)`
    _, err := tx.ExecContext(ctx, query, ...)
    return err
}
```

### 3. Service Layer

**Lokasi**: `/service`

Mengimplementasikan business logic aplikasi. Layer ini menggunakan repository untuk mengakses data.

Contoh service (`transaction_service.go`):
- Validasi bisnis
- Mengkoordinasikan beberapa operasi repository
- Menangani transaction management

```go
func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, transaction *model.Transaction, tenor int) error {
    // 1. Validasi limit
    // 2. Kurangi limit
    // 3. Simpan transaksi
    // 4. Commit transaction
}
```

### 4. Use Case Layer

**Lokasi**: `/usecase`

Berisi implementasi use case spesifik yang mengkoordinasikan alur kerja bisnis yang lebih kompleks.

### 5. Delivery Layer (HTTP Handlers)

**Lokasi**: `/delivery/controller`

Menangani HTTP request dan response. Bertanggung jawab untuk:
- Binding dan validasi input
- Memanggil service layer
- Mengubah response ke format yang diinginkan
- Menangani error

Contoh handler (`transaction_handler.go`):

```go
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    // 1. Bind request body ke DTO
    // 2. Panggil service layer
    // 3. Handle response/error
}
```

## Alur Kerja: Membuat Transaksi Baru

1. **HTTP Request**
   ```
   POST /transactions
   {
       "customer_nik": "1234567890123456",
       "otr": 500000,
       "admin_fee": 50000,
       "tenor": 1,
       "asset_name": "Laptop XYZ"
   }
   ```

2. **Delivery Layer** (`transaction_handler.go`)
   - Bind request body ke DTO
   - Validasi input
   - Panggil service layer

3. **Service Layer** (`transaction_service.go`)
   - Mulai database transaction
   - Dapatkan limit nasabah
   - Validasi limit
   - Kurangi limit
   - Simpan transaksi
   - Commit transaction

4. **Repository Layer**
   - Eksekusi query ke database
   - Return hasil/error

5. **Response**
   ```json
   {
       "status": 201,
       "message": "Transaction created successfully",
       "data": {
           "contract_number": "TRX-123456",
           "status": "success"
       }
   }
   ```

## Database Schema

### Tabel `customers`
- Menyimpan data nasabah
- Primary Key: `nik`


### Tabel `customer_limits`
- Menyimpan limit kredit per nasabah per tenor
- Primary Key: (`customer_nik`, `tenor`)
- Foreign Key ke `customers(nik)`

### Tabel `transactions`
- Mencatat semua transaksi
- Primary Key: `contract_number`
- Foreign Key ke `customers(nik)`

## Error Handling

Aplikasi menggunakan error handling yang konsisten:
- Error dari database di-wrap dengan pesan yang bermakna
- Error response mengikuti format standar
- HTTP status code yang sesuai dengan jenis error

## Keamanan

- Validasi input di semua layer
- Menggunakan transaction untuk operasi yang memerlukan konsistensi data
- Menggunakan mutex untuk mencegah race condition

## Testing

Setiap layer seharusnya memiliki test yang sesuai:
- Unit test untuk repository
- Unit test untuk service (dengan mock repository)
- Integration test untuk API endpoints
