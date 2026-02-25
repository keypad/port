```bash
> what is this?

 tcp port scanner utility.
 local cli + plain-text http output.
 built for terminal and curl workflows.

> features?

 ✓ scan custom port ranges on any host
 ✓ quick scan for common service ports
 ✓ plain-text http mode
 ✓ zero external services or credentials

> usage?

 port <host> <start> <end> [timeoutms]
 port common <host> [timeoutms]
 port list
 port serve [port]

> examples?

 go run ./src 127.0.0.1 1 1024
 go run ./src common 127.0.0.1
 curl "http://127.0.0.1:4176/scan?host=127.0.0.1&start=1&end=1024"
 curl "http://127.0.0.1:4176/common?host=127.0.0.1"

> stack?

 go 1.26 stdlib

> run?

 go run ./src 127.0.0.1 1 1024
 go run ./src serve 4176

> test?

 go test ./...

> proof?

 $ go test ./...
 ? github.com/keypad/port/src [no test files]
 ? github.com/keypad/port/src/port [no test files]
 ok github.com/keypad/port/test 0.221s

 $ go run ./src list | head -n 4
 20 ftp
 21 ftp
 22 ssh
 23 telnet

> links?

 https://github.com/keypad/port
```
