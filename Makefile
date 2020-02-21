.PHONY: packge clean build
package:
	echo "packaging"
	fyne package
clean:
	rm -rf amazon amazon.app

build:
	go build -o bin/amazon main.go