# OAuth2 multi-user token client

This package provides a multi-user token client for OAuth2.
Based on [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2).

<!-- TOC -->
* [OAuth2 multi-user token client](#oauth2-multi-user-token-client)
  * [Install](#install)
  * [Motivation / Why not golang.org/x/oauth2](#motivation--why-not-golangorgxoauth2)
    * [Multi User](#multi-user)
    * [Store Refresh Token](#store-refresh-token)
  * [How to use](#how-to-use)
    * [Create TokenStore](#create-tokenstore)
    * [Create tokenhttp.Client](#create-tokenhttpclient)
    * [Use tokenhttp.Client](#use-tokenhttpclient)
  * [Zoom helper](#zoom-helper)
<!-- TOC -->

## Install
```bash
go get github.com/araj-dev/muoauth@latest
```

## Motivation / Why not [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
### Multi User
[golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) provides a token client for a single user.
```go
config := &oauth2.Config{
	ClientID: "clientid",
	//...
}
tokenSource := config.TokenSource(ctx, &oauth2.token{RefreshToken: "xxxxxx"})
client := oauth2.NewClient(ctx, tokenSource)

client.Do(...)
```

that client refresh token automatically when it is expired.
but it has a single oauth user in one client.
we don't want to create a client for each user.

So this package add id argument in http.Client functions to identify a user.
```go
db, err := sql.Open("mysql", os.Getenv("DSN"))
if err != nil {
panic(err)
}
defer db.Close()

oc := &oauth2.Config{
    ClientID: "clientid",
    //...
}
// Need oauth2 config and db connection.
// Database connection is used as persistent store.
s := store.NewTokenStore(oc, db)

// Pass store to tokenhttp.NewClient
// tokenhttp.NewClient will refresh token automatically.
// and it save token to database.
c := tokenhttp.NewClient(s)

req := http.NewRequest(...)

// pass client to identify a user
// identity can be any string, it is used to identify a which token re-use or re-load from persistent.
c.Do("identity", req)
```

### Store Refresh Token
This package automatically refresh token when it is expired and save it to persistent store.

`store/persistent.go` provides MySQL persistent implementation.

If you want to use other type persistent, you can implement `store.TokenStore` interface.

## How to use
### Create TokenStore
```go
db, err := sql.Open("mysql", os.Getenv("DSN"))
if err != nil {
panic(err)
}
defer db.Close()
oc := &oauth2.Config{
    ClientID: "clientid",
    //...
	endpoint: oauth2.Endpoint{
        TokenURL: "https://example.com/token",
    },// or use pre-defined endpoints in golang.org/x/oauth2 (e.g: google.Endpoint, endpoints.Zoom,...)
}
// Need oauth2 config and db connection.
// Database connection is used as persistent store.
s := store.NewTokenStore(oc, db)
```

### Create tokenhttp.Client
```go
c := tokenhttp.NewClient(s)
```

### Use tokenhttp.Client
```go
// use http package to create request
req := http.NewRequest(...)
// pass client to identify a user
// identity can be any string, it is used to identify a which token re-use or re-load from persistent.
c.Do("identity", req)
```

## Zoom helper

If you want to use Zoom API, you can use `zoom.NewClient` to create a client.
`zoomapi` package contains a generated go file from OpenAPI schema provided from [Zoom API](https://marketplace.zoom.us/docs/api-reference/zoom-api).
Open API Schema file is in `oas/oas.yaml`
Actually it is not a complete implementation of Zoom API, but it is enough for my use case.
If you want to add mode API, you can add implementation in zoom package. `zoom/zoom.go`

See `example/main.go` for more detail.

```go
// create zoom client
db, err := sql.Open("mysql", os.Getenv("DSN"))
if err != nil {
  panic(err)
}
defer db.Close()
oc := oauth2.Config{
  ClientID: "clientid",
  ClientSecret: "secret"
  Endpoint: endpoints.Zoom,
}
s := store.NewTokenStore(oc, db)

c := zoom.NewClient(s)

// self info endpoint
// and response is typed by generated golang struct.
res := c.Me("token-id")
```

## TODO:
[] Add more test
[] Add more API for Zoom
[] Add more example
[] Add MySQL table definition for persistent store
[] Add more persistent store implementation (e.g: Redis, Postgres,...)
[] Improve Store Interface (around Crypto)
[] Improve golang map race condition.