package app

import (
	"encoding/json"
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	o := NewNagiosData("/tmp/nagios.txt")
	_b, _ := json.MarshalIndent(o, "", "    ")
	log.Printf("%s", _b)
	srv := GetServiceStatus("xvt-aws-ansible", "check_errcd_wa_api_int")
	_b, _ = json.MarshalIndent(srv, "", "    ")
	log.Printf("%s", _b)
}