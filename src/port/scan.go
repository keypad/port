package port

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Port int
	Name string
}

var names = map[int]string{
	20:   "ftp",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	123:  "ntp",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "smtp",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http",
}

func Help() string {
	return "usage: port <host> <start> <end> [timeoutms]\nusage: port common <host> [timeoutms]\nusage: port list\nusage: port serve [port]"
}

func List() string {
	keys := make([]int, 0, len(names))
	for key := range names {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	rows := make([]string, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, fmt.Sprintf("%d %s", key, names[key]))
	}
	return strings.Join(rows, "\n")
}

func Scan(host string, start int, end int, wait time.Duration) (string, error) {
	items, err := scan(host, start, end, wait)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "none\n", nil
	}
	rows := make([]string, 0, len(items))
	for _, item := range items {
		row := fmt.Sprintf("%d", item.Port)
		if item.Name != "" {
			row += " " + item.Name
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n") + "\n", nil
}

func Common(host string, wait time.Duration) (string, error) {
	keys := make([]int, 0, len(names))
	for key := range names {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	items, err := scanlist(host, keys, wait)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "none\n", nil
	}
	rows := make([]string, 0, len(items))
	for _, item := range items {
		row := fmt.Sprintf("%d", item.Port)
		if item.Name != "" {
			row += " " + item.Name
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n") + "\n", nil
}

func scan(host string, start int, end int, wait time.Duration) ([]Item, error) {
	if host == "" {
		return nil, fmt.Errorf("host required")
	}
	if start > end {
		return nil, fmt.Errorf("invalid range")
	}
	ports := make([]int, 0, end-start+1)
	for port := start; port <= end; port++ {
		ports = append(ports, port)
	}
	return scanlist(host, ports, wait)
}

func scanlist(host string, ports []int, wait time.Duration) ([]Item, error) {
	items := make([]Item, 0)
	lock := sync.Mutex{}
	group := sync.WaitGroup{}
	gate := make(chan struct{}, 256)
	for _, port := range ports {
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port range")
		}
		group.Add(1)
		go func(port int) {
			defer group.Done()
			gate <- struct{}{}
			defer func() {
				<-gate
			}()
			if !touch(host, port, wait) {
				return
			}
			item := Item{Port: port, Name: names[port]}
			lock.Lock()
			items = append(items, item)
			lock.Unlock()
		}(port)
	}
	group.Wait()
	sort.Slice(items, func(left int, right int) bool {
		return items[left].Port < items[right].Port
	})
	return items, nil
}

func touch(host string, port int, wait time.Duration) bool {
	join := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", join, wait)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
