package port

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handle)
	return http.ListenAndServe(addr, mux)
}

func Handle(write http.ResponseWriter, read *http.Request) {
	write.Header().Set("content-type", "text/plain; charset=utf-8")
	path := read.URL.Path
	if path == "/" {
		_, _ = write.Write([]byte(Help() + "\n"))
		return
	}
	if path == "/list" {
		_, _ = write.Write([]byte(List() + "\n"))
		return
	}
	host := read.URL.Query().Get("host")
	if host == "" {
		host = "127.0.0.1"
	}
	wait := 180 * time.Millisecond
	if raw := read.URL.Query().Get("timeout"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 10 || value > 10000 {
			write.WriteHeader(http.StatusBadRequest)
			_, _ = write.Write([]byte("invalid timeout\n"))
			return
		}
		wait = time.Duration(value) * time.Millisecond
	}
	if path == "/common" {
		text, err := Common(host, wait)
		if err != nil {
			write.WriteHeader(http.StatusBadRequest)
			_, _ = write.Write([]byte(err.Error() + "\n"))
			return
		}
		_, _ = write.Write([]byte(text))
		return
	}
	if path == "/scan" {
		start, err := readnum(read.URL.Query().Get("start"))
		if err != nil {
			write.WriteHeader(http.StatusBadRequest)
			_, _ = write.Write([]byte("invalid start\n"))
			return
		}
		end, err := readnum(read.URL.Query().Get("end"))
		if err != nil {
			write.WriteHeader(http.StatusBadRequest)
			_, _ = write.Write([]byte("invalid end\n"))
			return
		}
		text, err := Scan(host, start, end, wait)
		if err != nil {
			write.WriteHeader(http.StatusBadRequest)
			_, _ = write.Write([]byte(err.Error() + "\n"))
			return
		}
		_, _ = write.Write([]byte(text))
		return
	}
	write.WriteHeader(http.StatusNotFound)
	_, _ = write.Write([]byte(fmt.Sprintf("not found: %s\n", path)))
}

func readnum(text string) (int, error) {
	value, err := strconv.Atoi(text)
	if err != nil || value < 1 || value > 65535 {
		return 0, fmt.Errorf("invalid number")
	}
	return value, nil
}
