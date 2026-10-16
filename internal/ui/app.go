package ui

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"yuima-agent/internal/agent"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	fyneApp fyne.App
	x      float32
	y      float32
	name   string
	agent *agent.Agent
	logger *slog.Logger
	
	messagesBox *fyne.Container
	messagesScroll *container.Scroll
	inputField *widget.Entry
	sendButton *widget.Button
	statusLabel *widget.Label
}

func NewApp(agentInst *agent.Agent, log *slog.Logger) *App {
	a := &App{
		fyneApp: app.New(),
		x: 800,
		y: 600,
		name: "Yuima AI-Agent",
		agent: agentInst,
		logger: log,
	}
	return a
}

func (a *App) Run(ctx context.Context) {
	window := a.createWindow()
	a.setupWidgets(ctx)
	content := a.toCompose()
	window.SetContent(content)
	window.ShowAndRun()
}

func (a *App) createWindow() fyne.Window {
	window := a.fyneApp.NewWindow(a.name)
	window.Resize(fyne.NewSize(a.x, a.y))
	return window
}

func (a *App) setupWidgets(ctx context.Context) {
	a.messagesBox = container.NewVBox()
	a.messagesScroll = container.NewVScroll(a.messagesBox)
	a.inputField = a.createInputField(ctx)
	a.sendButton = a.createButton(ctx)
	a.statusLabel = a.createStatusLabel(ctx)
}

func (a *App) toCompose() *fyne.Container {
	inputPanel := container.NewBorder(nil, nil, nil, a.sendButton, a.inputField)

	bottomPanel := container.NewVBox(a.statusLabel, inputPanel)

	content := container.NewBorder(nil, bottomPanel, nil, nil, a.messagesScroll)
	return content
}

func (a *App) createStatusLabel(ctx context.Context) *widget.Label {
	statusLabel := widget.NewLabel("Yuima думает...")
	statusLabel.Hide()
	return statusLabel
}

func (a *App) createInputField(ctx context.Context) *widget.Entry{
	inputField := widget.NewEntry()
	inputField.MultiLine = true
	inputField.SetMinRowsVisible(3)
	inputField.Wrapping = fyne.TextWrapWord

	inputField.SetPlaceHolder("Enter text...")

	inputField.OnSubmitted = func(text string) {
		a.sendMessage(ctx)
	}

	return inputField
}

func (a *App) createButton(ctx context.Context) *widget.Button {
	sendButton := widget.NewButton("Send", func() {
		a.sendMessage(ctx)
	})

	return sendButton
}

func (a *App) sendMessage(ctx context.Context) {
	if a.sendButton.Disabled() {
		labl := widget.NewLabel("Yuima думает, подождите...")
		labl.Wrapping = fyne.TextWrapWord
		a.messagesBox.Add(labl)
		a.messagesScroll.ScrollToBottom()
		return
	}
	
	text := strings.TrimSpace(a.inputField.Text)
	if text == "" {
		return
	}

	labl := widget.NewLabel("You: "+text)
	labl.Wrapping = fyne.TextWrapWord
	a.messagesBox.Add(labl)

	a.sendButton.Disable()
	done := make(chan struct{})
	
	a.showStatus("Yuima думает", done)
	
	go func() {
		defer close(done)
		
		resp, err := a.agent.Run(ctx, text)
		if err != nil {
			fyne.Do(func() {
				labl := widget.NewLabel("System Error: " + err.Error())
				a.logger.Info("system error", "error", err.Error())
				a.hideStatus()
				labl.Wrapping = fyne.TextWrapWord
				a.messagesBox.Add(labl)
				a.messagesScroll.ScrollToBottom()
				a.sendButton.Enable()
			})
			
			return
		}
		
		fyne.Do(func() {
			labl := widget.NewLabel("Yuima: " + resp)
			a.hideStatus()
			labl.Wrapping = fyne.TextWrapWord
			a.messagesBox.Add(labl)
			a.messagesScroll.ScrollToBottom()
			if a.sendButton.Disabled() {
				a.sendButton.Enable()
			}
		})
	}()
	a.statusLabel.Refresh()
	a.inputField.SetText("")
	a.messagesScroll.ScrollToBottom()
}

func (a *App) showStatus(text string, done <-chan struct{}) {
	go func() {
		i := 0
		select {
			case <-time.After(300*time.Millisecond):
				fyne.Do(func() {
					a.statusLabel.SetText(text)
					a.statusLabel.Show()
				})
				ticker := time.NewTicker(500 * time.Millisecond)
				for {
					select {
					case <-ticker.C:
						select {
						case <-done:
							ticker.Stop()
							return
						default:
							i = (i + 1) % 4
							dots := strings.Repeat(".", i)
							fyne.Do(func() {
								a.statusLabel.SetText(text + dots)
							})
						}
					case <-done:
						ticker.Stop()
						return
					}
				}
			case <-done:
				return
		}
	}()
}

func (a *App) hideStatus() {
	a.statusLabel.Hide()
}

func (a *App) ShowReminder(text string) {
	if a.messagesBox == nil {
		a.logger.Error("system error", "error", "window has not created")
		return
	}

	fyne.Do(func() {
		labl := widget.NewLabel("Reminder: "+text)
		labl.Wrapping = fyne.TextWrapWord
		a.messagesBox.Add(labl)
		a.messagesScroll.ScrollToBottom()
		a.sendButton.Enable()

		a.fyneApp.SendNotification(&fyne.Notification{Title: "Yuima", Content: text})
	})
}