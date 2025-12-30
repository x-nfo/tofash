# 🧪 Test Add to Cart - Backend Connection

Dokumentasi untuk melakukan pengecekan koneksi fungsi add to cart dengan backend.

## 🔑 Cara Mendapatkan JWT Token

Sebelum menjalankan test, Anda memerlukan JWT token yang valid. Berikut adalah beberapa cara untuk mendapatkannya:

### Opsi 1: Menggunakan Halaman Login (Rekomendasi)

1. Buka file [`get_jwt_token.html`](get_jwt_token.html) di browser
2. Masukkan API URL: `http://localhost:8080/api/v1`
3. Masukkan email dan password user yang terdaftar
4. Klik tombol **"Login & Get Token"**
5. Token akan ditampilkan, klik **"Copy Token"** untuk menyalin
6. Paste token ke field "Access Token" di test cart

### Opsi 2: Menggunakan cURL

```bash
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"your_email@example.com","password":"your_password"}'
```

Response akan berisi:
```json
{
  "status": true,
  "message": "Success",
  "data": {
    "id": 1,
    "name": "User Name",
    "email": "your_email@example.com",
    "role": "customer",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

Salin nilai `access_token` dari response.

### Opsi 3: Menggunakan Postman atau API Client

1. Buat request baru dengan method `POST`
2. URL: `http://localhost:8080/api/v1/auth/signin`
3. Headers:
   - `Content-Type: application/json`
4. Body (raw JSON):
   ```json
   {
     "email": "your_email@example.com",
     "password": "your_password"
   }
   ```
5. Kirim request dan salin `access_token` dari response

### Opsi 4: Menggunakan Browser Dev Tools (Jika Sudah Login)

Jika Anda sudah login ke aplikasi, bisa mengambil token langsung dari cookies:

1. Buka aplikasi di browser
2. Tekan `F12` atau klik kanan > **Inspect** untuk membuka DevTools
3. Pilih tab **Application** atau **Storage**
4. Di sidebar, klik **Cookies**
5. Pilih domain aplikasi Anda
6. Cari cookie dengan nama `access_token`
7. Klik kanan pada `access_token` > **Copy Value**
8. Paste token ke field "Access Token" di test cart

**Catatan:** Ini adalah cara tercepat jika Anda sudah login ke aplikasi.

### Opsi 5: Membuat User Baru (Jika Belum Ada)

Jika Anda belum memiliki user, bisa mendaftar terlebih dahulu:

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "password_confirmation": "password123"
  }'
```

Setelah berhasil signup, gunakan Opsi 1, 2, atau 3 untuk login dan dapatkan token.

## 📋 Prasyarat

Sebelum menjalankan test, pastikan:

1. ✅ Backend server sudah berjalan di `http://localhost:8080`
2. ✅ Redis sudah berjalan
3. ✅ Memiliki JWT token yang valid (lihat cara mendapatkan di atas)
4. ✅ Ada produk yang tersedia di database

## 🚀 Cara Menjalankan Test

### Opsi 1: Browser Test (GUI)

1. Buka file [`test_add_to_cart.html`](test_add_to_cart.html) di browser
2. Konfigurasi:
   - **API URL**: `http://localhost:8080/api/v1` (default)
   - **Access Token**: Paste JWT token Anda di sini
   - **Product ID**: ID produk yang ingin ditest (default: 1)
   - **Quantity**: Jumlah produk (default: 1)
   - **Size**: Ukuran produk (default: M)
   - **Color**: Warna produk (default: Red)
3. Klik tombol **"Run All Tests"**
4. Lihat hasil test di halaman

### Opsi 2: Command Line Test (Node.js)

#### Tanpa Environment Variables:
```bash
node test_cart_connection.js
```

#### Dengan Environment Variables:
```bash
# Windows (Command Prompt)
set API_URL=http://localhost:8080/api/v1
set ACCESS_TOKEN=your_jwt_token_here
node test_cart_connection.js

# Windows (PowerShell)
$env:API_URL="http://localhost:8080/api/v1"
$env:ACCESS_TOKEN="your_jwt_token_here"
node test_cart_connection.js

# Linux/Mac
export API_URL=http://localhost:8080/api/v1
export ACCESS_TOKEN=your_jwt_token_here
node test_cart_connection.js
```

## 📊 Test Cases

### Test 1: Backend Connection
- Menguji apakah backend dapat diakses
- Endpoint: `GET /products?limit=1`

### Test 2: Add to Cart
- Menguji endpoint tambah ke keranjang
- Endpoint: `POST /carts`
- Payload:
  ```json
  {
    "product_id": 1,
    "quantity": 1,
    "size": "M",
    "color": "Red",
    "sku": "SKU-1-M-Red"
  }
  ```

### Test 3: Get Cart
- Menguji endpoint ambil keranjang
- Endpoint: `GET /carts`
- Mengembalikan semua item di keranjang user

### Test 4: Remove from Cart
- Menguji endpoint hapus dari keranjang
- Endpoint: `DELETE /carts?product_id=1`
- Menghapus item berdasarkan product_id

### Test 5: Clear Cart
- Menguji endpoint kosongkan keranjang
- Endpoint: `DELETE /carts/all`
- Menghapus semua item di keranjang

### Test 6: Full Workflow
- Menguji alur lengkap: Add → Get → Remove
- Memastikan semua endpoint bekerja bersama dengan baik

### Test 7: Multiple Items
- Menguji penambahan beberapa item sekaligus
- Memastikan keranjang dapat menampung multiple items

### Test 8: Update Quantity
- Menguji update quantity item yang sama
- Jika item yang sama ditambahkan, quantity harus bertambah
- Contoh: Add 1 item, lalu add lagi 2 item → Total quantity = 3

## 📈 Hasil Test

### Browser Test
- ✅ **Hijau**: Test passed
- ❌ **Merah**: Test failed
- 🟧 **Kuning**: Test pending
- 📊 **Summary**: Menampilkan total test, passed, failed, dan success rate

### Command Line Test
- ✅ **Hijau**: Test passed
- ❌ **Merah**: Test failed
- 🟡 **Kuning**: Informasi tambahan
- 📊 **Summary**: Menampilkan total test, passed, failed, dan success rate

## 🔧 Troubleshooting

### Error: "Backend connection failed"
- Pastikan backend server sudah berjalan
- Cek URL API yang benar
- Pastikan tidak ada firewall yang memblokir

### Error: "data token not found"
- Pastikan JWT token sudah di-set dengan benar
- Token harus valid dan belum expired
- Coba login ulang untuk mendapatkan token baru

### Error: "product not found"
- Pastikan product_id yang digunakan ada di database
- Cek endpoint `/products` untuk melihat daftar produk

### Error: "Cart is empty"
- Normal jika belum ada item di keranjang
- Lakukan test "Add to Cart" terlebih dahulu

## 📝 Contoh Output

### Browser Test:
```
🧪 Test Add to Cart - Backend Connection

⚙️ Configuration
API URL: http://localhost:8080/api/v1
Access Token: eyJhbGciOiJIUzI1NiIs...
Product ID: 1
Quantity: 1
Size: M
Color: Red

🚀 Run Tests

✅ Test 1: Backend Connection
   Details: {
     "message": "✅ Backend is accessible",
     "status": 200
   }

✅ Test 2: Add to Cart
   Details: {
     "message": "✅ Successfully added to cart",
     "status": 201
   }

...

📊 Test Summary
Total Tests: 8
Passed: 8
Failed: 0
Success Rate: 100.00%
```

### Command Line Test:
```
============================================================
🧪 CART API CONNECTION TEST SUITE
============================================================
API URL: http://localhost:8080/api/v1
Access Token: Set
============================================================

============================================================
TEST 1: Backend Connection
============================================================
✅ Backend is accessible
   Details: {
     "status": 200,
     "endpoint": "/products?limit=1"
   }

============================================================
TEST 2: Add to Cart
============================================================
✅ Add to cart endpoint works
   Details: {
     "status": 201,
     "message": "success",
     "data": {...}
   }

...

============================================================
📊 TEST SUMMARY
============================================================
Total Tests: 8
Passed: 8
Failed: 0
Success Rate: 100.00%
============================================================

🎉 All tests passed! Cart API is working correctly.
```

## 🎯 Kesimpulan

Jika semua test passed, berarti:
- ✅ Fungsi add to cart sudah terhubung dengan backend
- ✅ Backend dapat menerima dan memproses request cart
- ✅ Redis berfungsi dengan baik sebagai storage
- ✅ Semua endpoint cart bekerja dengan benar

Jika ada test yang failed, periksa:
- Backend server status
- Redis connection
- JWT token validity
- Product availability
- Network connection

## 📚 Referensi

- Frontend Store: [`storefront/src/lib/store.ts`](storefront/src/lib/store.ts)
- Backend Handler: [`internal/modules/product/handlers/cart_handler.go`](internal/modules/product/handlers/cart_handler.go)
- Backend Service: [`internal/modules/product/service/cart_service.go`](internal/modules/product/service/cart_service.go)
- Backend Repository: [`internal/modules/product/repository/cart_repository.go`](internal/modules/product/repository/cart_repository.go)
