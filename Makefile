nothing:

.PHONY: tao
tao: ./bin/tao
	@cd src && go build -o ../bin/tao

.PHONY: test
test: tao
	@cd tests && ../tools/run_test.sh

.PHONY: main
main: tao
	@./bin/tao main.tao

.PHONY: web
web:
	@cd web/src && go build -o ../../bin/web
