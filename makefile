# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd
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
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## lint: check for lint errors
.PHONY: lint
lint:
	golangci-lint run ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## css/watch: constantly generate css and watch for new changes
.PHONY: css/watch
css/watch:
	npx tailwindcss -i ./public/input.css -o ./public/style.css --watch

## css/output: generate output css from used classes in views/*.html
.PHONY: css/output
css/output:
	npx tailwindcss -i ./public/input.css -o ./public/style.css


## build: build the application
.PHONY: build
build: css/output
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=./Build/main ${MAIN_PACKAGE_PATH}

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