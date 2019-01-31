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
