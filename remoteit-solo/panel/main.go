// The add-on's panel in Home Assistant's sidebar, served through ingress: the device's name, state and UID, a link to
// it in the portal, and — while it is not registered — a field for its registration code.
//
// Home Assistant's ingress rules: the panel answers only Supervisor's ingress proxy (172.30.32.2), and its page uses
// relative URLs (it is served under /api/hassio_ingress/<token>/). Its one write, the registration code, takes JSON
// with a header of its own, which a page on another origin cannot send without a preflight the panel never answers.
//
// What it shows is the daemon's own account: its control channel's status (connectd control.go), else the status file
// beside the key; the claim code and a refused registration's reason from the files the daemon and the entrypoint keep
// there. The panel holds no credential and calls no remote.it API.
package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//go:embed index.html
var page []byte

type panel struct {
	state    string // the key's directory
	control  string // the daemon's control socket
	codeFile string // where a code entered here goes, for the entrypoint to take
	stage    string
	allow    map[string]bool
	// haPort and haScheme are Home Assistant's own web UI, as Supervisor says it serves it: the service to add.
	haPort   int
	haScheme string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "env" {
		envMain(os.Args[2:])
		return
	}
	addr := flag.String("listen", ":29190", "address to serve the panel on (the add-on's ingress_port)")
	state := flag.String("state", "/data/remoteit-device", "the device's key directory")
	control := flag.String("control", "/run/remoteit-device/control.sock", "the daemon's control socket")
	codeFile := flag.String("code-file", "/data/registration_code", "where a registration code entered here is written")
	stage := flag.String("stage", "prod", "the device's stage, for the portal's address (run.sh gives the add-on's)")
	haPort := flag.Int("ha-port", 8123, "Home Assistant's web UI port (Supervisor's /core/info)")
	haScheme := flag.String("ha-scheme", "http", "http, or https when Home Assistant serves TLS (Supervisor's /core/info)")
	allow := flag.String("allow", "172.30.32.2", "the only peers answered, comma-separated (Supervisor's ingress proxy)")
	flag.Parse()

	p := &panel{state: *state, control: *control, codeFile: *codeFile, stage: *stage, allow: map[string]bool{},
		haPort: *haPort, haScheme: *haScheme}
	for _, a := range strings.Split(*allow, ",") {
		if a = strings.TrimSpace(a); a != "" {
			p.allow[a] = true
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", p.index)
	mux.HandleFunc("GET /api/status", p.status)
	mux.HandleFunc("POST /api/register", p.register)
	srv := &http.Server{Addr: *addr, Handler: p.guard(mux), ReadHeaderTimeout: 10 * time.Second}
	log.Printf("panel: on %s, answering %s", *addr, *allow)
	log.Fatal(srv.ListenAndServe())
}

// guard answers Supervisor's ingress proxy alone.
func (p *panel) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		if !p.allow[host] {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func (p *panel) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'self'")
	w.Write(page)
}

type deviceView struct {
	Registered bool    `json:"registered"`
	UID        string  `json:"uid,omitempty"`
	Name       string  `json:"name,omitempty"`
	Owner      string  `json:"owner,omitempty"`
	State      string  `json:"state"`
	Connected  bool    `json:"connected"`
	Since      *string `json:"since,omitempty"`
	Stage      string  `json:"stage"`
	Connectd   string  `json:"connectd,omitempty"`
	Package    string  `json:"package,omitempty"`
	Claim      string  `json:"claim,omitempty"`
	Portal     string  `json:"portal"`
	// Waiting is a code entered here and not yet taken; Error the reason the last registration was refused.
	Waiting bool   `json:"waiting"`
	Error   string `json:"error,omitempty"`
	// HomeAssistant is the service for Home Assistant's own UI: on this machine (the add-on is on the host's network).
	HomeAssistant struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Scheme string `json:"scheme"`
	} `json:"homeAssistant"`
}

func (p *panel) portal() string {
	if p.stage == "prod" {
		return "https://app.remote.it"
	}
	return "https://app." + p.stage + ".remote.it"
}

func (p *panel) status(w http.ResponseWriter, r *http.Request) {
	v := deviceView{State: "starting", Stage: p.stage, Portal: p.portal()}
	v.HomeAssistant.Host, v.HomeAssistant.Port, v.HomeAssistant.Scheme = "127.0.0.1", p.haPort, p.haScheme
	key := filepath.Join(p.state, "device.key")
	if b, err := os.ReadFile(key + ".uid"); err == nil {
		v.Registered, v.UID = true, strings.TrimSpace(string(b))
		v.Portal += "/#/devices/" + v.UID
	}
	if st, err := p.controlStatus(); err == nil {
		v.Connected = st.Connected
		v.Since = st.Since
		v.UID = orElse(st.Device.UID, v.UID)
		v.Name, v.Owner = deref(st.Device.Name), deref(st.Device.Owner)
		v.Connectd, v.Package = st.Device.Connectd, deref(st.Device.Package)
		if st.Connected {
			v.State = "connected"
		} else {
			v.State = "connecting"
		}
	} else if b, err := os.ReadFile(key + ".status"); err == nil {
		var f struct {
			UID, State, Version, Package string
		}
		if json.Unmarshal(b, &f) == nil {
			v.State, v.Connected, v.Connectd, v.Package = f.State, f.State == "connected", f.Version, f.Package
		}
	} else if !v.Registered {
		v.State = "not registered"
	}
	if b, err := os.ReadFile(key + ".claim"); err == nil {
		v.Claim = strings.TrimSpace(string(b))
	}
	if !v.Registered {
		if b, err := os.ReadFile(filepath.Join(p.state, "registration-error")); err == nil {
			v.Error = strings.TrimSpace(string(b))
		}
		if _, err := os.Stat(p.codeFile); err == nil {
			v.Waiting = true
		}
	}
	writeJSON(w, http.StatusOK, v)
}

type controlStatus struct {
	Device struct {
		UID      string  `json:"uid"`
		Name     *string `json:"name"`
		Owner    *string `json:"owner"`
		Connectd string  `json:"connectd"`
		Package  *string `json:"package"`
	} `json:"device"`
	Connected bool    `json:"connected"`
	Since     *string `json:"since"`
}

// controlStatus asks the daemon (connectd's control channel: a JSON request a line, its answer a line).
func (p *panel) controlStatus() (*controlStatus, error) {
	conn, err := net.DialTimeout("unix", p.control, time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte(`{"id":1,"method":"status"}` + "\n")); err != nil {
		return nil, err
	}
	lines := bufio.NewScanner(conn)
	lines.Buffer(make([]byte, 64<<10), 4<<20)
	if !lines.Scan() {
		return nil, errors.New("no answer")
	}
	var answer struct {
		Result *controlStatus `json:"result"`
		Error  string         `json:"error"`
	}
	if err := json.Unmarshal(lines.Bytes(), &answer); err != nil {
		return nil, err
	}
	if answer.Error != "" || answer.Result == nil {
		return nil, errors.New(answer.Error)
	}
	return answer.Result, nil
}

var codeShape = regexp.MustCompile(`^[A-Za-z0-9-]{6,128}$`)

func (p *panel) register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Remoteit-Panel") != "1" || !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not from the panel"})
		return
	}
	if _, err := os.Stat(filepath.Join(p.state, "device.key.uid")); err == nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "this device is registered already"})
		return
	}
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "send {\"code\": …}"})
		return
	}
	code := strings.TrimSpace(body.Code)
	if !codeShape.MatchString(code) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "that is not a registration code"})
		return
	}
	tmp := p.codeFile + ".new"
	if err := os.WriteFile(tmp, []byte(code+"\n"), 0o600); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "cannot keep the code: " + err.Error()})
		return
	}
	if err := os.Rename(tmp, p.codeFile); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "cannot keep the code: " + err.Error()})
		return
	}
	os.Remove(filepath.Join(p.state, "registration-error"))
	log.Printf("panel: a registration code was entered")
	writeJSON(w, http.StatusAccepted, map[string]bool{"accepted": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func orElse(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
