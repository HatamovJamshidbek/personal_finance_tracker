package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	_ "github.com/redis/go-redis/v9"
)

func Email(email string) (string, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	source := rand.NewSource(time.Now().UnixNano())
	myRand := rand.New(source)

	randomNumber := myRand.Intn(900000) + 100000
	code := strconv.Itoa(randomNumber)
	err := client.Set(context.Background(), "Key-test", code, time.Minute*5).Err()
	if err != nil {
		return "", err
	}

	_, err = client.Get(context.Background(), "Key-test").Result()
	if err != nil {
		return "", err
	}

	err = SendCode(email, code)

	if err != nil {
		return "", err
	}

	return code, nil
}

func SendCode(email string, code string) error {
	from := "hatamovjamshid47@gmail.com"
	password := "rart megg xdym cqjk"

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("api/email/template.html")
	if err != nil {
		log.Fatal(err)
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Passwd string
	}{

		Passwd: code,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("Email send to:", email)
	return nil
}
