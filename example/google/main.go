package main

import (
	"database/sql"
	"fmt"
	"github.com/araj-dev/muoauth/store"
	"github.com/araj-dev/muoauth/tokenhttp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
	}
	s := store.NewTokenStore(config, db)
	client := tokenhttp.NewClient(s)

	// 1st
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/userinfo/v2/me", nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do("196", req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Printf("Status: %s\n", res.Status)
	// print body readble string
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Body: %s\n", string(body))

	// 2nd
	req2, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/userinfo/v2/me", nil)
	if err != nil {
		panic(err)
	}
	res2, err := client.Do("196", req2)
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()
	fmt.Printf("Status: %s\n", res2.Status)
	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Body: %s\n", string(body2))

	// 3rd
	res3, err := client.Get("196", "https://www.googleapis.com/userinfo/v2/me")
	if err != nil {
		panic(err)
	}
	defer res3.Body.Close()
	fmt.Printf("Status: %s\n", res3.Status)
	body3, err := io.ReadAll(res3.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Body: %s\n", string(body3))
}
