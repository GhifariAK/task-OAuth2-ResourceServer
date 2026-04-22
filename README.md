# Authorization Server (OAuth 2.0)

Repository ini merupakan implementasi **Authorization Server** yang dikembangkan untuk simulasi arsitektur keamanan _Microservices_. Proyek ini berfokus pada protokol komunikasi **Machine-to-Machine (M2M)**.

## --Deskripsi Proyek

Server ini bertindak sebagai "Security Network". Tugas utamanya adalah memverifikasi identitas mesin atau aplikasi klien (Client) dan menerbitkan _Access Token_ yang aman, menggunakan alur **Client Credentials Grant** sesuai dengan spesifikasi **RFC 6749 Section 4.4**.

## Tech Stack

- **Language:** Go (Golang)
- **Library Utama:** `github.com/go-oauth2/oauth2/v4`
- **Database (Client Store):** PostgreSQL
- **Database (Token Store):** Redis / In-Memory
- **Router:** Standard Native `net/http`

## Standar Keamanan & Kepatuhan RFC

Sistem ini dibangun dengan mematuhi standar keamanan IETF secara ketat:

- **RFC 6749 Section 4.4:** Implementasi Client Credentials murni (Tanpa _Refresh Token_ untuk skenario B2B, namun ditambahkan pada code sesuai request).
- **RFC 6749 Section 2.3.1:** Memaksa otentikasi klien menggunakan _Header_ `Authorization: Basic <base64(id:secret)>`.
- **RFC 6750:** Menerbitkan token bertipe `Bearer`.
- **RFC 7662:** Menyediakan _endpoint_ validasi token untuk Resource Server.

## --Dokumentasi API

### 1. Generate Access Token

- **Endpoint:** `POST /token`
- **Tujuan:** Menukar Client ID dan Secret dengan Token.
- **Headers Wajib:**
  - `Content-Type: application/x-www-form-urlencoded`
  - `Authorization: Basic <base64_credentials>`
- **Body:** `grant_type=client_credentials`

### 2. validasi Token (Internal)

- **Endpoint:** `GET /verify`
- **Tujuan:** Digunakan oleh resource server untuk mengecek keaslian token.

## Cara Menjalankan Server

1. Pastikan PostgreSQL sudah ter-install dan berjalan.
2. Atur _environment variable_ untuk koneksi database:
   ```bash
   export DB_CONN="postgres://user:password@localhost:5432/namadb?sslmode=disable"
   ```
3. Jalankan aplikasi
   ```bash
    go run main.go
   ```
4. Server akan berjalan di http://localhost:9096.
