package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"yuima-agent/internal/agent"
	"yuima-agent/internal/cleaner"
	"yuima-agent/internal/config"
	"yuima-agent/internal/llm"
	"yuima-agent/internal/reminder"
	"yuima-agent/internal/search"
	"yuima-agent/internal/storage"
	"yuima-agent/internal/storage/memory"
	"yuima-agent/internal/storage/sqlite"
	"yuima-agent/internal/tools"
	"yuima-agent/internal/tools/remindertools"
	"yuima-agent/internal/tools/todo"
	"yuima-agent/internal/tools/weathertools"
	"yuima-agent/internal/tools/webtools"
	"yuima-agent/internal/ui"
	"yuima-agent/internal/weather"
)

const (
	localTimeout = 5*time.Minute // Таймаут для локальной ЛЛМ
	cloudTimeout = 180*time.Second// Таймаут для облачной ЛЛМ
)

func main() {
	// Создаем зависимости
	cfg := config.Load()
	logger := newLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var store storage.Store
	if cfg.UseMemory {
		store = memory.NewStore()
	}else {
		path := filepath.Join("..", "..", "data.db")
		logger.Info("db.path", "path", path)
		
		db, err := sqlite.New(path, logger)
		if err != nil {
			logger.Error("db.error", "error", "failed create file db")
			os.Exit(1)
		}
		defer db.Close()
		
		store = sqlite.NewStore(db)
	}

	var searchClient webtools.SearchProvider
	if cfg.UseMockSearch {
		 searchClient = search.NewMock()
	}else {
		httpClient := http.Client{Timeout: 15*time.Second}
		searchClient = search.NewClient(&httpClient, cfg.SearxNGURL, cfg.FormatSearxNG, cfg.LangSerxNG, cfg.Safesearch)
		
	}
	
	weatherClient := weather.NewClient()
	wakeUp := make(chan struct{}, 1)

	// Создаем и запускаем фоновую очистку выполненых задач
	todoCleaner := cleaner.New(store, logger)
	go todoCleaner.Start(ctx)	

	// Проверка выбранного таймаута в конфиге
	var timeOut time.Duration
	if cfg.UseCloud {
		timeOut = cloudTimeout
	}else {
		timeOut = localTimeout
	}

	// Создаем реестр инструментов
	registry := tools.NewRegistry(
		todo.NewAddTool(store, 20),
		todo.NewListTool(store, 0),
		todo.NewCompleteTool(store, 10),
		todo.NewUncompleteTool(store, 5),
		todo.NewDeleteTool(store, 3),
		
		remindertools.NewListTool(store, 0),
		remindertools.NewSetTool(store, wakeUp, 10),

		weathertools.NewCheckTool(weatherClient, 0),

		webtools.NewSearchTool(searchClient, 0),
	)

	// Проверка и выбор ЛЛМ или Мок через конфиг
	// Создание агента с выбранным вариантом
	var agentInst *agent.Agent
	if cfg.UseMockAgent {
		mockClient := llm.NewMock()
		agentInst = agent.New(store, mockClient, registry, logger, timeOut)
	} else {
		llmClient := llm.NewClient(cfg.AIModel, cfg.AIAPIKey, cfg.AIBaseURL)
		agentInst = agent.New(store, llmClient, registry, logger, timeOut)
	}
	err := agentInst.LoadHistory()
	if err != nil {
		logger.Warn("main.load.history", "warning", "faeiled load history, error:" + err.Error())
	}

	go agentInst.CompressionHistory(ctx)

	// Создание приложения
	myApp := ui.NewApp(agentInst, logger)

	// Создание и запуск фонового процесса напоминаний
	remind := reminder.New(
		store,
		logger,
		wakeUp,
		func(text string) { // Callback функция для получения текста напоминания приложением
			myApp.ShowReminder(text)
		},
	)
	go remind.Start(ctx) // Фоновая горутина напоминаний
	
	myApp.Run(ctx) // Запуск приложения
}

// Функция создания логера
func newLogger(cfg config.Config) *slog.Logger {
	// Провера уровня по конфигу
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	// Выбор формата вывода по конфигу
	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}