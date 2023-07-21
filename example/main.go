package main

import (
	"database/sql"
	"fmt"
	"github.com/araj-dev/muoauth/gen/zoomapi"
	"github.com/araj-dev/muoauth/zoom"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"os"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     endpoints.Zoom,
	}

	zoomclient, err := zoom.New(config, db)
	if err != nil {
		panic(err)
	}

	res, err := zoomclient.Me("134")

	fmt.Println(res.JSON200.Id)

	res2, err := zoomclient.CreateMeeting("134", zoomapi.MeetingCreateJSONRequestBody{
		DefaultPassword: Ptr(false),
		Type:            Ptr(2),
	})

	fmt.Println(res2.JSON201.CreatedAt)
}

func Ptr[T any](v T) *T {
	return &v
}
