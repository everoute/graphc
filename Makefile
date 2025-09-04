.PHONY: image-generate generate docker-generate test docker-test publish

image-generate:
	docker build -f build/image/generate/Dockerfile -t localhost/generate ./build/image/generate/

generate:
	find . -name "*.go" -exec gci write --Section Standard --Section Default --Section "Prefix(github.com/everoute/graphc)" {} +

docker-generate: image-generate
	$(eval WORKDIR := /go/src/github.com/everoute/graphc)
	docker run --rm -iu 0:0 -w $(WORKDIR) -v $(CURDIR):$(WORKDIR) localhost/generate make generate

test:
	go test ./... --race --coverprofile coverage.out

docker-test:
	$(eval WORKDIR := /go/src/github.com/everoute/graphc)
	docker run --rm -iu 0:0 -w $(WORKDIR) -v $(CURDIR):$(WORKDIR) registry.smtx.io/sdn-base/golang:1.20 make test

debug-test:
	$(eval WORKDIR := /go/src/github.com/everoute/graphc)
	docker run --rm -iu 0:0 -w $(WORKDIR) -v $(CURDIR):$(WORKDIR) registry.smtx.io/sdn-base/golang:1.20 bash

publish:
	go build -o graphc_codegen tools/codegen/main.go
