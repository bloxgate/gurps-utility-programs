package astronomy

import (
	"fmt"
	"math"
)

// StarType is the spectral type of a star
type StarType int

const (
	// Yellow star color
	Yellow StarType = iota
	// Orange star color
	Orange
	// White star color
	White
	// Violet star color
	Violet
	// Blue star color
	Blue
	// Pink star color
	Pink
	// Red star color
	Red
)

func (s StarType) String() string {
	return [...]string{"Yellow (G)", "Orange (K)", "Blue-White (F)", "T Dwarf", "Blue (A/B/O)", "L Dwarf", "Red Dwarf (M)"}[s]
}

// Star represents a star on the galaxy map
type Star struct {
	RightAscension    float64
	Declination       float64
	DistanceFromEarth float64
	Name              string
	Type              StarType
}

// Distance calculates the distance to destination in parsecs
func (origin Star) Distance(destination Star) float64 {

	fmt.Printf("Orig %v\n", origin)
	fmt.Printf("Dest %v\n", destination)

	// Convert equatorial coordinates to cartesian
	// Note that this does differ from the normal conversion of spherical to cartesian for some reason
	xOr := origin.DistanceFromEarth * math.Cos(origin.RightAscension) * math.Cos(origin.Declination)
	yOr := origin.DistanceFromEarth * math.Sin(origin.RightAscension) * math.Cos(origin.Declination)
	zOr := origin.DistanceFromEarth * math.Sin(origin.Declination)

	xDe := destination.DistanceFromEarth * math.Cos(destination.RightAscension) * math.Cos(destination.Declination)
	yDe := destination.DistanceFromEarth * math.Sin(destination.RightAscension) * math.Cos(destination.Declination)
	zDe := destination.DistanceFromEarth * math.Sin(destination.Declination)

	distance := (xOr-xDe)*(xOr-xDe) + (yOr-yDe)*(yOr-yDe) + (zOr-zDe)*(zOr-zDe)
	distance = math.Sqrt(distance)

	//Convert to parsecs
	distance *= 0.306601
	return distance
}
