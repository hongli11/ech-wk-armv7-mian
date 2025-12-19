package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

// 提取HTTP处理函数，增强可读性
func handleRoot(w http.ResponseWriter, r *http.Request) {
    hostname, err := os.Hostname()
    if err != nil {
        log.Printf("获取主机名失败: %v", err)
        hostname = "unknown" // 降级处理，避免响应失败
    }
    fmt.Fprintf(w, "ech-workers is running!\n")
    fmt.Fprintf(w, "Hostname: %s\n", hostname)
    fmt.Fprintf(w, "Version: 1.0.0\n")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("OK"))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // 动态生成当前时间戳（UTC格式）
    timestamp := time.Now().UTC().Format(time.RFC3339)
    fmt.Fprintf(w, `{"status": "running", "timestamp": "%s"}`, timestamp)
}

func main() {
    // 命令行参数（默认值）
    defaultAddr := "0.0.0.0:30000"
    // 支持从环境变量读取配置（容器化部署更友好）
    if envPort := os.Getenv("PORT"); envPort != "" {
        defaultAddr = "0.0.0.0:" + envPort
    }

    // 命令行参数优先级高于环境变量
    listenAddr := flag.String("l", defaultAddr, "监听地址 (格式: 0.0.0.0:端口)")
    flag.Parse()

    // 注册路由
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/health", handleHealth)
    http.HandleFunc("/status", handleStatus)

    log.Printf("服务器启动中，监听地址: %s\n", *listenAddr)
    log.Printf("健康检查地址: http://%s/health\n", *listenAddr)

    // 启动服务（增强错误日志）
    if err := http.ListenAndServe(*listenAddr, nil); err != nil {
        log.Fatalf("服务器启动失败: %v", err)
    }
}
