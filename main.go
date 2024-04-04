package main

import (
	"PBP-API-Tools-1122011-1122027-1122037/controllers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"

	"github.com/gorilla/mux"
)

const (
	CONFIG_SMTP_HOST     = "smtp.gmail.com"
	CONFIG_SMTP_PORT     = 587
	CONFIG_SENDER_NAME   = "if-22027@students.ithb.ac.id"
	CONFIG_AUTH_EMAIL    = "if-22027@students.ithb.ac.id"
	CONFIG_AUTH_PASSWORD = "Ithb-2023"
)

func main() {
	// Inisialisasi Cron
	c := cron.New()

	// Menambahkan job Cron untuk mengirim email setiap 5 menit
	_, err := c.AddFunc("*/1 * * * *", func() {
		err := sendEmail()
		if err != nil {
			log.Println("Gagal mengirim email:", err)
		} else {
			log.Println("Email terkirim pada", time.Now())
		}
	})
	if err != nil {
		log.Fatal("Tidak bisa menambahkan job Cron:", err)
	}

	// Mulai Cron
	c.Start()

	// Inisialisasi router HTTP
	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")

	// Memulai server HTTP
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

// sendEmail mengirimkan email
func sendEmail() error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", "martin.mada1134@gmail.com")
	mailer.SetHeader("Subject", "Penawaran Spesial!")
	mailer.SetBody("text/html", "Selamat datang di toko kami! Kami senang ingin memberikan penawaran spesial kepada Anda.<b>Dapatkan diskon 20%</b> untuk semua produk kami! Jangan lewatkan kesempatan ini untuk berbelanja dengan harga terbaik.Gunakan kode diskon <strong>'SPESIAL20'</strong> saat checkout untuk mengaktifkan penawaran ini.")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
