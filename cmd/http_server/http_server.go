package main

import (
	"context"
	"log"

	"github.com/aria3ppp/rag-server/cmd/http_server/router"
	ollama_adapter "github.com/aria3ppp/rag-server/internal/rag/infras/ollama"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"
	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	e := echo.New()

	ollamaLLM, err := ollama.New(ollama.WithServerURL("http://127.0.0.1:11434"))
	if err != nil {
		log.Fatalf("failed to create ollama llm: %s", err)
	}

	ollamaModel := "llama3.1"
	_, err = ollamaLLM.Call(context.Background(), "ping", llms.WithModel(ollamaModel))
	if err != nil {
		log.Fatal(err)
	}

	ollamaQuerier := ollama_adapter.NewQueryer(ollamaLLM, ollamaModel)
	_ = ollamaQuerier
	queryService := usecase.NewUseCase(nil, nil, nil, nil)
	queryServer := router.New(queryService)

	e.POST("/query", queryServer.HandleQuery)

	log.Fatal(e.Start(":8080"))
}
