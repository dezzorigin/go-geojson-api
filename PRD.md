# PRD — go-geojson-api

## 1. Informasi Proyek

| Field | Nilai |
|---|---|
| Nama Proyek | go-geojson-api |
| Tech Stack | Go 1.26+ (diuji 1.27), Fiber v2, tanpa database |
| Status | Berjalan (v1.1.0) |
| Penanggung Jawab | Nabil Aulia Dika |
| Versi Dokumen | 2.0 (ditulis ulang jujur, 2026-11-18) |
| Tanggal | 2026-11-18 |

> Catatan: PRD v1.0 adalah template generik yang tidak menggambarkan
> produk ini (menyebut "database", "npm install", dsb.). Dokumen ini
> menggantikannya dengan deskripsi yang sesuai kode yang sebenarnya.

## 2. Latar Belakang

Aplikasi peta 3D butuh data gedung dalam format GeoJSON. Repo ini
menyediakan API kecil yang mengembalikan FeatureCollection GeoJSON berisi
dua gedung contoh, supaya frontend peta bisa dikembangkan dan diuji tanpa
menunggu data survei atau backend penuh.

## 3. Masalah yang Diselesaikan

- Frontend peta 3D butuh endpoint GeoJSON yang bisa dipanggil langsung.
- Developer butuh contoh respons GeoJSON yang valid (FeatureCollection +
  Polygon + properties) untuk menguji rendering.
- Perlu health check sederhana untuk memastikan servis hidup.

## 4. Tujuan & Sasaran

- Menyajikan `GET /api/v1/spatial/buildings` sebagai GeoJSON valid (RFC 7946).
- Mendukung filter `?type=` (residential / academic).
- Filter tanpa hasil mengembalikan `"features": []`, bukan `null`.
- Health check `GET /api/v1/health` mengembalikan `{"status":"healthy"}`.
- Menjalankan di lokal (`go run`) dan di Docker tanpa konfigurasi rumit.

## 5. Ruang Lingkup

### In Scope
- 3 route: `/` (indeks), `/api/v1/health`, `/api/v1/spatial/buildings`.
- Filter `?type=` pada endpoint buildings.
- Port lewat env `PORT` (default 8080).
- Dockerfile multi-stage + CI (fmt, vet, test, build, Docker build +
  smoke test, HTML validate).
- Halaman dokumentasi statis di `landing/` (Cloudflare Pages).

### Out of Scope
- Database dan endpoint tulis (POST/PUT/DELETE) — tidak ada.
- Autentikasi, rate limiting, logging terstruktur — tidak ada.
- Data gedung nyata / hasil survei — data yang ada adalah mock contoh.
- Hosting API live — Pages hanya menayangkan dokumentasi statis.

## 6. Fitur & Persyaratan Fungsional

| ID | Fitur | Kriteria |
|---|---|---|
| F1 | Health check | `GET /api/v1/health` → 200, `{"status":"healthy"}` |
| F2 | Daftar gedung | `GET /api/v1/spatial/buildings` → FeatureCollection, 2 feature, id `bld-001` + `bld-002` |
| F3 | Filter type | `?type=residential` → 1 feature (`bld-001`); `?type=academic` → 1 feature (`bld-002`) |
| F4 | Filter tanpa hasil | `?type=<lainnya>` → 200, `"features": []` (bukan null) |
| F5 | Indeks API | `GET /` → 200, berisi nama, versi, daftar endpoint |
| F6 | Route tak dikenal | → HTTP 404 |
| F7 | Port konfigurabel | Env `PORT` dipakai; default 8080 |
| F8 | Docker | `docker build` berhasil; container menjawab health check |

Data gedung (mock, hardcode di `geojson.go`):

| id | name | type | height | capacity |
|---|---|---|---|---|
| bld-001 | Student Dormitory Block A | residential | 15 | 120 |
| bld-002 | Main Library | academic | 24 | 500 |

Koordinat: poligon contoh di kawasan Darussalam, Banda Aceh
(longitude 95.3697–95.3715, latitude 5.5681–5.5690).

## 7. Persyaratan Non-Fungsional

- **Performa:** respons < 10 ms untuk data statis sekecil ini (diuji lokal).
- **Keamanan:** tidak ada secret di repo (`.gitignore` menutup `.env`);
  tidak ada input pengguna yang diteruskan ke query — filter hanya
  membandingkan string.
- **Maintainability:** handler dan data terpisah dari bootstrap
  (`app.go`, `geojson.go`, `main.go`); tes di `main_test.go`.
- **Portabilitas:** Linux/macOS/Windows via Go; container via Docker.

## 8. Arsitektur & Tech Stack

Go + Fiber v2 (router, logger, CORS). Tanpa database: data mock
in-memory di `geojson.go`. Diagram: `docs/diagram.svg`.

## 9. Milestone / Roadmap

| Fase | Status |
|---|---|
| Inisialisasi & scaffolding | ✅ (v1.0.0) |
| Perbaikan filter null → `[]`, route indeks, PORT env | ✅ (v1.1.0, 2026-11-18) |
| Tes unit + Dockerfile + CI | ✅ (v1.1.0) |
| Dokumentasi + landing Pages | ✅ (v1.1.0) |
| Data gedung nyata | ⏳ belum diputuskan |
| Endpoint tulis / database | ⏳ belum diputuskan |

## 10. Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| Data mock dikira data nyata | README + PRD menyatakan mock secara eksplisit |
| CORS terbuka disalahgunakan | API read-only, tanpa data sensitif |
| Dependency usang | `go.mod` terkunci; audit via `go list -m -u` berkala |

## 11. Kriteria Keberhasilan

- Semua endpoint F1–F8 lolos play-test (tercatat di `demo-output.txt`).
- `go test -race ./...` hijau, coverage ≥ 70%.
- CI hijau di main (Go + Docker + landing).
- Halaman dokumentasi live di Cloudflare Pages, HTTP 200.

## 12. Lampiran

- Struktur: lihat README.md
- Referensi endpoint: `docs/api.md`
- Deployment: `docs/deployment.md`
- Lisensi: MIT (LICENSE)
