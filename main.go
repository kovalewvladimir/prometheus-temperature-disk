package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

type SmartctlOutput struct {
	Temperature struct {
		Current int `json:"current"`
	} `json:"temperature"`
}

var (
	devices   string
	deviceDir string
)

func init() {
	flag.StringVar(&devices, "devices", "", "Comma-separated list of devices to monitor (e.g., sda,sdb,sdc)")
	flag.StringVar(&deviceDir, "device-dir", "/dev", "Directory containing device files")
}

func metrics(w http.ResponseWriter, r *http.Request) {
	deviceList := strings.Split(devices, ",")
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP disk_temperature_celsius Current temperature of the disk\n")
	fmt.Fprintf(w, "# TYPE disk_temperature_celsius gauge\n")
	
	for _, device := range deviceList {
		device = strings.TrimSpace(device)
		devicePath := filepath.Join(deviceDir, device)
		cmd := exec.Command("smartctl", "-n", "standby", "-a", "-j", devicePath)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Error executing smartctl for device %s: %v", devicePath, err)
			continue
		}

		var data SmartctlOutput
		if err := json.Unmarshal(output, &data); err != nil {
			log.Printf("Error parsing JSON for device %s: %v", devicePath, err)
			continue
		}

		fmt.Fprintf(w, "disk_temperature_celsius{device=%q,path=%q} %d\n", device, devicePath, data.Temperature.Current)
	}
}

func main() {
	flag.Parse()
	
	if devices == "" {
		log.Fatal("Error: -devices parameter is required. Please specify comma-separated list of devices to monitor")
	}
	
	log.Printf("Starting server with devices: %s (in directory: %s)", devices, deviceDir)
	http.HandleFunc("/metrics", metrics)
	log.Fatal(http.ListenAndServe("0.0.0.0:9586", nil))
}
