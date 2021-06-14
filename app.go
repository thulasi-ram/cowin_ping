package main

import (
	"bufio"
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type App struct {
	app    fyne.App
	window fyne.Window

	headerLabelWidget  *canvas.Text
	dataTableContainer *fyne.Container

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
	a.headerLabelWidget = headerWidget()
	ctx, cancel := context.WithCancel(context.Background())
	centersChan := make(chan VaccineCenters, 0)

	form := NewCowinForm()

	onSubmit := func() {
		if a.running {
			return
		}
		a.running = true
		request := ToSearchRequest(form)
		go func(ctx context.Context) {
			PeriodicPushData(ctx, request, centersChan)
		}(ctx)

		go func(ctx context.Context) {
			for p := range centersChan {
				buf := GetFormattedDataAndMakeSound(ctx, p)
				scanner := bufio.NewScanner(&buf)
				a.AddLog("As on: "+time.Now().Format(time.RFC1123), &fyne.TextStyle{Italic: true})
				for scanner.Scan() {
					a.AddLog(scanner.Text(), &fyne.TextStyle{Monospace: true})
				}
				a.AddLog("", nil)
				a.AddLog("", nil)
			}
		}(ctx)

		a.changingRunBox = a.newChangingRunBox(request)
		a.mainBox.Add(a.changingRunBox)
	}
	onCancel := func() {
		if !a.running {
			return
		}
		cancel()
		a.mainBox.Remove(a.changingRunBox)
		a.mainBox.Refresh()
		a.running = false
	}

	fyneForm := ToFyneForm(form, onSubmit, onCancel)

	a.mainBox = container.NewPadded(
		container.NewVBox(
			container.NewGridWithRows(2, a.headerLabelWidget, newToolBar(a.window.Canvas())),
			container.NewVBox(fyneForm),
		),
	)
}

func (a *App) Run() {
	a.window.SetContent(a.mainBox)
	a.window.Resize(fyne.NewSize(500, 600))
	a.window.ShowAndRun()
}

func (a *App) newChangingRunBox(r *SearchRequest) *fyne.Container {
	paddingTopContainer1 := container.NewGridWrap(fyne.Size{
		Width:  500,
		Height: 250,
	})
	paddingTopContainer2 := container.NewGridWrap(fyne.Size{
		Width:  500,
		Height: 100,
	})
	a.dataTableContainer = container.NewVBox()

	resizedDTC := container.NewVScroll(a.dataTableContainer)
	pgBar := widget.NewProgressBarInfinite()
	return container.NewBorder(
		paddingTopContainer1, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Running For :"+r.Pincode),
			pgBar,
		),
		container.NewMax(
			container.NewBorder(
				paddingTopContainer2, nil, nil, nil,
				resizedDTC,
			),
		),
	)
}

func (a *App) AddLog(message string, textStyle *fyne.TextStyle) {
	t := canvas.NewText(message, color.White)
	if textStyle != nil {
		t.TextStyle = *textStyle
	}
	a.dataTableContainer.Add(t)
}
