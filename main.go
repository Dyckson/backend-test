package main

import (
	_ "backend-test/internal/cmd/server"
	"backend-test/internal/http/handler"
	"backend-test/internal/http/router"
)

func main() {
	r := router.NewRouter()
	handler.HandleRequests(r)
	r.Run(":1111")
}
