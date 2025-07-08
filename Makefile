APP_NAME := crawly
.PHONY: run build clean

run:
	go run .

build:
	go build -o $(APP_NAME) .

clean:
	rm -f $(APP_NAME)
