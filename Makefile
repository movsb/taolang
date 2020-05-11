IMAGE = taocker/taolang:latest

nothing:

.PHONY: tao
tao:
	@cd main && go build -o ../bin/tao

.PHONY: web
web:
	@cd web/src && go build -o ../../bin/web

.PHONY: tests
tests: tao
	@cd tests && ../tools/run_test.sh

.PHONY: examples
examples: tao
	@cd web/examples && ../../tools/run_test.sh

.PHONY: main
main: tao
	@./bin/tao --main main.tao

.PHONY: repl
repl: tao
	@./bin/tao

.PHONY: all
all: tao web tests examples

.PHONY: build-image
build-image:
	cd main && GOOS=linux GOARCH=amd64 go build -o ../docker/bin/tao
	cd web/src && GOOS=linux GOARC=amd64 go build -o ../../docker/bin/web
	rsync -aPvh ./web/{examples,html} docker/web
	cd docker && docker build -t ${IMAGE} .

.PHONY: push-image
push-image:
	docker push ${IMAGE}
