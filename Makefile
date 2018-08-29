nothing:

build:
	cd src && go build -o ../tao

main: build
	./tao main.tao
