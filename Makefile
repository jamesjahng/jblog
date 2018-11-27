run:
	./bin/terrace

build:
	go build -x -o ./bin/terrace ./src
	sudo setcap CAP_NET_BIND_SERVICE=+eip ./bin/terrace

all:
	@echo "Doing all"

deploy:
	@echo "Pushing to production"
	@git push git@example.com:~/testapp master

update:
	@echo "Makefile: Doing UPDATE stuff like grunt, gulp, rake,..."
	@whoami
	@pwd
