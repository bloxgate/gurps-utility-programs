package astronomy

import "math/big"

//DMS represents a measurement in Degrees Minutes Seconds
type DMS struct {
	Degrees int64
	Minutes int64
	Seconds float64
	Sign    bool //Sign - if sign is true this number is negative
}

var (
	minutesToDegrees *big.Float = big.NewFloat(60)   //60 minutes per degree
	secondsToDegrees *big.Float = big.NewFloat(3600) //60*60 seconds per degree
)

//ToDegrees converts a DMS measure to degrees
func (measurement DMS) ToDegrees() *big.Float {
	deg := big.NewFloat(float64(measurement.Degrees))
	mins := big.NewFloat(float64(measurement.Minutes))
	secs := big.NewFloat(measurement.Seconds)

	mins.Quo(mins, minutesToDegrees)
	secs.Quo(secs, secondsToDegrees)

	deg.Add(deg, mins)
	deg.Add(deg, secs)

	if measurement.Sign {
		deg.Neg(deg)
	}

	return deg
}

//HMS represents a measurement in Hours Minutes Seconds
type HMS struct {
	Hours   int64
	Minutes int64
	Seconds float64
}

//ToDegrees converts an HMS measurement to degrees
func (measurement HMS) ToDegrees() *big.Float {
	asDMS := &DMS{Degrees: measurement.Hours * 15, Minutes: measurement.Minutes * 15, Seconds: measurement.Seconds * 15, Sign: false}

	return asDMS.ToDegrees()
}
