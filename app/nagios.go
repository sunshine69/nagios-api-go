package app

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"strings"
	"github.com/sunshine69/nagios-api-go/conf"
)

//NewNagiosData - Parse the status file and populate the data structure. map<nagios_type><uniq_key_to_lookup><key:val>
func NewNagiosData(nagiosStatusFilePath string) map[string]map[string]map[string]string {
	var output = make(map[string]map[string]map[string]string, 0)
	var curTokenStack []string
	//When we did not the unique key, we store info into this
	var tempMapStack []map[string]string
	//Part of unique key
	var uniqKey map[string]string
	var uniqKeyTypeLookup = map[string]map[string]bool {
		"hoststatus": {"host_name":false},
		"servicestatus": {"host_name":false, "service_description":false},
	}
	var foundUniqKey = false
	var ignoreBlock = false

	processLines := func(line string) {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "#") { return }
		// fmt.Printf("Line content: '%s'\n", line)
		if strings.HasSuffix(line, "{") {
			token := strings.Split(line, " ")
			if _, ok := uniqKeyTypeLookup[token[0]]; !ok {
				ignoreBlock = true
				return
			}
			// fmt.Printf("Beginning block: '%s'\n", token[0])
			if _, ok := output[token[0]]; ! ok {
				output[token[0]] = make(map[string]map[string]string, 0)
			}
			curTokenStack = append(curTokenStack, token[0])

		} else if strings.Contains(line, "=") && ! ignoreBlock {//We are in the middle of block
			// fmt.Printf("Middle block: '%s'\n", line)
			token := strings.Split(line, "=")
			//Always Populate tempMapStack for the current block
			if len(tempMapStack) == 0 {
				// fmt.Printf("Creating new entry in tempMapStack\n")
				var _tempMap = map[string]string {token[0]: token[1]}
				tempMapStack = append(tempMapStack, _tempMap)
			} else {
				//fmt.Printf("Update entry to tempMapStack\n")
				tempMapStack[0][token[0]] = token[1]
			}

			//Look for unique key if found then create it in the output
			if ! foundUniqKey {
				if _, ok := uniqKeyTypeLookup[curTokenStack[0]][token[0]]; ok {
					// fmt.Printf("Current token is in uniqKey map\n")
					switch curTokenStack[0] {
					case "hoststatus": //Found single key
						// fmt.Printf("Going to make new unique in output\n")
						uniqKey[token[0]] = token[1]
						output[curTokenStack[0]][uniqKey["host_name"]] = make(map[string]string, 0)
						foundUniqKey = true
					case "servicestatus":
						uniqKey[token[0]] = token[1]
						if len(uniqKey) == 2 {
							foundUniqKey = true
							// fmt.Printf("Going to make new unique in output\n")
							_key := fmt.Sprintf("%s-%s", uniqKey["host_name"], uniqKey["service_description"])
							output[curTokenStack[0]][_key] = make(map[string]string, 0)
						}
					}
				}
			}
		} else if line == "}" {// Closing block
			if ! ignoreBlock {
				// fmt.Printf("Closing stack for %s\n    tempMap: %v\ncurTokenStak: %v\nfoundUniqKey: %v\nuniqKey: %v\n", curTokenStack[0], tempMapStack, curTokenStack, foundUniqKey, uniqKey)
				switch curTokenStack[0] {
					case "hoststatus":
						output[curTokenStack[0]][uniqKey["host_name"]] = tempMapStack[0]
					case "servicestatus":
						_key := fmt.Sprintf("%s-%s", uniqKey["host_name"], uniqKey["service_description"])
						output[curTokenStack[0]][_key] = tempMapStack[0]
				}
			}
			// fmt.Printf("output now is: %v\n", output)
			//Reset state
			tempMapStack = tempMapStack[:0]
			curTokenStack = curTokenStack[:0]
			foundUniqKey = false
			uniqKey = make(map[string]string)
			ignoreBlock = false
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

// GetServiceStatus - Get service status
func GetServiceStatus(nagiosHost, serviceName string) map[string]string {
	n := NewNagiosData(conf.Config.NagiosStatusFilePath)
	_key := fmt.Sprintf("%s-%s", nagiosHost, serviceName)
	return n["servicestatus"][_key]
}