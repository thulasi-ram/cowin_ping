package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type App struct {
	app    fyne.App
	window fyne.Window

	headerLabelWidget       *widget.Label
	pincodeWidget           *widget.Entry
	dateWidget              *widget.Entry
	showOnlyAvailableWidget *widget.Check
	formWidget              *widget.Form
	dataTableContainer      *fyne.Container

	addDataWidget *widget.Button

	mainBox        *fyne.Container
	changingRunBox *fyne.Container

	running bool
}

func NewCowinApp() *App {
	fyneApp := app.New()
	window := fyneApp.NewWindow("Cowin Reminder")
	a := &App{app: fyneApp, window: window}
	a.setupWidgets()
	return a
}

func (a *App) setupWidgets() {
	a.headerLabelWidget = widget.NewLabel("Welcome to Cowin Reminder")

	a.pincodeWidget = PincodeWidget()
	a.pincodeWidget.SetText("517501")
	a.dateWidget = DateWidget()
	a.dateWidget.SetText("13-06-2021")

	a.showOnlyAvailableWidget = widget.NewCheck("", nil)

	ctx, cancel := context.WithCancel(context.Background())
	centersChan := make(chan VaccineCenters, 0)

	a.formWidget = &widget.Form{
		Items: FormItems{
			{Text: "pincode", Widget: a.pincodeWidget},
			{Text: "date", Widget: a.dateWidget},
			{Text: "only available", Widget: a.showOnlyAvailableWidget},
		},
		OnSubmit: func() {
			if a.running {
				return
			}
			a.running = true
			request := a.ToSearchRequest()
			go func(ctx context.Context) {
				PeriodicPushData(ctx, request, centersChan)
			}(ctx)

			go func(ctx context.Context) {
				for p := range centersChan {
					a.AddLog(GetFormattedDataAndMakeSound(ctx, p))
				}
			}(ctx)

			a.changingRunBox = a.newChangingRunBox()
			a.mainBox.Add(a.changingRunBox)
		},
		SubmitText: "start",
		OnCancel: func() {
			if !a.running {
				return
			}
			cancel()
			a.mainBox.Remove(a.changingRunBox)
			a.mainBox.Refresh()
			a.running = false
		},
		CancelText: "stop",
	}

	a.addDataWidget = widget.NewButton("Append", func() {
		val := fmt.Sprintf("Item")
		a.dataTableContainer.Add(canvas.NewText(val, color.White))
	})

	a.mainBox = container.NewVBox(
		a.headerLabelWidget,
		a.formWidget,
	)
}

func (a *App) Run() {
	a.window.SetContent(a.mainBox)
	a.window.Resize(fyne.NewSize(500, 500))
	a.window.ShowAndRun()
}

func (a *App) ToSearchRequest() *SearchRequest {
	return &SearchRequest{
		Pincode:             a.pincodeWidget.Text,
		Date:                a.dateWidget.Text,
		IsSecondDose:        false,
		IsFor45Plus:         false,
		OnlyShowIfAvailable: a.showOnlyAvailableWidget.Checked,
	}
}

func (a *App) newChangingRunBox() *fyne.Container {
	paddingTopContainer := container.NewGridWrap(fyne.Size{
		Width:  500,
		Height: 100,
	})
	a.dataTableContainer = container.NewVBox()
	return container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Running For :"+a.pincodeWidget.Text),
			widget.NewProgressBarInfinite(),
		),
		container.NewBorder(paddingTopContainer, nil, nil, nil,
			container.NewVScroll(a.dataTableContainer),
		))
}

func (a *App) AddLog(message string) {
	a.dataTableContainer.Add(canvas.NewText(message, color.White))
}

func (a *App) AddLogs(messages []string) {
	for _, m := range messages {
		a.AddLog(m)
	}
}
