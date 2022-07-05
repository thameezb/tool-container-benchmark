package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Server()
}

type HTTPServer struct {
	engine           *gin.Engine
	runBenchmarkFunc func()
}

func NewServer(runBenchmarkFunc func()) Interface {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(
		gin.Recovery(),
	)
	engine.RedirectTrailingSlash = true

	r := &HTTPServer{
		engine:           engine,
		runBenchmarkFunc: runBenchmarkFunc,
	}

	route := engine.Group("/container-benchmark/")
	route.GET("ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
	route.POST("run", r.handleRunBenchmark)

	return r
}

func (h *HTTPServer) handleRunBenchmark(ctx *gin.Context) {
	h.runBenchmarkFunc()
	ctx.String(http.StatusCreated, "benchmark running")
}

func (h *HTTPServer) Server() {
	server := &http.Server{
		Addr:         ":3000",
		Handler:      h.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
