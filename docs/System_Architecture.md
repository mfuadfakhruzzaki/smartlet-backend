1. Perangkat & koneksi lapangan

Node Server Lantai 1 (ikon Wi-Fi/walet).

Node Server Lantai 2.

Node Gateway (pengumpul data).

Koneksi antar-node: kedua Node Server terhubung ke Node Gateway lewat ESP-NOW (ditunjukkan dengan garis putus-putus). Artinya komunikasi jarak dekat, tanpa infrastruktur Wi-Fi/internet.

2. Aplikasi pengguna & antarmuka

Client memakai Mobile App.

Admin memakai Website.

Baik Mobile App maupun Website tersambung ke API Gateway (ikon biru di tengah) lewat koneksi IP (garis utuh). API Gateway adalah satu pintu masuk (reverse proxy/routing/terminasi TLS) untuk semua trafik dari aplikasi dan dari Node Gateway ke layanan di sisi server.

3. Layanan server (kanan)

Terdapat dua klaster layanan berwarna kuning yang berada di belakang API Gateway:

3.1 IoT Ingestion Gateway

MQTT Broker — pusat publish/subscribe untuk telemetri perangkat.

Ingestion Worker — konsumer dari broker yang memvalidasi, enrich, melakukan transformasi, dan menulis data ke basis data.

Klaster ini terhubung ke Database (lihat bagian 4). Jalur ini adalah data path utama untuk telemetri sensor/perangkat.

3.2 Backend Server (layanan bisnis)

Modul-modul layanan monolit/mikroservis yang melayani kebutuhan aplikasi:

Auth — otentikasi/otorisasi (mis. JWT).

User — profil & manajemen pengguna.

Farming — domain “rumah/gedung swiflet”, perangkat, sensor, panen, dsb.

Market — harga mingguan, listing/penjualan, transaksi.

Content — artikel, e-book, video, konten edukasi/promo.

Request — pengajuan pemasangan, perawatan, hingga pencopotan.
Backend ini juga terhubung ke Database.

4. Lapisan penyimpanan (Database)

Satu blok besar “Database” memuat tiga komponen penyimpanan:

TimescaleDB — time-series database (di atas PostgreSQL) untuk menyimpan telemetri/sensor (timestep, suhu, kelembapan, dsb) secara efisien dan bisa downsample/retention.

PostgreSQL — data relasional untuk entitas bisnis (user, rumah swiflet, artikel, harga, panen, order, dsb).

S3 Storage — objek/file (mis. foto bukti, gambar artikel/produk, video/thumbnail, arsip raw log).

IoT Ingestion Gateway menulis terutama ke TimescaleDB (dan bila perlu ke S3), sedangkan Backend Server membaca/menulis ke PostgreSQL/TimescaleDB/S3 sesuai modulnya.

5. Alur data utama

Telemetri perangkat (uplink)

Node Server (L1/L2) ➜ (ESP-NOW) ➜ Node Gateway ➜ API Gateway ➜ MQTT Broker ➜ Ingestion Worker ➜ TimescaleDB (opsional: arsip raw ke S3).

Aplikasi (Mobile/Web) kemudian membaca ringkasan/riwayat dari Backend Server yang mengambil data dari TimescaleDB.

Akses fitur aplikasi

Client/Admin ➜ (Mobile App/Website) ➜ API Gateway ➜ Backend Server (Auth, User, Farming, Market, Content, Request) ➜ PostgreSQL/TimescaleDB/S3.

Contoh: lihat harga mingguan (Market), input/lihat panen (Farming), baca artikel (Content), ajukan maintenance (Request).

Kontrol perangkat (downlink, jika diimplementasi)

Admin/otomasi ➜ Backend Server ➜ (publish perintah) MQTT Broker ➜ Node Gateway ➜ (protokol lokal) ➜ Node Server.

Catatan: jalur downlink ini logis dari komponen yang ada; detail protokol/pola QoS tidak ditampilkan di gambar, tetapi arsitektur mendukungnya.

6. Peran komponen & batas tanggung jawab

ESP-NOW segment: jaringan lokal antar perangkat (hemat daya, rendah latensi, tanpa infrastruktur).

API Gateway: keamanan (TLS, rate-limit), routing ke ingestion dan backend, titik audit.

IoT Ingestion: fokus streaming/ingest (skala tinggi, idempoten, schema validation), optimasi time-series.

Backend Server: API bisnis dan pengalaman pengguna (auth, konten, pasar, permintaan layanan).

Database tier: dipisah per karakteristik data—time-series vs relasional vs objek.
