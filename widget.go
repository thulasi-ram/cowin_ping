package main

import (
	"errors"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type FormItems []*widget.FormItem

func PincodeWidget() *widget.Entry {
	pinCodeValidator := func(p string) error {
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
	pinCodeInput := widget.NewEntry()
	pinCodeInput.Validator = pinCodeValidator
	return pinCodeInput
}

func DateWidget() *widget.Entry {
	dateInput := widget.NewEntry()
	dateInput.SetPlaceHolder("DD-MM-YYYY")
	return dateInput
}
