package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const minBrightness = 3
const maxBrightness = 100
const brightnessStep = 5

const minTemperature = 143
const maxTemperature = 344
const temperatureStep = 5

type keyLightOptions struct {
	NumberOfLights int `json:"numberOfLights"`
	Lights         []struct {
		On          int `json:"on"`
		Brightness  int `json:"brightness"`
		Temperature int `json:"temperature"`
	} `json:"lights"`
}

func main() {
	lightIPsStr := flag.String("light-ips", "192.168.0.181,192.168.0.182", "Comma-separated list of Elgato Key Light IPs.")
	cmd := flag.String("command", "toggle-power", "Command to run. May be: toggle-power, decrease-brightness, increase-brightness, decrease-temperature, increase-temperature, set-min-brightness, set-max-brightness, set-min-temperature, set-max-temperature, set-brightness, set-temperature.")
	value := flag.String("value", "", "Numeric value to use for 'set-brightness' and 'set-temperature' commands.")
	flag.Parse()

	lightIPs := strings.Split(*lightIPsStr, ",")

	for _, ip := range lightIPs {
		ip = strings.TrimSpace(ip)

		// Fetch current light options.
		lightURL := fmt.Sprintf("http://%s", net.JoinHostPort(ip, "9123"))
		resp, err := http.Get(lightURL + "/elgato/lights")
		if err != nil {
			log.Fatalln("Error fetching lights:", err)
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Error reading light info:", err)
		}
		var opts keyLightOptions
		if err = json.Unmarshal(buf, &opts); err != nil {
			log.Fatalln("Error unmarshaling options:", err)
		}

		// Apply command to the fetched options.
		switch *cmd {
		case "toggle-power":
			opts.Lights[0].On = 1 - opts.Lights[0].On
		case "decrease-brightness":
			opts.Lights[0].Brightness -= brightnessStep
			if opts.Lights[0].Brightness < minBrightness {
				opts.Lights[0].Brightness = minBrightness
			}
		case "increase-brightness":
			opts.Lights[0].Brightness += brightnessStep
			if opts.Lights[0].Brightness > maxBrightness {
				opts.Lights[0].Brightness = maxBrightness
			}
		case "increase-temperature":
			opts.Lights[0].Temperature -= temperatureStep
			if opts.Lights[0].Temperature < minTemperature {
				opts.Lights[0].Temperature = minTemperature
			}
		case "decrease-temperature":
			opts.Lights[0].Temperature += temperatureStep
			if opts.Lights[0].Temperature > maxTemperature {
				opts.Lights[0].Temperature = maxTemperature
			}
		case "set-min-brightness":
			opts.Lights[0].Brightness = minBrightness
		case "set-max-brightness":
			opts.Lights[0].Brightness = maxBrightness
		case "set-min-temperature":
			opts.Lights[0].Temperature = minTemperature
		case "set-max-temperature":
			opts.Lights[0].Temperature = maxTemperature
		case "set-brightness":
			val, err := strconv.Atoi(*value)
			if err != nil {
				log.Fatalf("Error parsing provided value %q: %v", *value, err)
			}
			opts.Lights[0].Brightness = val
		case "set-temperature":
			val, err := strconv.Atoi(*value)
			if err != nil {
				log.Fatalf("Error parsing provided value %q: %v", *value, err)
			}
			opts.Lights[0].Temperature = val
		default:
			log.Fatalf("Unknown command %q", *cmd)
		}

		// Set the new options.
		jsonOpts, err := json.Marshal(opts)
		if err != nil {
			log.Fatalln("Error marshaling JSON:", err)
		}

		req, err := http.NewRequest("PUT", lightURL+"/elgato/lights", bytes.NewReader(jsonOpts))
		if err != nil {
			log.Fatalln("Error building update request:", err)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln("Error updating light options:", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalln("Error updating light options:", resp.Status)
		}
	}
}
