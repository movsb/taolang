nothing:

.PHONY: test
test: ./tao
	@./run_test.sh

.PHONY: tao
tao: 
	@cd src && go build -o ../bin/tao

.PHONY:main
main: tao
	@./bin/tao main.tao

.PHONY: web
web:
	@cd web/src && go build -o ../../bin/web
