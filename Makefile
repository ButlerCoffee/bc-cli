compile:
	go build -o bc-cli .

install:
	pip install pip -U
	pip install -r requirements-dev.txt
	pre-commit install
	pre-commit autoupdate
	pre-commit install-hooks
	pre-commit install --hook-type commit-msg

upgrade:
	go get -u ./...
	go mod tidy
	pip-upgrade

.PHONY:
.ONESHELL:
