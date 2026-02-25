package test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	app "github.com/keypad/port/src/port"
)

func hold(t *testing.T) (net.Listener, int) {
	t.Helper()
	listen, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	port := listen.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return listen, port
}

func TestScan(t *testing.T) {
	listen, port := hold(t)
	defer func() {
		_ = listen.Close()
	}()
	start := port - 1
	if start < 1 {
		start = 1
	}
	text, err := app.Scan("127.0.0.1", start, port+1, 80*time.Millisecond)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if !strings.Contains(text, fmt.Sprintf("%d", port)) {
		t.Fatalf("missing port: %s", text)
	}
}

func TestCommon(t *testing.T) {
	text, err := app.Common("127.0.0.1", 20*time.Millisecond)
	if err != nil {
		t.Fatalf("common failed: %v", err)
	}
	if text == "" {
		t.Fatal("missing text")
	}
}

func TestHttp(t *testing.T) {
	listen, port := hold(t)
	defer func() {
		_ = listen.Close()
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handle)
	server := httptest.NewServer(mux)
	defer server.Close()
	url := fmt.Sprintf("%s/scan?host=127.0.0.1&start=%d&end=%d&timeout=80", server.URL, port, port)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if !strings.Contains(string(body), fmt.Sprintf("%d", port)) {
		t.Fatalf("missing port in response: %s", string(body))
	}
}
