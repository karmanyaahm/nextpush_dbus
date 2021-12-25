build-local:
	go build -o nextpush_dbus
test: build-local
	go test ./...

build-docker:
	mkdir bin
	chmod ugo+rwx bin
	docker run --rm -v `pwd`:/app ghcr.io/karmanyaahm/mega_go_arch_xcompiler:v0.2.1 build nextpush

