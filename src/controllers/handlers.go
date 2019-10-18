package controllers

import (
	// "log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	app "../app"
)

//GetServiceStatus - handler
var GetServiceStatus = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service_name"]
	nagiosHost := vars["nagios_host"]
	serviceStatus := app.GetServiceStatus(nagiosHost, serviceName)
	// output, err := json.MarshalIndent(serviceStatus, "", "    ")
	// log.Printf("%v", serviceStatus)
	je := json.NewEncoder(w)
	je.SetIndent("", "    ")
	je.Encode(serviceStatus)
}