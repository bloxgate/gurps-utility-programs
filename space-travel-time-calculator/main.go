package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"

	"./astronomy"
)

func main() {
	os.Setenv("FYNE_THEME", "dark")

	var galaxy []astronomy.Star
	var starNames []string

	var origin, dest astronomy.Star
	var ftlVelocity float64

	starDB, err := os.Open("stars.dat")
	if err != nil {
		log.Fatal("Unable to read stars database!")
	}
	defer starDB.Close()
	starReader := bufio.NewScanner(starDB)
	starReader.Scan()

	for starReader.Scan() {
		line := starReader.Text()

		if strings.Index(line, "#") == 0 {
			continue
		}

		ind := strings.Index(line, "\"")
		trim1 := line[ind:]
		ind = strings.LastIndex(trim1, "\"") + 1
		trim2 := trim1[:ind]
		starNames = append(starNames, strings.ReplaceAll(trim2, "\"", ""))
		trim2 = strings.ReplaceAll(trim2, " ", "_")
		re, err := regexp.Compile("\".*\"")
		if err != nil {
			log.Fatal("Error in space removal regex", err)
		}
		line = re.ReplaceAllString(line, trim2)
		line = strings.Replace(line, "\"", "", 2)
		s := astronomy.Star{}

		_, err = fmt.Sscanln(line, &s.RightAscension, &s.Declination, &s.DistanceFromEarth, &s.Name, &s.Type)
		if err != nil {
			log.Fatalf("Error parsing star DB: %v", err)
		}

		//Convert angles to radians
		s.RightAscension = s.RightAscension * (math.Pi / 180.0)
		s.Declination = s.Declination * (math.Pi / 180.0)

		galaxy = append(galaxy, s)
	}

	a := app.New()
	mainWin := a.NewWindow("Space Travel Time Calculator")

	distLabel := widget.NewLabel("Distance: ")
	distValue := widget.NewLabel("")
	travelTimeLabel := widget.NewLabel("Travel Time:")
	travelTimeValue := widget.NewLabel("")

	originDrop := widget.NewSelect(starNames, func(changed string) {
		for _, s := range galaxy {
			if s.Name == strings.ReplaceAll(changed, " ", "_") {
				origin = s
				log.Printf("Origin set to %s\n", s.Name)
				break
			}
		}
		travelTimeValue.Text = ""
		travelTimeValue.Refresh()
		distValue.Text = ""
		distValue.Refresh()
	})
	destDrop := widget.NewSelect(starNames, func(changed string) {
		for _, s := range galaxy {
			if s.Name == strings.ReplaceAll(changed, " ", "_") {
				dest = s
				log.Printf("Destination set to %s\n", s.Name)
				break
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
		ftlVelocity, err = strconv.ParseFloat(ftlVelEntry.Text, 64)
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

	mainWin.SetContent(container.NewVBox(
		form,
		container.NewHBox(distLabel, distValue),
		container.NewHBox(travelTimeLabel, travelTimeValue),
	))

	ftlVelEntry.Validate()
	form.Refresh()

	mainWin.ShowAndRun()
}
