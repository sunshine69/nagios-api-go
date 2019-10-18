package app

import (
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	srv := GetServiceStatus("xvt-aws-ansible", "check_errcd_wa_api_int")
	log.Printf("%v", srv.Performance_data)
}