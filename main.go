package main

import "context"
import "fmt"
import "log"
import "net"
import "net/http"
import "net/http/httputil"
import "os"
import "strings"
import "time"

import "golang.org/x/crypto/acme/autocert"

type App struct {
	rundir   string
	baseHost string
}

func (app *App) hostPolicy(ctx context.Context, host string) error {
	if !strings.HasSuffix(host, app.baseHost) {
		return fmt.Errorf("not a subhost of %s: %s", app.baseHost, host)
	}

	sockpath := app.rundir + "/" + strings.TrimSuffix(host, app.baseHost)
	_, err := os.Stat(sockpath)
	if err != nil {
		return fmt.Errorf("not found: %s", host)
	}

	return nil
}

func (app *App) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	if !strings.HasSuffix(host, app.baseHost) {
		return nil, fmt.Errorf("not a subhost of %s: %s", app.baseHost, host)
	}

	sockpath := app.rundir + "/" + strings.TrimSuffix(host, app.baseHost)
	return net.Dial("unix", sockpath)
}

func main() {
	rundir := os.Getenv("RUNTIME_DIRECTORY")
	if rundir == "" {
		rundir = os.TempDir() + "/henk"
	}
	err := os.MkdirAll(rundir, 0775)
	if err != nil {
		panic(err)
	}

	statedir := os.Getenv("STATE_DIRECTORY")
	if statedir == "" {
		panic("bad STATE_DIRECTORY")
	}

	certCacheDir := statedir + "/autocert"
	err = os.MkdirAll(certCacheDir, 0700)
	if err != nil {
		panic(err)
	}

	var baseHost string
	if len(os.Args) > 1 {
		baseHost = os.Args[1]
	}
	app := App{rundir, baseHost}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(certCacheDir),
		HostPolicy: app.hostPolicy,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           app.dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = r.Host
		},
		Transport: transport,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("reverse proxy error: %v", err)
			http.Error(w, "not found", http.StatusNotFound)
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	http.Serve(certManager.Listener(), proxy)
}
