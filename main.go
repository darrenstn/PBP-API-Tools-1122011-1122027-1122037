package main

import (
	"PBP-API-Tools-1122011-1122027-1122037/controllers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Inisialisasi Cron
	c := cron.New()

	config, errConfig := loadConfig("config.json")
	if errConfig != nil {
		log.Fatal("Error loading configuration:", errConfig)
	}

	// Menambahkan job Cron untuk mengirim email setiap 5 menit
	_, err := c.AddFunc("*/1 * * * *", func() {
		recipient := controllers.GetEmailWithContent("SILVER")
		err := sendEmail(config, recipient.Email, recipient.Content)
		if err != nil {
			log.Println("Gagal mengirim email:", err)
		} else {
			log.Println("Email terkirim pada", time.Now())
		}
		recipient2 := controllers.GetEmailWithContent("GOLD")
		err2 := sendEmail(config, recipient2.Email, recipient2.Content)
		if err2 != nil {
			log.Println("Gagal mengirim email:", err2)
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

type Config struct {
	SMTPHost     string `json:"SMTP_HOST"`
	SMTPPort     int    `json:"SMTP_PORT"`
	AuthEmail    string `json:"AUTH_EMAIL"`
	AuthPassword string `json:"AUTH_PASSWORD"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)

	if err != nil {
		return config, err
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// sendEmail mengirimkan email
func sendEmail(config Config, email string, content string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.AuthEmail)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Penawaran Spesial!")
	mailer.SetBody("text/html", content)

	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.AuthEmail,
		config.AuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
