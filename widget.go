package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
	"time"
)

func headerWidget() *canvas.Text {
	return &canvas.Text{
		Alignment: fyne.TextAlignCenter,
		Color:     color.White,
		Text:      "Welcome to Cowin Reminder",
		TextSize:  18,
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}
}

// Form Items - Input Widgets

func newPincodeWidget() *widget.Entry {
	validator := func(p string) error {
		if p == "" {
			return errors.New("pincode cannot be empty")
		}
		if _, err := strconv.Atoi(p); err != nil {
			return errors.New("pincode should be a number")
		}
		if len(p) != 6 {
			return errors.New("pincode should be 6 digits")
		}

		return nil
	}
	input := widget.NewEntry()
	input.Validator = validator
	return input
}

func newDateWidget() *widget.Entry {

	validator := func(p string) error {
		if p == "" {
			return errors.New("date cannot be empty")
		}
		_, err := time.Parse("02-01-2006", p)
		if err != nil {
			return errors.New("date should be of format dd-mmm-yyyy")
		}

		return nil
	}
	input := widget.NewEntry()
	input.SetPlaceHolder("DD-MM-YYYY")
	input.Validator = validator
	return input
}

func newShowOnlyAvailableConfigWidget() *widget.Check {
	w := widget.NewCheck("", nil)
	w.Text = "Show only available slots"
	return w
}

func newIsSecondDoseConfigWidget() *widget.Check {
	w := widget.NewCheck("", nil)
	w.Text = "Is Second Dose?"
	return w
}

func newAgeGoupRadioConfigWidget() *widget.RadioGroup {
	w := widget.NewRadioGroup([]string{AgeGroup18Plus.Text, AgeGroup45Plus.Text}, nil)
	w.Horizontal = true
	return w
}

func newInfoPopUp(c fyne.Canvas) *widget.PopUp {
	var p *widget.PopUp
	text := canvas.NewText("Info", color.White)
	btn := widget.NewButton("Close", func() {
		p.Hide()
	})
	p = widget.NewModalPopUp(
		container.NewVBox(
			text,
			canvas.NewText("Contact: me@ahiravan.dev", color.White),
			btn,
		),
		c,
	)
	return p
}

func newHelpPopUp(c fyne.Canvas) *widget.PopUp {
	var p *widget.PopUp
	text := canvas.NewText("Help", color.White)
	btn := widget.NewButton("Close", func() {
		p.Hide()
	})
	p = widget.NewModalPopUp(container.NewVBox(
		text,
		canvas.NewText("Hits covin portal every minute", color.White),
		canvas.NewText("Makes an alarm if a slot is found", color.White),
		btn,
	), c)
	return p
}

func newToolBar(c fyne.Canvas) *widget.Toolbar {
	infoPopup := newInfoPopUp(c)
	helpPopup := newHelpPopUp(c)
	return widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			infoPopup.Show()
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			helpPopup.Show()
		}),
	)
}
