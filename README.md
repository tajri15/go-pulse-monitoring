# Go-Pulse: Real-Time Uptime Monitoring System 🚀

Go-Pulse adalah sistem monitoring uptime website yang dibangun dengan Go (Golang) untuk backend dan Vue.js untuk frontend. Aplikasi ini memungkinkan pengguna untuk memantau status (Up/Down) dan waktu respons dari situs web mereka secara real-time melalui dashboard interaktif.

Proyek ini mendemonstrasikan arsitektur full-stack modern, memanfaatkan kekuatan konkurensi Go untuk *health checks* di latar belakang dan WebSockets untuk komunikasi real-time.


*(Tips: Ambil screenshot dashboard Anda yang sudah berjalan dan letakkan di sini!)*

---

### ## ✨ Fitur Utama

* **Autentikasi Pengguna**: Sistem registrasi dan login yang aman menggunakan JWT (JSON Web Tokens).
* **Manajemen Situs**: Pengguna dapat melakukan operasi CRUD (Create, Read, Delete) untuk situs-situs yang ingin mereka pantau.
* **Dashboard Real-Time**: Status situs diperbarui secara langsung di antarmuka pengguna tanpa perlu refresh, berkat WebSockets.
* **Worker Konkuren**: Pemeriksaan kesehatan situs dilakukan di latar belakang secara bersamaan menggunakan *worker pool* (Goroutines), memastikan performa yang efisien.
* **Containerization**: Seluruh aplikasi (backend Go & database PostgreSQL) di-container-kan menggunakan Docker untuk kemudahan development dan deployment.

---

### ## 🛠️ Tumpukan Teknologi (Tech Stack)

| Kategori      | Teknologi                                        |
|---------------|----------------------------------------------------|
| **Backend** | Go (Golang), Gin Web Framework, Gorilla WebSocket  |
| **Frontend** | Vue.js, Vite, Tailwind CSS                         |
| **Database** | PostgreSQL                                         |
| **Infrastruktur** | Docker, Docker Compose                           |

---

### ## 🚀 Cara Menjalankan Proyek

Untuk menjalankan proyek ini di lingkungan lokal, ikuti langkah-langkah berikut.

**Prasyarat:**
* [Git](https://git-scm.com/)
* [Go](https://go.dev/dl/) (versi 1.21+)
* [Node.js & npm](https://nodejs.org/) (versi 18+)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/)

**Langkah-langkah Instalasi:**

1.  **Clone repositori ini:**
    ```bash
    git clone [https://github.com/](https://github.com/)<USERNAME_ANDA>/go-pulse-monitoring.git
    cd go-pulse-monitoring
    ```

2.  **Pastikan Docker Desktop sedang berjalan.**

3.  **Jalankan Backend & Database:**
    Buka terminal di direktori root proyek dan jalankan:
    ```bash
    docker-compose up --build
    ```
    Backend akan berjalan di `http://localhost:8080`.

4.  **Jalankan Frontend:**
    Buka terminal **kedua**, navigasikan ke folder frontend, instal dependensi, dan jalankan server development:
    ```bash
    cd web/frontend-vue
    npm install
    npm run dev
    ```
    Frontend akan dapat diakses di `http://localhost:5173`.

5.  Buka browser Anda dan kunjungi **`http://localhost:5173`** untuk mulai menggunakan aplikasi!

---

### ## 📁 Struktur Proyek