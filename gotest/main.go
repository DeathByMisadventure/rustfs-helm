package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	httpPort   = flag.Int("http", 9999, "HTTP POST port")
	tcpPort    = flag.Int("tcp", 0, "TCP port (0 to disable)")
	filePath   = flag.String("file", "", "File to tail (optional)")
	addTS      = flag.Bool("timestamp", false, "Prepend timestamp to each line (default: false)")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// HTTP server for POSTs (like MinIO audit webhook)
	if *httpPort > 0 {
		go startHTTPServer(ctx, *httpPort)
	}

	// TCP server (raw stream)
	if *tcpPort > 0 {
		go startTCPServer(ctx, *tcpPort)
	}

	// File tail
	if *filePath != "" {
		go tailFile(ctx, *filePath)
	}

	fmt.Println("Audit sidecar started - forwarding to stdout")

	<-sigCh
	fmt.Println("Shutting down...")
}

func printLine(line string) {
	if *addTS {
		ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		fmt.Printf("%s %s\n", ts, line)
	} else {
		fmt.Println(line)
	}
}



func startHTTPServer(ctx context.Context, port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			printLine(string(body))


		}
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP error: %v\n", err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(context.Background())
}

func startTCPServer(ctx context.Context, port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "TCP listen error: %v\n", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		go handleTCPConn(ctx, conn)
	}
}

func handleTCPConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		printLine(scanner.Text())

	}
}

func tailFile(ctx context.Context, path string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var offset int64 = 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f, err := os.Open(path)
			if err != nil {
				continue
			}
			if stat, _ := f.Stat(); stat.Size() > offset {
				f.Seek(offset, 0)
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					printLine(scanner.Text())
				}
				offset = stat.Size()
			}
			f.Close()
		}
	}
}

func isJSONLine(b []byte) bool {
	return len(b) > 0 && (b[0] == '{' || b[0] == '[')
}
