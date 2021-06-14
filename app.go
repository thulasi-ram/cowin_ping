package main

import (
	"bufio"
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type App struct {
	app    fyne.App
	window fyne.Window

	headerLabelWidget       *canvas.Text
	pincodeWidget           *widget.Entry
	dateWidget              *widget.Entry
	showOnlyAvailableWidget *widget.Check
	isSecondDoseWidget      *widget.Check
	formWidget              *widget.Form
	ageGoupRadioWidget      *widget.RadioGroup
	dataTableContainer      *fyne.Container

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
	a.headerLabelWidget = &canvas.Text{
		Alignment: fyne.TextAlignCenter,
		Color:     color.White,
		Text:      "Welcome to Cowin Reminder",
		TextSize:  18,
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}

	a.pincodeWidget = PincodeWidget()
	a.pincodeWidget.SetText("517501")
	a.dateWidget = DateWidget()
	a.dateWidget.SetText(time.Now().Format("02-01-2006"))

	a.showOnlyAvailableWidget = widget.NewCheck("", nil)
	a.showOnlyAvailableWidget.SetChecked(true)
	a.showOnlyAvailableWidget.Text = "Show only available slots"
	a.ageGoupRadioWidget = widget.NewRadioGroup([]string{AgeGroup18Plus.Text, AgeGroup45Plus.Text}, nil)
	a.ageGoupRadioWidget.SetSelected(AgeGroup18Plus.Text)
	a.ageGoupRadioWidget.Horizontal = true
	a.isSecondDoseWidget = widget.NewCheck("", nil)
	a.isSecondDoseWidget.Text = "Is Second Dose?"

	ctx, cancel := context.WithCancel(context.Background())
	centersChan := make(chan VaccineCenters, 0)

	a.formWidget = &widget.Form{
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

	layout1 := container.NewGridWithColumns(
		2,
		container.NewVBox(widget.NewLabel("Pincode"), a.pincodeWidget),
		container.NewVBox(widget.NewLabel("Date"), a.dateWidget),
	)
	a.formWidget.Append("", layout1)

	layout := container.NewHBox(
		a.ageGoupRadioWidget,
		widget.NewSeparator(),

		a.showOnlyAvailableWidget,
		widget.NewSeparator(),

		a.isSecondDoseWidget,
	)
	a.formWidget.Append("", layout)
	a.formWidget.Append("", canvas.NewText("", color.White))

	var infoPopup *widget.PopUp
	var helpPopup *widget.PopUp
	text := canvas.NewText("Info", color.White)
	btn := widget.NewButton("Close", func() {
		infoPopup.Hide()
	})
	text1 := canvas.NewText("Help", color.White)
	btn1 := widget.NewButton("Close", func() {
		helpPopup.Hide()
	})
	infoPopup = widget.NewModalPopUp(container.NewVBox(text, btn), a.window.Canvas())
	helpPopup = widget.NewModalPopUp(container.NewVBox(text1, btn1), a.window.Canvas())
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			infoPopup.Show()
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			helpPopup.Show()
		}),
	)

	a.mainBox = container.NewPadded(
		container.NewVBox(
			container.NewGridWithRows(2, a.headerLabelWidget, toolBar),
			container.NewVBox(
				a.formWidget,
			)),
	)
}

func (a *App) Run() {
	a.window.SetContent(a.mainBox)
	a.window.Resize(fyne.NewSize(500, 600))
	a.window.ShowAndRun()
}

func (a *App) ToSearchRequest() *SearchRequest {
	var isFor45Plus bool
	if a.ageGoupRadioWidget.Selected == AgeGroup45Plus.Text {
		isFor45Plus = true
	}
	return &SearchRequest{
		Pincode:             a.pincodeWidget.Text,
		Date:                a.dateWidget.Text,
		IsSecondDose:        a.isSecondDoseWidget.Checked,
		IsFor45Plus:         isFor45Plus,
		OnlyShowIfAvailable: a.showOnlyAvailableWidget.Checked,
	}
}

func (a *App) newChangingRunBox() *fyne.Container {
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
			widget.NewLabel("Running For :"+a.pincodeWidget.Text),
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
