package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

var randSource = rand.NewSource(time.Now().UnixNano())
var randGen = rand.New(randSource)

func mulSlice(slice []int) int {
	product := 1

	for _, m := range slice {
		product *= m
	}

	return product
}

func findExpRadius(numDice ...int) int {
	return 2 * mulSlice(numDice)
}

func explosionDamage(distance int, damage int) int {
	return damage / (3 * distance)
}

func findFragRadius(numDice int) int {
	return 5 * numDice
}

func fragDamage(numDice int) int {
	sum := 0
	d := 0
	for d < numDice {
		sum += rand.Intn(6) + 1
		d++
	}

	return sum
}

func main() {
	os.Setenv("FYNE_THEME", "dark")

	a := app.New()
	mainWin := a.NewWindow("Explosion Damage Calculator")

	mainDiceEntry := widget.NewEntry()
	//mainDiceEntry.Validator = validation.NewRegexp("\\d+", "Enter number")
	mainDiceFormItem := widget.NewFormItem("Damage Dice", mainDiceEntry)
	mainDiceFormItem.HintText = "6dx2 [3d] = 6"

	diceMulEntry := widget.NewEntry()
	//diceMulEntry.Validator = validation.NewRegexp("\\d+", "Enter a number")
	diceMulFormItem := widget.NewFormItem("Dice Multiplier", diceMulEntry)
	diceMulFormItem.HintText = "6dx2 [3d] = 2"

	fragDiceEntry := widget.NewEntry()
	//fragDiceEntry.Validator = validation.NewRegexp("\\d+", "Enter a number")
	fragDiceFormItem := widget.NewFormItem("Fragmentation Dice", fragDiceEntry)
	fragDiceFormItem.HintText = "6dx2 [3d] = 3"

	distanceEntry := widget.NewEntry()
	//distanceEntry.Validator = validation.NewRegexp("\\d+", "Enter a number")
	distanceFormItem := widget.NewFormItem("Distance", distanceEntry)
	distanceFormItem.HintText = "Distance from explosion in meters"

	damageEntry := widget.NewEntry()
	damageEntry.Validator = validation.NewRegexp("\\d+", "Enter a number")
	damageFormItem := widget.NewFormItem("Total Damage", damageEntry)

	expDistance := widget.NewLabel("Explosion Radius:")
	expDistanceNum := widget.NewLabel("")

	expDamageLabel := widget.NewLabel("Explosion Damage:")
	expDamageNum := widget.NewLabel("")

	fragDistance := widget.NewLabel("Fragmentation Distance:")
	fragDistanceNum := widget.NewLabel("")

	fragDamageLabel := widget.NewLabel("Fragmentation Damage:")
	fragDamageNum := widget.NewLabel("")

	form := widget.NewForm(mainDiceFormItem, diceMulFormItem, fragDiceFormItem, damageFormItem, distanceFormItem)
	form.OnSubmit = func() {
		dmgDice, _ := strconv.Atoi(mainDiceEntry.Text)
		diceMul, _ := strconv.Atoi(diceMulEntry.Text)
		fragDice, _ := strconv.Atoi(fragDiceEntry.Text)
		distance, _ := strconv.Atoi(distanceEntry.Text)
		damage, _ := strconv.Atoi(damageEntry.Text)
		eRadius := findExpRadius(dmgDice, diceMul)
		fRadius := findFragRadius(fragDice)

		expDistanceNum.Text = fmt.Sprintf("%d", eRadius)
		fragDistanceNum.Text = fmt.Sprintf("%d", fRadius)

		expDamageNum.Text = fmt.Sprintf("%d", explosionDamage(distance, damage))
		fragDamageNum.Text = fmt.Sprintf("%d", fragDamage(fragDice))

		expDistanceNum.Refresh()
		fragDistanceNum.Refresh()
		expDamageNum.Refresh()
		fragDamageNum.Refresh()
	}
	form.SubmitText = "Calculate"

	mainWin.SetContent(container.NewVBox(
		form,
		container.NewAdaptiveGrid(2, expDistance, expDistanceNum, expDamageLabel, expDamageNum, fragDistance, fragDistanceNum, fragDamageLabel, fragDamageNum),
	))

	damageEntry.Validate()
	form.Refresh()

	mainWin.ShowAndRun()

}
