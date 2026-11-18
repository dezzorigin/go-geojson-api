# Changelog

Semua perubahan penting pada proyek ini dicatat di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
dan proyek ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-11-18

### Added
- Route indeks `GET /` yang mengembalikan nama, versi, dan daftar endpoint
- Tes unit endpoint di `main_test.go` (health, daftar gedung, filter,
  filter tanpa hasil, indeks, 404) — dijalankan dengan `go test -race`
- `Dockerfile` multi-stage (Go 1.27 → Alpine, user non-root)
- CI GitHub Actions: `gofmt` + `go vet` + `go test -race` + `go build` +
  `docker build` + smoke test container + `html-validate` untuk landing
- `.gitignore` Go
- Halaman dokumentasi statis di `landing/` (deploy Cloudflare Pages)
- Dokumentasi: `docs/api.md`, `docs/architecture.md`,
  `docs/deployment.md`, `docs/diagram.svg`
- `demo-output.txt`: output curl verbatim dari play-test endpoint

### Changed
- README dan PRD ditulis ulang secara jujur: README lama berisi instruksi
  `npm install` untuk proyek Go dan mengklaim ada database — keduanya
  tidak pernah ada
- Port server kini lewat env `PORT` (default 8080)
- Kode dipecah: `main.go` (bootstrap), `app.go` (route), `geojson.go`
  (data + handler) — sebelumnya semua di `main.go`

### Fixed
- Filter `?type=` yang tidak cocok gedung mana pun sebelumnya
  mengembalikan `"features": null` (GeoJSON tidak valid); sekarang
  mengembalikan `"features": []`

### Removed
- Klaim fiktif di CHANGELOG lama (rate limiting, Docker, validasi SQL
  injection, "optimasi query database") — tidak satu pun pernah ada di kode

## [1.0.0] - 2026-11-05

### Added
- Endpoint `GET /api/v1/health`
- Endpoint `GET /api/v1/spatial/buildings` dengan filter `?type=`
- Data mock GeoJSON 2 gedung (asrama + perpustakaan)
- README, PRD, LICENSE MIT
