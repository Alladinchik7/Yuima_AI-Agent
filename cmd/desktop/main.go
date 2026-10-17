package main

import (
	"context"

	"yuima-agent/internal/bootstrap"
	"yuima-agent/internal/cleaner"
	"yuima-agent/internal/config"
	"yuima-agent/internal/log"
	"yuima-agent/internal/reminder"
	"yuima-agent/internal/tools"
	"yuima-agent/internal/tools/remindertools"
	"yuima-agent/internal/tools/todo"
	"yuima-agent/internal/tools/weathertools"
	"yuima-agent/internal/tools/webtools"
	"yuima-agent/internal/ui"
)

func main() {
	// Создаем зависимости
	cfg := config.Load()
	logger := log.New(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := bootstrap.InitStore(cfg, logger)
	

	searchClient := bootstrap.InitSearchClient(cfg, logger)
	
	weatherClient := bootstrap.InitWeatherClient()
	wakeUp := make(chan struct{}, 1)

	// Создаем и запускаем фоновую очистку выполненых задач
	todoCleaner := cleaner.New(store, logger)
	go todoCleaner.Start(ctx)	

	// Проверка выбранного таймаута в конфиге
	timeOut := bootstrap.SetTimeOut(cfg) 

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

	llm := bootstrap.InitLLM(cfg)
	agentInst := bootstrap.InitAgent(store, llm, registry, logger, timeOut) 
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