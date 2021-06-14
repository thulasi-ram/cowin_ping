package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/tabwriter"
	"time"
)

const (
	PrintDataFormat = "%s\t%s\t%s\t%v\t%s\n"
)

type VaccineCenters []*VaccineCenter

type VaccineCenter struct {
	Name              string `json:"name"`
	Fee               string `json:"fee"`
	Vaccine           string `json:"vaccine"`
	MinAgeLimit       int    `json:"min_age_limit"`
	AvailableCapacity int    `json:"available_capacity"`
	AvailableDose1    int    `json:"available_capacity_dose_1"`
	AvailableDose2    int    `json:"available_capacity_dose_2"`
}

func (v *VaccineCenter) Comments() string {
	return fmt.Sprintf("dose1: %v, dose2: %v", v.AvailableDose1, v.AvailableDose2)
}

type AgeGroup struct {
	Min  int
	Max  int
	Text string
}

var (
	//AgeGroup18Minus = AgeGroup{0, 17}
	AgeGroup18Plus = AgeGroup{18, 45, "18+ Group"}
	AgeGroup45Plus = AgeGroup{45, 200, "45+ Group"}
)

type SearchRequest struct {
	Pincode             string
	Date                string
	IsSecondDose        bool
	IsFor45Plus         bool
	OnlyShowIfAvailable bool
}

func (r *SearchRequest) AgeGroup() AgeGroup {
	if r.IsFor45Plus {
		return AgeGroup45Plus
	}
	return AgeGroup18Plus
}

func makeRequest(pinCode string, date string) VaccineCenters {
	resp, err := http.Get(fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByPin?pincode=%s&date=%s", pinCode, date))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var response struct {
		Sessions VaccineCenters `json:"sessions"`
	}
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	return response.Sessions
}

func printData(res VaccineCenters, writer io.Writer) {
	w := tabwriter.NewWriter(writer, 8, 8, 3, ' ', 0)
	_, _ = fmt.Fprintf(w, PrintDataFormat, "Name", "Fee", "Vaccine", "Available", "Comments")
	for _, c := range res {
		_, _ = fmt.Fprintf(w, PrintDataFormat, c.Name, c.Fee, c.Vaccine, c.AvailableCapacity, c.Comments())
	}
	_ = w.Flush()
}

func filterCenters(vc VaccineCenters, r *SearchRequest) VaccineCenters {
	filtered := VaccineCenters{}
	for _, c := range vc {
		if c.MinAgeLimit >= r.AgeGroup().Min && c.MinAgeLimit < r.AgeGroup().Max {
			if r.OnlyShowIfAvailable {
				if c.AvailableCapacity > 0 {
					filtered = append(filtered, c)
				}
			} else {
				filtered = append(filtered, c)
			}
		}
	}
	return filtered
}

func PeriodicPushData(ctx context.Context, r *SearchRequest, centers chan VaccineCenters) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			resp := makeRequest(r.Pincode, r.Date)
			validCenters := filterCenters(resp, r)
			centers <- validCenters
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func GetFormattedDataAndMakeSound(ctx context.Context, centers VaccineCenters) (message bytes.Buffer) {
	var b bytes.Buffer
	printData(centers, &b)
	if len(centers) > 0 {
		go func(ctx context.Context) {
			makeSound(ctx)
		}(ctx)
	}
	return b
}
