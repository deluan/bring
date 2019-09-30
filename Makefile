
run:
	go run app/main.go rdp `ipconfig getifaddr en0` 3389

offline:
	go run app/main.go vnc 10.0.0.11 5901

qemu:
	go run app/main.go vnc `ipconfig getifaddr en0` 5900

watch:
	goconvey -cover -excludedDirs testdata .

.PHONY: test
test:
	go test -cover -v .

bench:
	go test -bench=. -run=XXX ./...

doc:
	@echo "Doc server address: http://localhost:6060"
	godoc -http=":6060" -goroot=$$GOPATH
