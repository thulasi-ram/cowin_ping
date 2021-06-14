package main

import (
	"errors"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func PincodeWidget() *widget.Entry {
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

func DateWidget() *widget.Entry {

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
