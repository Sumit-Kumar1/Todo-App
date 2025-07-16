# Change these variables as necessary.
BINARY_NAME := todoapp 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
    @echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## setup : to install required go tooling and air
.PHONY: setup
setup:
	go install github.com/air-verse/air@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest

## mocks: to generate mock interfaces
.PHONY: mocks
mocks:
	go generate ./...

## lint: check for lint errors
.PHONY: lint
lint:
	golangci-lint run ./... --timeout=5m


#tests: run unit tests with gotestsum
.PHONY: tests
tests:
	gotestsum --format testname -- -count=1 -p 1 -coverprofile=cover.out ./...
	go tool cover -html=cover.out


## css/watch: constantly generate css and watch for new changes
.PHONY: css/watch
css/watch:
	npx tailwindcss -i ./public/app.css -o ./public/style.css --watch

## css/output: generate output css from used classes in views/*.html
.PHONY: css/output
css/output:
	npx tailwindcss -i ./public/app.css -o ./public/style.css


## build: build the application
.PHONY: build
build: css/output
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=/tmp/bin/${BINARY_NAME} .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=./Build/main .

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/air-verse/air@v1.61.1 \
        --build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
        --build.exclude_dir "" \
        --build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
        --misc.clean_on_exit "true"

## docker/image : build the docker image
.PHONY: docker/image
docker/image: build
	docker buildx build -t todoapp . --no-cache --progress=plain

## run/container : run the docker container from the image build
.PHONY: run/container
run/container : docker/image
	docker run --name todoapp -p 9001:9001 -d todoapp:latest

## deploy/local : deploys the container image on local instance using kubectl
.PHONY: deploy/local
deploy/local: docker/image
	kubectl apply -f deployment.yaml
