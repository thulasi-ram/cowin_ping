package main

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type CowinForm struct {
	pincodeWidget           *widget.Entry
	dateWidget              *widget.Entry
	showOnlyAvailableWidget *widget.Check
	isSecondDoseWidget      *widget.Check
	ageGoupRadioWidget      *widget.RadioGroup

	//form *widget.Form
}

func NewCowinForm() *CowinForm {
	f := &CowinForm{
		pincodeWidget:           newPincodeWidget(),
		dateWidget:              newDateWidget(),
		showOnlyAvailableWidget: newShowOnlyAvailableConfigWidget(),
		isSecondDoseWidget:      newIsSecondDoseConfigWidget(),
		ageGoupRadioWidget:      newAgeGoupRadioConfigWidget(),
	}
	f.setDefaults()
	return f
}

func (f *CowinForm) setDefaults() {
	f.pincodeWidget.SetText("517501")
	f.dateWidget.SetText(time.Now().Format("02-01-2006"))
	f.ageGoupRadioWidget.SetSelected(AgeGroup18Plus.Text)
	f.showOnlyAvailableWidget.SetChecked(true)
}

func ToFyneForm(f *CowinForm, onSubmitCallBack func(), onCancelCallback func()) *widget.Form {
	form := &widget.Form{
		SubmitText: "start",
		CancelText: "stop",
		OnSubmit:   onSubmitCallBack,
		OnCancel:   onCancelCallback,
	}

	mainInputLayout := container.NewGridWithColumns(
		2,
		container.NewVBox(widget.NewLabel("Pincode"), f.pincodeWidget),
		container.NewVBox(widget.NewLabel("Date"), f.dateWidget),
	)

	configLayout := container.NewHBox(
		f.ageGoupRadioWidget,
		widget.NewSeparator(),
		f.showOnlyAvailableWidget,
		widget.NewSeparator(),
		f.isSecondDoseWidget,
	)

	form.Append("", mainInputLayout)
	form.Append("", configLayout)
	form.Append("", canvas.NewText("", color.White))

	return form
}

func ToSearchRequest(f *CowinForm) *SearchRequest {
	var isFor45Plus bool
	if f.ageGoupRadioWidget.Selected == AgeGroup45Plus.Text {
		isFor45Plus = true
	}
	return &SearchRequest{
		Pincode:             f.pincodeWidget.Text,
		Date:                f.dateWidget.Text,
		IsSecondDose:        f.isSecondDoseWidget.Checked,
		IsFor45Plus:         isFor45Plus,
		OnlyShowIfAvailable: f.showOnlyAvailableWidget.Checked,
	}
}
