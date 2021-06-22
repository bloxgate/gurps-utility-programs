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

	var galaxy []*astronomy.Star
	var userStar *astronomy.Star = &astronomy.Star{RightAscension: 0, Declination: 0, DistanceFromEarth: 0, Name: "User Input", Type: 0}
	var starNames []string

	var origin, dest *astronomy.Star
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

		//Replaces all spaces inside the quoted name with _, and then removes the quotes
		//This lets us scan in the line easily
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

		galaxy = append(galaxy, &s)
	}
	starNames = append(starNames, "")
	copy(starNames[2:], starNames[1:])
	starNames[1] = "User Star"

	a := app.New()
	mainWin := a.NewWindow("Space Travel Time Calculator")

	// Begin Calculator Tab
	distLabel := widget.NewLabel("Distance: ")
	distValue := widget.NewLabel("")
	travelTimeLabel := widget.NewLabel("Travel Time:")
	travelTimeValue := widget.NewLabel("")

	originDrop := widget.NewSelect(starNames, func(changed string) {
		for _, s := range galaxy {
			if changed == "User Star" {
				origin = userStar
			} else if s.Name == strings.ReplaceAll(changed, " ", "_") {
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
			if changed == "User Star" {
				dest = userStar
			} else if s.Name == strings.ReplaceAll(changed, " ", "_") {
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
	calculatorContainer := container.NewVBox(
		form,
		container.NewHBox(distLabel, distValue),
		container.NewHBox(travelTimeLabel, travelTimeValue),
	)
	calculatorTab := container.NewTabItem("Calculator", calculatorContainer)
	//End Calculator Tab

	//Begin Add Star Tab
	raHours := widget.NewEntry()
	raMins := widget.NewEntry()
	raSecs := widget.NewEntry()
	rightAscensionEntry := container.NewHBox(raHours, widget.NewLabel("h"), raMins, widget.NewLabel("m"), raSecs, widget.NewLabel("s"))

	decSign := widget.NewCheck("Negative?", nil)
	decDegs := widget.NewEntry()
	decMins := widget.NewEntry()
	decSecs := widget.NewEntry()
	declinationEntry := container.NewHBox(decSign, decDegs, widget.NewLabel("\u00B0"), decMins, widget.NewLabel("\u2032"), decSecs, widget.NewLabel("\u2033"))

	addDistEntry := widget.NewEntry()

	var addStarType astronomy.StarType
	addStarTypeSelect := widget.NewSelect(astronomy.StarTypeStrings, nil)

	alpha := widget.NewLabel("")
	alphaLabel := widget.NewLabel("\u03B1:")
	delta := widget.NewLabel("")
	deltaLabel := widget.NewLabel("\u03B4:")
	addDist := widget.NewLabel("")
	addDistLabel := widget.NewLabel("d:")
	addStarTypeLabel := widget.NewLabel("Palette: ")
	addStarPalette := widget.NewLabel("")

	setAsUserInput := widget.NewButton("Set As User Input", func() {
		scanText := fmt.Sprintf("%s %s %s %d", alpha.Text, delta.Text, addDist.Text, addStarType)

		n, err := fmt.Sscan(scanText, &userStar.RightAscension, &userStar.Declination, &userStar.DistanceFromEarth, &userStar.Type)
		if n < 4 {
			fmt.Printf("Error reading user input: %v\n", err)
		}
		userStar.RightAscension *= (math.Pi / 180)
		userStar.Declination *= (math.Pi / 180)
	})
	setAsUserInput.Disable()

	raFormItem := widget.NewFormItem("Right Ascension", rightAscensionEntry)
	decFormItem := widget.NewFormItem("Declination", declinationEntry)
	addDistFormItem := widget.NewFormItem("Distance", addDistEntry)
	addStarTypeFormItem := widget.NewFormItem("Star Type", addStarTypeSelect)
	addStarForm := widget.NewForm(raFormItem, decFormItem, addDistFormItem, addStarTypeFormItem)
	addStarForm.SubmitText = "Convert Measurements"
	addStarForm.OnSubmit = func() {
		addStarType = astronomy.StarType(addStarTypeSelect.SelectedIndex())
		addStarPalette.SetText(fmt.Sprintf("%d", addStarType))

		hadError := false

		d, err := strconv.ParseFloat(addDistEntry.Text, 64)
		if err != nil {
			addDist.SetText("Invalid Distance")
			addDist.Refresh()
			hadError = true
		} else {
			addDist.SetText(fmt.Sprint(d))
		}

		var h, m int64
		var s float64
		var raMeasure *astronomy.HMS
		h, err = strconv.ParseInt(raHours.Text, 10, 64)
		if err != nil {
			alpha.SetText("Invalid Right Ascension")
			alpha.Refresh()
			hadError = true
		} else {
			m, err = strconv.ParseInt(raMins.Text, 10, 64)
			if err != nil {
				alpha.SetText("Invalid Right Ascension")
				alpha.Refresh()
				hadError = true
			} else {
				s, err = strconv.ParseFloat(raSecs.Text, 64)
				if err != nil {
					alpha.SetText("Invalid Right Ascension")
					alpha.Refresh()
					hadError = true
				} else {
					raMeasure = &astronomy.HMS{Hours: h, Minutes: m, Seconds: s}
					alpha.SetText(raMeasure.ToDegrees().String())
					alpha.Refresh()
				}
			}
		}

		h, err = strconv.ParseInt(decDegs.Text, 10, 64)
		if err != nil {
			delta.SetText("Invalid Declination")
			delta.Refresh()
			return
		} else {
			m, err = strconv.ParseInt(decMins.Text, 10, 64)
			if err != nil {
				delta.SetText("Invalid Declination")
				delta.Refresh()
				return
			} else {
				s, err = strconv.ParseFloat(decSecs.Text, 64)
				if err != nil {
					delta.SetText("Invalid Declination")
					delta.Refresh()
					return
				} else {
					decMeasure := &astronomy.DMS{Degrees: h, Minutes: m, Seconds: s, Sign: decSign.Checked}
					delta.SetText(decMeasure.ToDegrees().String())
					delta.Refresh()
				}
			}
		}

		if !hadError {
			setAsUserInput.Enable()
		}
	}

	addStarContainer := container.NewVBox(
		addStarForm,
		container.NewHBox(alphaLabel, alpha),
		container.NewHBox(deltaLabel, delta),
		container.NewHBox(addDistLabel, addDist),
		container.NewHBox(addStarTypeLabel, addStarPalette),
		setAsUserInput,
	)
	addStarTab := container.NewTabItem("Converter", addStarContainer)
	//End Add Star Tab

	tabs := container.NewAppTabs(calculatorTab, addStarTab)

	mainWin.SetContent(tabs)

	ftlVelEntry.Validate()
	form.Refresh()

	mainWin.ShowAndRun()
}
