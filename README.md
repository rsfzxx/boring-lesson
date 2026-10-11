# go-starter

Starter backend Go dengan MongoDB, Redis, dan RabbitMQ. Baru berisi satu endpoint.

git clone <repo-anda>
cd boring-lesson
docker compose up -d --build
docker build -t usernameanda/boring-lesson-api:1.0 .
docker push usernameanda/boring-lesson-api:1.0

## Struktur

```
cmd/api/main.go        Titik masuk: config, koneksi, wiring, graceful shutdown
internal/config/       Konfigurasi dari environment variable
internal/infra/        Koneksi MongoDB, Redis, RabbitMQ
internal/dto/          Bentuk data request / response
internal/service/      Logika bisnis (tidak tahu soal HTTP)
internal/handler/      HTTP handler (tidak berisi logika bisnis)
internal/router/       Daftar route
```

Alur request: `router -> handler -> service -> infra`.

## Menjalankan

```bash
docker compose up -d
go mod tidy
go run ./cmd/api
curl -i localhost:8080/health
```

Response saat semua sehat (HTTP 200):

```json
{"status":"ok","services":{"mongodb":"up","rabbitmq":"up","redis":"up"}}
```

Coba `docker compose stop redis` lalu panggil lagi: status menjadi `degraded` dengan HTTP 503.

## Menambah fitur baru

1. Buat DTO di `internal/dto/`.
2. Buat service di `internal/service/`.
3. Buat handler di `internal/handler/`.
4. Daftarkan route di `internal/router/router.go`.
5. Rangkai di `cmd/api/main.go`.

Saat mulai menyimpan data, tambahkan `internal/model/` (bentuk dokumen MongoDB) dan
`internal/repository/` (query ke database), supaya service tidak memanggil driver langsung.
