# Referensi API — go-geojson-api

Base URL saat dijalankan lokal: `http://localhost:8080`
(port bisa diubah lewat env `PORT`).

Semua respons berformat JSON. Tidak ada autentikasi.

## `GET /`

Indeks API.

```bash
$ curl -s localhost:8080/
{"endpoints":["GET /api/v1/health","GET /api/v1/spatial/buildings","GET /api/v1/spatial/buildings?type=residential","GET /api/v1/spatial/buildings?type=academic"],"name":"go-geojson-api","version":"1.1.0"}
```

| Field | Isi |
|---|---|
| `name` | `go-geojson-api` |
| `version` | versi API (semver) |
| `endpoints` | daftar endpoint yang tersedia |

## `GET /api/v1/health`

Health check.

```bash
$ curl -s localhost:8080/api/v1/health
{"status":"healthy"}
```

## `GET /api/v1/spatial/buildings`

FeatureCollection GeoJSON (RFC 7946) berisi gedung contoh.

```bash
$ curl -s localhost:8080/api/v1/spatial/buildings
{"features":[{"geometry":{"coordinates":[[[95.3697,5.5683],[95.3699,5.5683],[95.3699,5.5681],[95.3697,5.5681],[95.3697,5.5683]]],"type":"Polygon"},"properties":{"capacity":120,"height":15,"id":"bld-001","name":"Student Dormitory Block A","type":"residential"},"type":"Feature"},{"geometry":{"coordinates":[[[95.371,5.569],[95.3715,5.569],[95.3715,5.5685],[95.371,5.5685],[95.371,5.569]]],"type":"Polygon"},"properties":{"capacity":500,"height":24,"id":"bld-002","name":"Main Library","type":"academic"},"type":"Feature"}],"type":"FeatureCollection"}
```

Struktur tiap feature:

| Field | Tipe | Keterangan |
|---|---|---|
| `type` | string | selalu `"Feature"` |
| `properties.id` | string | `bld-001` / `bld-002` |
| `properties.name` | string | nama gedung |
| `properties.type` | string | `residential` / `academic` |
| `properties.height` | number | tinggi (angka contoh) |
| `properties.capacity` | number | kapasitas orang (angka contoh) |
| `geometry.type` | string | selalu `"Polygon"` |
| `geometry.coordinates` | number[][][] | `[longitude, latitude]`, ring tertutup |

### Query param `type`

Filter berdasarkan `properties.type`:

```bash
$ curl -s "localhost:8080/api/v1/spatial/buildings?type=residential"
{"features":[{"geometry":{"coordinates":[[[95.3697,5.5683],[95.3699,5.5683],[95.3699,5.5681],[95.3697,5.5681],[95.3697,5.5683]]],"type":"Polygon"},"properties":{"capacity":120,"height":15,"id":"bld-001","name":"Student Dormitory Block A","type":"residential"},"type":"Feature"}],"type":"FeatureCollection"}
```

Nilai `type` yang tidak cocok dengan gedung mana pun mengembalikan array
kosong — tetap GeoJSON valid:

```bash
$ curl -s "localhost:8080/api/v1/spatial/buildings?type=apartemen"
{"features":[],"type":"FeatureCollection"}
```

## Error

| Situasi | Respons |
|---|---|
| Route tidak dikenal | HTTP 404, body `Not Found` |
| Method salah (mis. POST ke endpoint GET) | HTTP 405, body `Method Not Allowed` |

## Batasan

- Data mock statis — selalu dua gedung yang sama.
- Tidak ada pagination (hanya 2 feature).
- Tidak ada endpoint tulis.
