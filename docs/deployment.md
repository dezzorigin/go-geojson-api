# Deployment — go-geojson-api

## Menjalankan lokal

```bash
git clone https://github.com/dezzorigin/go-geojson-api.git
cd go-geojson-api
go run .
# http://localhost:8080
```

Port lain: `PORT=9000 go run .`

## Binary

```bash
go build -o go-geojson-api .
PORT=8080 ./go-geojson-api
```

Binary statis (CGO disabled di Dockerfile; `go build` lokal
menghasilkan binary dinamis biasa — cukup untuk Linux).

## Docker

```bash
docker build -t go-geojson-api .
docker run -p 8080:8080 go-geojson-api
```

Smoke test setelah jalan:

```bash
curl -s localhost:8080/api/v1/health
# {"status":"healthy"}
curl -s "localhost:8080/api/v1/spatial/buildings?type=residential" | grep -q bld-001 && echo OK
```

CI menjalankan langkah yang sama (`docker build` + smoke test) di
setiap push ke `main`.

## Halaman dokumentasi (Cloudflare Pages)

Halaman di `landing/index.html` adalah situs **statis** — dokumentasi dan
demo data contoh. Bukan API-nya (API tidak bisa jalan di Pages).

Deploy manual dengan wrangler (butuh token Cloudflare):

```bash
npx wrangler pages deploy landing --project-name go-geojson-api
```

URL live: https://go-geojson-api.pages.dev

> Nama project Pages adalah namespace global. Jika suatu saat nama
> `go-geojson-api` sudah dipakai akun lain, deploy akan menghasilkan
> subdomain bersufiks — saat itu terjadi, pakai nama fallback yang
> tercatat di PR deploy, jangan memaksa nama bersufiks.

## Yang TIDAK di-deploy

- **API live.** Tidak ada hosting API (VPS/Cloud Run/dll.) yang
  dipasang dalam repo ini. Untuk memakai API sungguhan, jalankan
  binary/Docker di server sendiri.
