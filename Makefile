
run:
	cd sample; go run . vnc 10.0.0.11 5901

rdp:
	cd sample; go run . rdp `ipconfig getifaddr en0` 3389

qemu:
	cd sample; go run . vnc `ipconfig getifaddr en0` 5900

watch:
	goconvey -cover -excludedDirs testdata .

.PHONY: test
test:
	go test -cover -v ./...

bench:
	go test -bench=. -run=XXX ./...

coverage:
	mkdir -p reports
	go test -coverprofile=reports/coverage.out
	go tool cover -func=reports/coverage.out
	go tool cover -html=reports/coverage.out -o reports/index.html
	open reports/index.html

doc:
	@echo "Doc server address: http://localhost:6060"
	godoc -http=":6060" -goroot=$$GOPATH

release:
	@if [[ ! "${V}" =~ ^[0-9]+\.[0-9]+\.[0-9]+.*$$ ]]; then echo "Usage: make release V=X.X.X"; exit 1; fi
	go mod tidy
	make test
	git add .
	git ci -m "Release v${V}"
	git tag v${V}
	git push origin v${V}
	git push master
