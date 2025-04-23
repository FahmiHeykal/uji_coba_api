# Golang API dengan Registrasi Pengguna, Verifikasi Email, dan Reset Kata Sandi

Ini adalah API Go (Golang) sederhana yang menunjukkan fungsionalitas autentikasi pengguna dasar, termasuk :
- Registrasi pengguna dengan OTP untuk verifikasi email
- Verifikasi email dengan OTP
- Fungsionalitas lupa kata sandi dengan OTP untuk reset kata sandi
- Reset kata sandi menggunakan OTP

API ini menggunakan PostgreSQL sebagai basis data dan `JWT` untuk pembuatan token. Tujuan utama proyek ini adalah untuk menunjukkan bagaimana menangani autentikasi dan reset kata sandi dalam aplikasi web Go.

## Fitur
- **Registrasi Pengguna** : Memungkinkan pengguna baru untuk mendaftar dengan memberikan email dan kata sandi.
- **Verifikasi Email** : Mengirimkan OTP ke email pengguna untuk verifikasi email.
- **Lupa Kata Sandi** : Memungkinkan pengguna untuk mereset kata sandi mereka dengan memberikan email dan OTP.
- **Reset Kata Sandi** : Pengguna dapat mereset kata sandi mereka dengan memberikan OTP dan kata sandi baru.
