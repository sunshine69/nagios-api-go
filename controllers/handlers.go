package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	app "github.com/sunshine69/nagios-api-go/app"
)

//GetServiceStatus - handler
var GetServiceStatus = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service_name"]
	nagiosHost := vars["nagios_host"]
	serviceStatus := app.GetServiceStatus(nagiosHost, serviceName)
	je := json.NewEncoder(w)
	je.SetIndent("", "    ")
	je.Encode(serviceStatus)
}