package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sunshine69/nagios-api-go/controllers"
	"net/http"
	"github.com/sunshine69/nagios-api-go/conf"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{nagios_host}/service/{service_name}", controllers.GetServiceStatus).Methods("GET")

	// router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := conf.Config.Port
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}