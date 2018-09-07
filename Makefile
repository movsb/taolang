nothing:

.PHONY: test
test: ./bin/tao
	@cd tests && ../tools/run_test.sh

.PHONY: tao
tao: ./bin/tao
	@cd src && go build -o ../bin/tao

.PHONY: main
main: ./bin/tao
	@./bin/tao main.tao

.PHONY: web
web:
	@cd web/src && go build -o ../../bin/web
