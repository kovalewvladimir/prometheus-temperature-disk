package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var excludedDevicesList []string

type SmartctlOutput struct {
	Temperature struct {
		Current int `json:"current"`
	} `json:"temperature"`
}

func init() {
	if excluded := os.Getenv("EXCLUDED_DEVICES"); excluded != "" {
		excludedDevicesList = strings.Split(excluded, ",")
		for i, device := range excludedDevicesList {
			excludedDevicesList[i] = strings.TrimSpace(device)
		}
	}
}

func isExcludedDevice(name string) bool {
	for _, device := range excludedDevicesList {
		if device == name {
			return true
		}
	}
	return false
}

func isValidDiskDevice(name string) bool {
	if isExcludedDevice(name) {
		return false
	}

	// Проверяем тип устройства
	isSataDevice := strings.HasPrefix(name, "sd") && !strings.HasPrefix(name, "sdz")
	isNvmeDevice := strings.HasPrefix(name, "nvme")

	if isSataDevice {
		// Проверяем, что SATA устройство не является разделом (нет цифр в конце)
		hasNoDigitSuffix := !strings.ContainsAny(name[len(name)-1:], "0123456789")
		return hasNoDigitSuffix
	}

	return isNvmeDevice
}

func getDevices() []string {
	entries, err := os.ReadDir("/dev")
	if err != nil {
		log.Printf("Error reading /dev directory: %v", err)
		return nil
	}

	var devices []string
	for _, entry := range entries {
		name := entry.Name()
		if isValidDiskDevice(name) {
			devices = append(devices, name)
		}
	}
	return devices
}

func metrics(w http.ResponseWriter, r *http.Request) {
	deviceList := getDevices()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP disk_temperature_celsius Current temperature of the disk\n")
	fmt.Fprintf(w, "# TYPE disk_temperature_celsius gauge\n")

	for _, device := range deviceList {
		devicePath := filepath.Join("/dev", device)
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
	excludedDevices := os.Getenv("EXCLUDED_DEVICES")
	if excludedDevices != "" {
		log.Printf("Excluded devices: %s", excludedDevices)
	}

	devices := getDevices()
	if len(devices) == 0 {
		log.Fatal("Error: No devices found in /dev directory")
	}

	log.Printf("Starting server with detected devices: %s", strings.Join(devices, ", "))
	http.HandleFunc("/metrics", metrics)
	log.Fatal(http.ListenAndServe("0.0.0.0:9586", nil))
}
