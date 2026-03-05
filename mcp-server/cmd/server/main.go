package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/taskon/mcp-server/internal/api"
	"github.com/taskon/mcp-server/internal/llm"
	"github.com/taskon/mcp-server/internal/mcp"
)

func main() {
	// 命令行参数
	port := flag.String("port", "8080", "服务端口")
	claudeAPIKey := flag.String("claude-key", "", "Claude API Key")
	taskonAPIBase := flag.String("taskon-api", "https://api.taskon.xyz", "TaskOn API Base URL")
	debug := flag.Bool("debug", false, "启用调试模式")
	flag.Parse()

	// 配置日志
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg("TaskOn Quest MCP Server 启动中...")

	// 检查必要配置
	if *claudeAPIKey == "" {
		*claudeAPIKey = os.Getenv("CLAUDE_API_KEY")
	}
	if *claudeAPIKey == "" {
		log.Fatal().Msg("请设置 Claude API Key (--claude-key 或 CLAUDE_API_KEY 环境变量)")
	}

	// 初始化组件
	taskonClient := api.NewTaskOnClient(*taskonAPIBase)
	claudeClient := llm.NewClaudeClient(*claudeAPIKey)
	mcpServer := mcp.NewServer(taskonClient, claudeClient)

	// 配置HTTP路由
	mux := http.NewServeMux()

	// MCP WebSocket端点
	mux.HandleFunc("/mcp", mcpServer.HandleWebSocket)

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 对话API (REST接口，用于简单集成)
	mux.HandleFunc("/api/chat", mcpServer.HandleChat)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second, // LLM响应可能较慢
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info().Msg("收到关闭信号，正在优雅关闭...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("关闭服务器失败")
		}
	}()

	log.Info().Str("port", *port).Msg("服务器已启动")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("服务器启动失败")
	}

	log.Info().Msg("服务器已关闭")
}
