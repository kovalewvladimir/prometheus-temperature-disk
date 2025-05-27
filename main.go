package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type SmartctlOutput struct {
	Temperature struct {
		Current int `json:"current"`
	} `json:"temperature"`
}

func metrics(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("smartctl", "-n", "standby", "-a", "-j", "/dev/sda")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing smartctl: %v", err), http.StatusInternalServerError)
		return
	}

	var data SmartctlOutput
	if err := json.Unmarshal(output, &data); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP disk_temperature_celsius Current temperature of the disk\n")
	fmt.Fprintf(w, "# TYPE disk_temperature_celsius gauge\n")
	fmt.Fprintf(w, "disk_temperature_celsius{device=\"sda\"} %d\n", data.Temperature.Current)
}

func main() {
	http.HandleFunc("/metrics", metrics)
	log.Fatal(http.ListenAndServe("0.0.0.0:9586", nil))
}
