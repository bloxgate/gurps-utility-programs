package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"travelcalc/astronomy"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

var (
	origin, dest *astronomy.Star
	ftlVelocity  float64
)

func CreateCalculatorTab(starNames []string, galaxy []*astronomy.Star, userStar *astronomy.Star) *container.TabItem {
	distLabel := widget.NewLabel("Distance: ")
	distValue := widget.NewLabel("")
	travelTimeLabel := widget.NewLabel("Travel Time:")
	travelTimeValue := widget.NewLabel("")

	originDrop := widget.NewSelect(starNames, func(changed string) {
		if changed == "User Star" {
			origin = userStar
			log.Println("Origin set to user input")
		} else {
			for _, s := range galaxy {
				if s.Name == strings.ReplaceAll(changed, " ", "_") {
					origin = s
					log.Printf("Origin set to %s\n", s.Name)
					break
				}
			}
		}
		travelTimeValue.Text = ""
		travelTimeValue.Refresh()
		distValue.Text = ""
		distValue.Refresh()
	})
	destDrop := widget.NewSelect(starNames, func(changed string) {
		if changed == "User Star" {
			dest = userStar
			log.Println("Destination set to user input")
		} else {
			for _, s := range galaxy {
				if s.Name == strings.ReplaceAll(changed, " ", "_") {
					dest = s
					log.Printf("Destination set to %s\n", s.Name)
					break
				}
			}
		}
		travelTimeValue.Text = ""
		travelTimeValue.Refresh()
		distValue.Text = ""
		distValue.Refresh()
	})
	ftlVelEntry := widget.NewEntry()
	ftlVelEntry.Validator = validation.NewRegexp("^\\d+\\.?\\d*$", "Not a valid number!")

	originFormItem := widget.NewFormItem("Origin", originDrop)
	destFormItem := widget.NewFormItem("Destination", destDrop)
	velFormItem := widget.NewFormItem("FTL Velocity (pc/day)", ftlVelEntry)

	form := widget.NewForm(originFormItem, destFormItem, velFormItem)
	form.OnSubmit = func() {
		ftlVelocity, err := strconv.ParseFloat(ftlVelEntry.Text, 64)
		if err != nil {
			travelTimeValue.Text = "Invalid Speed"
			travelTimeValue.Refresh()
			return
		}

		dist := origin.Distance(dest)
		distValue.Text = fmt.Sprintf("%.4f pc", dist)
		distValue.Refresh()

		durString := fmt.Sprintf("%fh", (dist/ftlVelocity)*24)
		duration, err := time.ParseDuration(durString)
		if err != nil {
			travelTimeValue.Text = "Bad Time"
			travelTimeValue.Refresh()
			return
		}

		travelTimeValue.Text = duration.String()
		travelTimeValue.Refresh()
	}
	form.SubmitText = "Calculate"
	calculatorContainer := container.NewVBox(
		form,
		container.NewHBox(distLabel, distValue),
		container.NewHBox(travelTimeLabel, travelTimeValue),
	)
	calculatorTab := container.NewTabItem("Calculator", calculatorContainer)

	ftlVelEntry.Validate()
	form.Refresh()

	return calculatorTab
}
