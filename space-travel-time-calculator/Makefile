all: travelcalc

travelcalc: GOOS=win
travelcalc: GOARCH=64
travelcalc: main.go
	go build -o travelcalc.exe -ldflags="-w -s -H=windowsgui" main.go

clean:
	rm travelcalc.exe