package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"travelcalc/astronomy"
	"travelcalc/ui"
)

func main() {
	os.Setenv("FYNE_THEME", "dark")

	var galaxy []*astronomy.Star
	var userStar *astronomy.Star = &astronomy.Star{RightAscension: 0, Declination: 0, DistanceFromEarth: 0, Name: "User Input", Type: 0}
	var starNames []string

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
	calculatorTab := ui.CreateCalculatorTab(starNames, galaxy, userStar)
	//End Calculator Tab

	//Begin Add Star Tab
	addStarTab := ui.CreateConverterTab(userStar)
	//End Add Star Tab

	tabs := container.NewAppTabs(calculatorTab, addStarTab)

	mainWin.SetContent(tabs)

	mainWin.ShowAndRun()
}
