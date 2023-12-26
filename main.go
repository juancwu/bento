package main

import (
	"fmt"
	"net/http"

	"github.com/juancwu/bento/shared"
	"github.com/labstack/echo/v4"
)

func main() {
    p := shared.PORT

    e := echo.New()

    e.GET("/health-check", func(c echo.Context) error {
        return c.String(http.StatusOK, http.StatusText(http.StatusOK))
    })

    fmt.Printf("Health Check: http://127.0.0.1:%d/health-check\n", p)
    e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", p)))
}
