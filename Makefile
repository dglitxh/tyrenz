all: clean build linkbin

clean:
	echo "cleaning......"
	go mod tidy

build: 
	echo "building app"
	go build .

linkbin: 
	echo "Creating symlink..."
	sudo ln -s  $(realpath ./tyrenz) "/usr/bin/tyrenz"