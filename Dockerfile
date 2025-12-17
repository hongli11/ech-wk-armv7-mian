# ==================== 第一阶段：构建 ====================
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 配置 Go 环境变量和缓存
ARG TARGETARCH
ARG TARGETVARIANT
ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

# 先复制依赖文件，利用缓存层
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 复制源码
COPY . .

# 根据目标架构设置编译环境并编译
RUN if [ "$TARGETARCH" = "arm" ] && [ "$TARGETVARIANT" = "v7" ]; then \
        export GOARCH=arm GOARM=7; \
    else \
        export GOARCH=$TARGETARCH; \
    fi; \
    export GOOS=linux && \
    go build -ldflags="-s -w -extldflags '-static'" -trimpath -o ech-workers .

# ==================== 第二阶段：运行 ====================
FROM arm32v7/alpine:3.18

# 添加应用元数据
LABEL maintainer="your-name" \
      description="ech-workers for ARMv7" \
      org.opencontainers.image.source="https://github.com/hongli11/ech-wk-armv7-mian"

# 创建非root用户（安全优化）
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件，并设置正确的权限
COPY --from=builder --chown=appuser:appuser /app/ech-workers .

# 切换到非root用户运行
USER appuser

# 暴露端口
EXPOSE 30000

# 启动命令
ENTRYPOINT ["./ech-workers"]
CMD ["-l", "0.0.0.0:30000"]
