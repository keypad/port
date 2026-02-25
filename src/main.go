package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/keypad/port/src/port"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, port.Help())
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Println(port.Help())
		return nil
	}
	if args[0] == "serve" {
		addr := ":4176"
		if len(args) > 1 {
			value, err := strconv.Atoi(args[1])
			if err != nil || value < 1 || value > 65535 {
				return fmt.Errorf("invalid port")
			}
			addr = fmt.Sprintf(":%d", value)
		}
		return port.Serve(addr)
	}
	if args[0] == "list" {
		fmt.Println(port.List())
		return nil
	}
	if args[0] == "common" {
		if len(args) < 2 || len(args) > 3 {
			return fmt.Errorf("usage: port common <host> [timeoutms]")
		}
		wait, err := readwait(args, 2)
		if err != nil {
			return err
		}
		text, err := port.Common(args[1], wait)
		if err != nil {
			return err
		}
		fmt.Print(text)
		return nil
	}
	if len(args) < 3 || len(args) > 4 {
		return fmt.Errorf("usage: port <host> <start> <end> [timeoutms]")
	}
	start, err := readnum(args[1])
	if err != nil {
		return err
	}
	end, err := readnum(args[2])
	if err != nil {
		return err
	}
	wait, err := readwait(args, 3)
	if err != nil {
		return err
	}
	text, err := port.Scan(args[0], start, end, wait)
	if err != nil {
		return err
	}
	fmt.Print(text)
	return nil
}

func readnum(text string) (int, error) {
	value, err := strconv.Atoi(text)
	if err != nil || value < 1 || value > 65535 {
		return 0, fmt.Errorf("invalid port range")
	}
	return value, nil
}

func readwait(args []string, at int) (time.Duration, error) {
	if len(args) <= at {
		return 180 * time.Millisecond, nil
	}
	value, err := strconv.Atoi(args[at])
	if err != nil || value < 10 || value > 10000 {
		return 0, fmt.Errorf("invalid timeout")
	}
	return time.Duration(value) * time.Millisecond, nil
}
