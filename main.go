package main

import (
	"flag"
	"log"
	"github.com/gorilla/mux"
	"github.com/sevlyar/go-daemon"
	"net/http"
	"github.com/sunshine69/nagios-api-go/controllers"
	"github.com/sunshine69/nagios-api-go/conf"
)

func main() {

	Forking := flag.Bool("d", false, "Forking as deamon. Need to look into the config.json to set the log file and pid path correctly")
	flag.Parse()

	if *Forking {
		cntxt := &daemon.Context{
			PidFileName: conf.Config.PidFilePath,
			PidFilePerm: 0644,
			LogFileName: conf.Config.LogFilePath,
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[go-daemon nagios-api]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		log.Print("- - - - - - - - - - - - - - -")
		log.Print("daemon started")

		serveHTTP()
	} else {
		serveHTTP()
	}
}

func serveHTTP() {
	router := mux.NewRouter()

	router.HandleFunc("/{nagios_host}/service/{service_name}", controllers.GetServiceStatus).Methods("GET")

	// router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := conf.Config.Port
	if port == "" {
		port = "8000" //localhost
	}

	log.Printf("Port: %v\n", port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		log.Println(err)
	}
}