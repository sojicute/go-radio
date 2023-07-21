build:
	go env -w GOARCH=wasm GOOS=js
	go build -o web/app.wasm ./app

	go env -w GOARCH=amd64 GOOS=windows
	go build -o hello ./server

run: build
	./hello