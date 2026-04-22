````markdown
# Resource Server (Protected API dummy)

Repositori ini merupakan **Resource Server** yang mensimulasikan penyimpanan data esensial perusahaan, seperti _Executive Dashboard_ dan laporan infrastruktur.

## --Deskripsi Proyek

Server ini bertindak sebagai penyedia data utama. Karena data yang disimpan bersifat rahasia, server ini dilindungi oleh lapisan **Auth Middleware**. Middleware ini bertugas mencegat setiap _request_ yang masuk dan memvalidasi "Tiket Akses" klien ke Authorization Server sebelum membuka gerbang data.

## Tech Stack

- **Language:** Go (Golang)
- **Router:** Standard Native `net/http`
- **Keamanan:** Custom Auth Middleware (Bearer Token Verification)
- **Komunikasi:** Internal HTTP Request (Introspection)

## Alur Keamanan (Middleware Workflow)

Setiap kali ada permintaan akses (misal dari aplikasi _Dashboard_ atau Postman), middleware akan melakukan langkah berikut:

1. Memeriksa keberadaan _Header_ `Authorization: Bearer <token>`.
2. Mengekstrak token dan menanyakan validitasnya langsung ke Authorization Server (BE 1) melalui HTTP `GET /verify`.
3. Jika Auth Server mengonfirmasi token tersebut valid dan aktif, barulah _handler_ utama dieksekusi dan data dikembalikan dalam format JSON.
4. Jika token palsu atau kedaluwarsa, akses langsung ditolak dengan status `401 Unauthorized`.

## --Dokumentasi API

### 1. Mengambil Data Laporan

- **Endpoint:** `GET /api/reports`
- **Tujuan:** Menarik data metrik Laporan Q1 dan Infrastruktur.
- **Headers Wajib:**
  - `Authorization: Bearer <access_token_dari_be_1>`
- **Respons Sukses (200 OK):** Mengembalikan struktur JSON berisi _array_ laporan.

## Cara Menjalankan Server

1. **Prasyarat Penting:** Pastikan [Authorization Server ](https://github.com/GhifariAK/task-OAuth2-AuthServer) sudah berjalan di port 9096.
2. Jalankan aplikasi ini:
   ```bash
   go run main.go
   ```
3. Server akan berjalan di http://localhost:9097.
````
