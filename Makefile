.EXPORT_ALL_VARIABLES:

PGTOOLS_CONFIG_PATH = local.config.yaml

run:
	go run main.go

run-docker:
	docker build -f dev.Dockerfile -t pgtools:local .
	docker run --rm --name dev-pgtools -h pgtools-local -it -v ./container:/opt/container pgtools:local /bin/bash