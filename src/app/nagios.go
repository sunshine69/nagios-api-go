package app

import (
	"strconv"
	"bufio"
	"log"
	"os"
	"strings"
	"../conf"
)


//ServiceStatus - nagios ServiceStatus data structure. We only care about what we want for now
type ServiceStatus struct {
	Host_name string
	Service_description string
	Current_state int8
	Plugin_output string
	Long_plugin_output string
	Performance_data string
	Process_performance_data int8
	Is_flapping int8
}

//HostStatus - nagios hoststatus
type HostStatus struct {
	Host_name string
	Current_state int8
	Performance_data string
	Plugin_output string
	Long_plugin_output string
	Process_performance_data int8
}

//ServiceDowntime - nagios ServiceDowntime
type ServiceDowntime struct {
    Host_name string
    Service_description string
    downtime_id int
    comment_id int
    entry_time int64
    start_time int64
    flex_downtime_start int8
    end_time int64
    triggered_by int8
    fixed int8
    duration int16
    is_in_effect int8
    start_notification_sent int8
    comment string
}

//Nagios - Hold current state of the server
type Nagios struct {
	ServiceStatus []ServiceStatus
	HostStatus []HostStatus
	ServiceDowntime []ServiceDowntime
}

func parseInt(s string) int {
	o, e := strconv.Atoi(s)
	if e != nil {
		log.Printf("Can not parse %s", s)
		panic("parseInt")
	}
	return o
}

//NewNagios - Parse the status file and populate the data structure
func NewNagios(nagiosStatusFilePath string) Nagios {
	var output Nagios
	var curTokenStack []string
	var curServiceStatusStack []ServiceStatus
	var curHostStatusStack []HostStatus
	var curServiceDowntimeStack []ServiceDowntime

	processLines := func(line string) {
		line = strings.TrimSpace(line)
		if line == "" { return }
		// fmt.Printf("Line content: '%s'\n", line)
		if strings.HasSuffix(line, "{") {
			token := strings.Split(line, " ")
			// fmt.Printf("1 curTokenStack: %v\n", curTokenStack)
			// fmt.Printf("2 token: %v\n", token)
			if len(curTokenStack) > 0 {//Push the previous obect and clear things out
				// fmt.Printf("2 %v", curTokenStack)
				switch curTokenStack[0] {
					case "hoststatus":
						output.HostStatus = append(output.HostStatus, curHostStatusStack[0])
						curHostStatusStack = nil
					case "servicestatus":
						output.ServiceStatus = append(output.ServiceStatus, curServiceStatusStack[0])
						curServiceStatusStack = nil
					case "servicedowntime":
						output.ServiceDowntime = append(output.ServiceDowntime, curServiceDowntimeStack[0])
						curServiceDowntimeStack = nil
				}
				curTokenStack = nil
			}
			//Add new token holding the current one
			switch token[0] {
				case "hoststatus", "servicestatus", "servicedowntime":
					curTokenStack = append(curTokenStack, token[0])
					switch token[0] {
					case "hoststatus":
						curHostStatusStack = append(curHostStatusStack, HostStatus{})
					case "servicestatus":
						curServiceStatusStack = append(curServiceStatusStack, ServiceStatus{})
					case "servicedowntime":
						curServiceDowntimeStack = append(curServiceDowntimeStack, ServiceDowntime{})
					}
				default:
					// fmt.Printf("Not supported %s\n", token[0])
					return
			}
		} else if strings.Contains(line, "=") {
			if len(curTokenStack) == 0 { return }
			token := strings.Split(line, "=")
			switch curTokenStack[0] {
				case "hoststatus":
					switch token[0] {
						case "host_name":
							curHostStatusStack[0].Host_name = token[1]
						case "current_state":
							curHostStatusStack[0].Current_state = int8(parseInt(token[1]))
						case "performance_data":
							curHostStatusStack[0].Performance_data = token[1]
						case "plugin_output":
							curHostStatusStack[0].Plugin_output = token[1]
						case "long_plugin_output":
							curHostStatusStack[0].Long_plugin_output = token[1]
						case "process_performance_data":
							curHostStatusStack[0].Process_performance_data = int8(parseInt(token[1]))
						}

				case "servicestatus":
					switch token[0] {
						case "host_name":
							curServiceStatusStack[0].Host_name = token[1]
						case "service_description":
							curServiceStatusStack[0].Service_description = token[1]
						case "current_state":
							_state, _ := strconv.Atoi(token[1])
							curServiceStatusStack[0].Current_state = int8(_state)
						case "plugin_output":
							curServiceStatusStack[0].Plugin_output = token[1]
						case "long_plugin_output":
							curServiceStatusStack[0].Long_plugin_output = token[1]
						case "performance_data":
							curServiceStatusStack[0].Performance_data = token[1]
						case "process_performance_data":
							curServiceStatusStack[0].Process_performance_data = int8(parseInt(token[1]))
						case "is_flapping":
							curServiceStatusStack[0].Is_flapping = int8(parseInt(token[1]))
					}
				case "servicedowntime":
					switch token[0] {
					}
			}
		}
	}

	file, err := os.Open(nagiosStatusFilePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        processLines(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
	return output
}

//GetServiceStatus - Get service status
func GetServiceStatus(nagiosHost, serviceName string) ServiceStatus {
	n := NewNagios(conf.Config.NagiosStatusFilePath)
	for _, d := range(n.ServiceStatus) {
		if d.Host_name == nagiosHost && d.Service_description == serviceName {
			return d
		}
	}
	return ServiceStatus{}
}