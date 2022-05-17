GO_MOD_NAME := $(shell go list -m)
GO_MOD_DOMAIN := $(shell echo $(GO_MOD_NAME) | awk -F '/' '{print $$1}')
GO_MOD_BASE_NAME := $(shell echo $(GO_MOD_NAME) | awk -F '/' '{print $$NF}')
SERVICE_NAME := $(GO_MOD_BASE_NAME)

NOW := $(shell date +'%Y%m%d%H%M%S')
TAG := $(shell git describe --always --tags --abbrev=0 | tr -d "[\r\n]")
COMMIT := $(shell git rev-parse --short HEAD| tr -d "[ \r\n\']")
VERSION_PKG := github.com/prometheus/common/version

IMPORTANT_GO_ENV_VARS := "GOPATH|GO111MODULE|GOARCH|GOCACHE|GOMODCACHE|GONOPROXY|GONOSUMDB|GOPRIVATE|GOPROXY|GOSUMDB|GOMOD|CGO"

LD_FLAGS_BASE := -X main.ServiceName=$(SERVICE_NAME) \
				-X $(VERSION_PKG).Version=$(TAG) \
				-X $(VERSION_PKG).Revision=$(COMMIT) \
				-X $(VERSION_PKG).Branch=$(shell git rev-parse --abbrev-ref HEAD) \
				-X $(VERSION_PKG).BuildUser=$(shell whoami) \
				-X $(VERSION_PKG).BuildDate=$(shell date +%Y%m%d-%H%M%S)
LD_FLAGS := -s -w $(LD_FLAGS_BASE)


.PHONY: build/debug
build/debug:
	@echo "\n building debug binary $(SERVICE_NAME)"
	@go env | grep -E $(IMPORTANT_GO_ENV_VARS)
	@go version
	go build -gcflags "all=-N -l" -ldflags="$(LD_FLAGS_BASE)" -o $(SERVICE_NAME)

.PHONY: build/release
build/release: export CGO_ENABLED=0
build/release:
	@echo "\n building release binary $(SERVICE_NAME)"
	@go env | grep -E $(IMPORTANT_GO_ENV_VARS)
	/usr/bin/time -f '%Us user %Ss system %P cpu %es total' go build -trimpath -ldflags="$(LD_FLAGS)" -o $(SERVICE_NAME)

.PHONY: clean
clean:
	@rm -f $(SERVICE_NAME)
	@rm -f $(SERVICE_NAME).tar.gz

.PHONY: fmt
fmt:
	@command -v gofumpt || go install mvdan.cc/gofumpt@latest
	# gofumpt: -s is deprecated as it is always enabled
	gofumpt -extra -w -d .
	@command -v gci || go install github.com/daixiang0/gci@latest
	gci write -s Std -s Def -s 'Prefix($(GO_MOD_DOMAIN))' -s 'Prefix($(GO_MOD_NAME))' .

.PHONY: lint
lint:
	golangci-lint version
	golangci-lint run -v --color always --out-format colored-line-number

.PHONY: cfg
cfg: build/debug
	./$(SERVICE_NAME) --dump > ./config/config.dist.toml

.PHONY: changelog
changelog:
	#git-chglog > CHANGELOG.md
	git cliff | tee CHANGELOG.md | bat -l markdown -P

.PHONY: gitleaks
gitleaks:
	gitleaks detect -v

.PHONY: makefile/test
makefile/test:
	@echo "GO_MOD_NAME:      $(GO_MOD_NAME)"
	@echo "GO_MOD_DOMAIN:    $(GO_MOD_DOMAIN)"
	@echo "GO_MOD_BASE_NAME: $(GO_MOD_BASE_NAME)"
	@echo "SERVICE_NAME:     $(SERVICE_NAME)"
