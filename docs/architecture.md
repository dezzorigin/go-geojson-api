# Arsitektur — go-geojson-api

Diagram: [`diagram.svg`](./diagram.svg)

## Gambaran

```
Client (curl / browser / aplikasi peta)
        │  HTTP, JSON
        ▼
Fiber v2 — router + middleware logger + CORS
        │
        ├── GET /                    → indeks API
        ├── GET /api/v1/health       → {"status":"healthy"}
        └── GET /api/v1/spatial/buildings → listBuildings()
                                                │
                                                ▼
                                     geojson.go: buildings()
                                     (data mock in-memory,
                                      2 gedung hardcode)
```

## Komponen

| File | Peran |
|---|---|
| `main.go` | Bootstrap: baca env `PORT` (default 8080), `app.Listen` |
| `app.go` | `newApp()`: Fiber app, middleware (logger, CORS), semua route |
| `geojson.go` | `buildings()`: data mock 2 gedung; `listBuildings()`: handler + filter `?type=` |
| `main_test.go` | Tes endpoint lewat `app.Test` (tanpa server sungguhan) |
| `Dockerfile` | Multi-stage: `golang:1.27-alpine` build → `alpine:3.20` runtime, user non-root |
| `landing/index.html` | Halaman dokumentasi statis (bukan bagian dari API) |

## Keputusan desain

- **Tanpa database.** Data hanya 2 gedung contoh dan tidak berubah;
  hardcode di `geojson.go` lebih sederhana dan jujur daripada
  berpura-pura ada database (README lama mengklaim ada database —
  tidak pernah ada).
- **Fiber v2.** Router cepat, middleware logger + CORS bawaan, `app.Test`
  memudahkan tes tanpa port sungguhan.
- **Filter di handler, bukan di query.** `listBuildings` menyaring slice
  in-memory berdasarkan `properties.type`. Filter tanpa hasil
  mengembalikan slice kosong (`[]fiber.Map{}`) supaya JSON-nya
  `"features": []`, bukan `null`.
- **Port lewat env.** `PORT` dibaca di `main.go` supaya container bisa
  memetakan port apa pun tanpa ubah kode.

## Alur request

1. Client `GET /api/v1/spatial/buildings?type=residential`.
2. Middleware logger mencatat request; CORS menambahkan header.
3. Router Fiber meneruskan ke `listBuildings`.
4. Handler mengambil `buildings()`, menyaring yang `type`-nya
   `residential`, mengembalikan `FeatureCollection` sebagai JSON.

## Batasan arsitektur

- Satu proses, satu instance; tidak ada state bersama antar instance
  (tidak relevan — datanya statis).
- Tidak ada cache: tiap request membangun ulang slice dari data hardcode
  (biayanya negligible untuk 2 feature).
- Tidak ada observability selain log request Fiber (metode, path, status,
  latensi).
