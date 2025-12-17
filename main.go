package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    // 命令行参数
    listenAddr := flag.String("l", "0.0.0.0:30000", "监听地址")
    flag.Parse()
    
    // 简单的HTTP服务器
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        hostname, _ := os.Hostname()
        fmt.Fprintf(w, "ech-workers is running!\n")
        fmt.Fprintf(w, "Hostname: %s\n", hostname)
        fmt.Fprintf(w, "Version: 1.0.0\n")
    })
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status": "running", "timestamp": "%s"}`, "2024-01-01T00:00:00Z")
    })
    
    log.Printf("Starting server on %s...\n", *listenAddr)
    log.Printf("Health check: http://%s/health\n", *listenAddr)
    
    if err := http.ListenAndServe(*listenAddr, nil); err != nil {
        log.Fatal("Server failed: ", err)
    }
}
