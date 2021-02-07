package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
	})
	e.GET("/oauth/line", func(c echo.Context) error {
		clientID := os.Getenv("CLIENT_ID")
		callbackURL := "http://localhost:1323"
		uri := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=code&state=xxx&client_id=%s&redirect_uri=%s&scope=profile openid", clientID, callbackURL)
		return c.Redirect(http.StatusPermanentRedirect, uri)
	})

	e.POST("/oauth/line/token", func(c echo.Context) error {
		var uri = "https://api.line.me/oauth2/v2.1/token"
		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("code", c.FormValue("code"))
		data.Set("redirect_uri", "http://localhost:1323")
		data.Set("client_id", os.Getenv("CLIENT_ID"))
		data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
		req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		/// http client end ///
		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)
		return c.JSON(resp.StatusCode, res)
	})
	e.POST("/oauth/line/getInfo", func(c echo.Context) error {
		uri := "https://api.line.me/v2/profile"
		req, err := http.NewRequest("GET", uri, nil)
		req.Header.Set("Authorization", "Bearer "+c.Request().Header.Get("Authorization"))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		/// http client end ///
		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)
		return c.JSON(resp.StatusCode, res)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
