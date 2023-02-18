all: explosioncalc

explosioncalc: GOOS=win
explosioncalc: GOARCH=64
explosioncalc: main.go
	go build -o explosioncalc.exe -ldflags="-w -s -H=windowsgui" main.go

clean:
	rm explosioncalc.exe
