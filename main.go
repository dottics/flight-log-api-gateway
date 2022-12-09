package main

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flight-log-api-gateway/src"
	"log"
	"net/http"
)

func main() {
	env := dutil.Env{}
	env.Load(".env")
	log.Println("Go API Gateway listening on port:", env.Vars["API_GW_PORT"])
	s := src.NewServer()

	log.Fatal(http.ListenAndServe(":"+env.Vars["API_GW_PORT"], s.Router))
}
