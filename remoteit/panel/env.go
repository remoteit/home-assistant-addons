package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// `remoteit-panel env <options.json>`: the add-on's options, and Home Assistant's own UI as Supervisor says it serves
// it, as shell assignments for the start script (run.sh) to eval — so the image needs no jq, and nothing runs on the
// target platform at build.
func envMain(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: remoteit-panel env <options.json>")
		os.Exit(2)
	}
	var o struct {
		RegistrationCode string `json:"registration_code"`
		DeviceName       string `json:"device_name"`
		Stage            string `json:"stage"`
		AnyPort          bool   `json:"any_port"`
	}
	if b, err := os.ReadFile(args[0]); err == nil {
		if err := json.Unmarshal(b, &o); err != nil {
			fmt.Fprintf(os.Stderr, "remoteit-panel: %s: %v\n", args[0], err)
		}
	}
	port, scheme := 8123, "http"
	if token := os.Getenv("SUPERVISOR_TOKEN"); token != "" {
		if p, ssl, err := coreInfo(token); err == nil {
			if p > 0 {
				port = p
			}
			if ssl {
				scheme = "https"
			}
		} else {
			fmt.Fprintf(os.Stderr, "remoteit-panel: Supervisor's /core/info: %v\n", err)
		}
	}
	stage := o.Stage
	if stage == "" {
		stage = "solo"
	}
	// What the device may do beside its services: no console (a shell in this container, on the host's network), no
	// proxy role; any port only when asked for. The subnet and exit node need a tun device the add-on does not have.
	policy := "console_control off;proxy_control off"
	if !o.AnyPort {
		policy += ";any_port_control off"
	}
	set := func(k, v string) { fmt.Printf("export %s=%s\n", k, quote(v)) }
	set("R3_DEVICE_STAGE", stage)
	set("R3_DEVICE_POLICY", policy)
	if o.DeviceName != "" {
		set("R3_DEVICE_NAME", o.DeviceName)
	}
	if o.RegistrationCode != "" {
		set("R3_REGISTRATION_CODE", o.RegistrationCode)
	}
	set("HA_PORT", fmt.Sprint(port))
	set("HA_SCHEME", scheme)
}

func coreInfo(token string) (port int, ssl bool, err error) {
	req, _ := http.NewRequest("GET", "http://supervisor/core/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, false, fmt.Errorf("%s", resp.Status)
	}
	var info struct {
		Data struct {
			Port int  `json:"port"`
			SSL  bool `json:"ssl"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, false, err
	}
	return info.Data.Port, info.Data.SSL, nil
}

// quote is a value as one shell word: single-quoted, its own single quotes closed and escaped.
func quote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }
