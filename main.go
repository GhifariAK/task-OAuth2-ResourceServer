// BE 2: Resource Server
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// 1. Membuat Middleware (AuthMiddleware)
// Fungsi ini akan mencegat setiap request yang masuk ke API.
// Tugasnya mengecek kualifikasi request sebelum diizinkan mengakses data utama.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// A. Mengambil isi dari header "Authorization"
		// Setiap request dari BE 3 WAJIB membawa header "Authorization".
		// Jika kosong, langsung hentikan proses (Status 401).
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Akses Ditolak: Header Authorization tidak ditemukan", http.StatusUnauthorized)
			return
		}

		// B. Memastikan format token bearer
		// Format standar Token adalah "Bearer <string_token_acak>".
		// pecah string berdasarkan spasi untuk mengambil tokennya saja.
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Akses Ditolak: Format token salah", http.StatusUnauthorized)
			return
		}
		tokenString := splitToken[1]

		// C. Bertanya ke Authorization Server apakah token ini valid
		// Karena API Resource  ini tidak punya database token, kita harus melakukan HTTP Request internal
		// ke Authorization Server untuk memverifikasi keaslian token.
		// dengan cara Membuat HTTP GET request ke localhost:9096/verify
		req, err := http.NewRequest("GET", "http://localhost:9096/verify", nil)
		if err != nil {
			log.Printf("Internal Error: Gagal membuat request ke Auth Server: %v", err)
			http.Error(w, "Terjadi kesalahan internal server", http.StatusInternalServerError)
			return
		}

		// kirimkan kembali token yang diterima dari client ke Authorization Server melalui header
		req.Header.Set("Authorization", "Bearer "+tokenString)

		// Lakukan request ke Authorization Server untuk memverifikasi token.
		client := &http.Client{}
		resp, err := client.Do(req)

		// Membaca response body dari Authorization Server untuk keperluan logging jika token tidak valid.
		bodyBytes, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		// D. Mengevaluasi status kebenaran token
		// Jika Authentication Server merespons selain Status OK (200), berarti token palsu, expired, atau salah.
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Otorisasi Gagal: Auth Server menolak token. Alasan: %s", string(bodyBytes))
			http.Error(w, "Akses Ditolak: Token tidak valid atau sudah expired", http.StatusUnauthorized)
			return
		}

		// E. Jika lolos semua pengecekan, panggil 'next.ServeHTTP' untuk masuk ke fungsi getReports.
		next.ServeHTTP(w, r)
	})
}

// 2. ENDPOINT: Resource Data
// Membuat Endpoint API (Data yang dilindungi)

// Fungsi ini hanya bisa diakses jika request sudah lolos dari AuthMiddleware.
func getReports(w http.ResponseWriter, r *http.Request) {
	// Contoh payload data rahasia perusahaan yang akan dikirim ke client jika token valid.
	data := map[string]interface{}{
		"status":  "success",
		"message": "Autentikasi Berhasil! Selamat datang di API Resource",
		"data": []map[string]string{
			{"id": "R001", "name": "Dummy Laporan Penjualan Q1", "date": "2026-03-31"},
			{"id": "R002", "name": "Dummy Laporan Audit Infrastruktur", "date": "2026-04-10"},
		},
	}

	// Mengatur response agar dalam format JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	// Mux berfungsi sebagai router
	mux := http.NewServeMux()

	// 3. Mendaftarkan Route dan Membungkusnya dengan Middleware
	// Endpoint /api/reports dilindungi oleh AuthMiddleware
	mux.Handle("/api/reports", AuthMiddleware(http.HandlerFunc(getReports)))

	// Server jalan di port 9097
	log.Println("BE 2 (Resource Server) berjalan di http://localhost:9097")
	log.Fatal(http.ListenAndServe(":9097", mux))
}
