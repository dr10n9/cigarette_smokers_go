build:
	go build -o app.out

run: build
	./app.out

clean:
	rm -rf *.out