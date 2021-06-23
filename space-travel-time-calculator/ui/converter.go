package ui

import (
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"travelcalc/astronomy"
)

func CreateConverterTab(userStar *astronomy.Star) *container.TabItem {
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

	return addStarTab
}
