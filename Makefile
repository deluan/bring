
.PHONY: run
run:
	cd sample; go run . vnc 10.0.0.11 5901

.PHONY: rdp
rdp:
	cd sample; go run . rdp `ipconfig getifaddr en0` 3389

.PHONY: qemu
qemu:
	cd sample; go run . vnc `ipconfig getifaddr en0` 5900

.PHONY: watch
watch:
	ginkgo watch -notify ./...

.PHONY: test
test:
	ginkgo -coverprofile=reports/coverage.out ./... -v

.PHONY: bench
bench:
	go test -bench=. -run=XXX ./...

.PHONY: coverage
coverage:
	mkdir -p reports
	go test -coverprofile=reports/coverage.out -v -coverpkg ./... ./...
	go tool cover -func=reports/coverage.out
	go tool cover -html=reports/coverage.out -o reports/index.html
	open reports/index.html

.PHONY: test
doc:
	@echo "Doc server address: http://localhost:6060/pkg"
	godoc -http=:6060

.PHONY: release
release:
	@if [[ ! "${V}" =~ ^[0-9]+\.[0-9]+\.[0-9]+.*$$ ]]; then echo "Usage: make release V=X.X.X"; exit 1; fi
	go mod tidy
	make test
	@if [ -n "`git status -s`" ]; then echo "\n\nThere are pending changes. Please commit first"; exit 1; fi
	git tag v${V}
	git push origin v${V}
	git push origin master
