.DEFAULT_GOAL := run
COMMIT_MESSAGE=" - commit message"

.EXPORT_ALL_VARIABLES:

DUMPUTILS_CONFIG_PATH = test.config.yaml

run:
	go run main.go

run-docker:
	docker build -f dev.Dockerfile -t pgtools:local .
	docker run --rm --name dev-pgtools -h pgtools-local -it -v ./container:/opt/container pgtools:local /bin/bash

clean-commit:
	go clean
	gofmt -s -w .
	git add .
	git commit -m "$(COMMIT_MESSAGE)"
	git push