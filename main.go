package main

import (
	"flag"
	"fmt"
	"stockpricews/controller"
	"stockpricews/handler"
	"stockpricews/repository"
)

// In order to run the program the following params must be supplied
// go run . -server.port=8080 -db.user=<user> -db.pass=<pass> -db.port=8181
func main() {
	serverPort := flag.Int("server.port", 8080, "port to listen for incoming http requests")
	dbUser := flag.String("db.user", "root", "username to access the local mysql instance")
	dbPass := flag.String("db.pass", "", "password to access the local mysql instance")
	dbPort := flag.Int("db.port", 8181, "port of the local mysql instance")

	flag.Parse()

	// init and wire components following Onion Architecture. In a real-life app a DI framework might be used to do the job
	r, err := repository.New(*dbUser, *dbPass, *dbPort)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository %w", err))
	}
	c := controller.New(r)
	_, err = handler.New(c, *serverPort)
	if err != nil {
		fmt.Errorf("failed to initialize handler%w", err)
		panic(err)
	}

}
