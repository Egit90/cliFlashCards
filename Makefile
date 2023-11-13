
Bin_name=flashCards

build:
	go build -o $(Bin_name) -v

run: build
	 ./$(Bin_name) 

help: build
	./$(Bin_name) --help