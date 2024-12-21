package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nerolizm/portone-payment/internal/config"
	"github.com/nerolizm/portone-payment/internal/handler"
	v1 "github.com/nerolizm/portone-payment/internal/infrastructure/http/v1"
	"github.com/nerolizm/portone-payment/internal/service"
	"github.com/rs/zerolog/log"
)

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Server is shutting down")

	// 서버 종료
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
	close(done)
}

func main() {
	// 로거 초기화
	config.InitLogger()
	log.Info().Msg("Logger initialized")

	// 환경 변수 초기화
	if err := config.Init(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize config")
	}
	log.Info().Msg("Config initialized")

	// HTTP 서버 설정
	mux := http.NewServeMux()

	// 서비스와 핸들러 초기화
	client := v1.NewClient()
	paymentService := service.NewPaymentService(client)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Swagger UI 설정
	opts := middleware.SwaggerUIOpts{
		BasePath: "/",
		SpecURL:  "/openapi.json",
		Path:     "/swagger",
	}
	sh := middleware.SwaggerUI(opts, nil)
	mux.Handle("/swagger", sh)
	mux.Handle("/openapi.json", http.FileServer(http.Dir("api")))

	// API 라우팅 설정
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/cancel-payment", paymentHandler.HandlePaymentCancel)

	server := &http.Server{
		Addr:    config.Env.Port,
		Handler: mux,
	}

	// 서버 종료 시그널을 처리하기 위한 채널
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	// 고루틴으로 종료 처리 실행
	go gracefulShutdown(server, quit, done)

	// 서버 시작
	log.Info().Str("port", config.Env.Port).Msg("Server is starting")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server failed to start")
	}

	// 서버가 완전히 종료될 때까지 대기
	<-done
}
