nothing:

test: ./tao
	@./run_test.sh

build: ./tao

./tao:
	cd src && go build -o ../tao

main: build
	./tao main.tao
