GOALIAS     ?= go

COVER_OUT   := coverage.out
COVER_HTML  := coverage.html
PKG_COVER   := pkg.cover
BIN         := traning-oauth-go
REPO        := bitbucket.org/hebertthome/traning-oauth-go
MOCKGEN     := $(GOBIN)/mockgen
LINT        := $(GOBIN)/golint
GOTOOLDIR   := $(shell $(GOALIAS) env GOTOOLDIR)
VET         := $(GOTOOLDIR)/vet
COVER       := $(GOTOOLDIR)/cover
BUILD       := $(shell git rev-parse --short HEAD)
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
SOURCES     := $(shell find . -type f -name '*.go')
PKGS        := $(shell $(GOALIAS) list ./...)
CONFIG      := $(BASE_DIR)/config.yaml
CONFIG_LOCAL := $(BASE_DIR)/config-local.yaml



.PHONY: build
build: check_gopath $(BIN)


.PHONY: all
all: clean cover lint vet build


.PHONY: check_gopath
check_gopath:
ifndef GOPATH
	@echo "ERROR!! GOPATH must be declared. Check http://golang.org/doc/code.html#GOPATH"
	@exit 1
endif
ifeq ($(shell $(GOALIAS) list ./... | grep -q '^_'; echo $$?), 0)
	@echo "ERROR!! This directory should be at $(GOPATH)/src/$(REPO)"
	@exit 1
endif
	@exit 0


.PHONY: check_gobin
check_gobin:
ifndef GOBIN
	@echo "ERROR!! GOBIN must be declared. Check http://golang.org/doc/code.html#GOBIN"
	@exit 1
endif
	@exit 0


.PHONY: env
env:
	@echo "Dumping Go environment vars..."
	@GOPATH=$(GOPATH) $(GOALIAS) env


.PHONY: clean
clean:
	@echo "Removing temp files..."
	@rm -fv $(BIN) $(COVER_HTML) $(COVER_OUT) $(PKG_COVER)
	@rm -fv *.cover *.out
	@find . -name '.*.swp' -exec rm -fv {} \;
	@GOPATH=$(GOPATH) $(GOALIAS) clean -v


.PHONY: test
test: check_gopath
	@for pkg in $(PKGS); do \
		GOPATH=$(GOPATH) $(GOALIAS) test -v -race $$pkg || exit 1; \
	done


.PHONY: cover
cover: check_gopath $(COVER)
	@echo Running tests with coverate report...
	@echo 'mode: set' > $(COVER_OUT)
	@touch $(PKG_COVER)
	@for pkg in $(PKGS); do \
		GOPATH=$(GOPATH) $(GOALIAS) test -v -coverprofile=$(PKG_COVER) $$pkg || exit 1; \
		grep -v 'mode: set' $(PKG_COVER) >> $(COVER_OUT); \
	done
	@echo Generating HTML report in $(COVER_HTML)...
	@GOPATH=$(GOPATH) $(GOALIAS) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@(which -s open && open $(COVER_HTML)) || (which -s gnome-open && gnome-open $(COVER_HTML)) || (exit 0)


.PHONY: run
run: $(BIN)
	@echo "Starting $(BIN) with $(CONFIG_LOCAL)..."
	@./$(BIN) -conf $(CONFIG_LOCAL)

.PHONY: lint
lint: $(LINT)
	@for src in $(SOURCES); do \
			GOPATH=$(GOPATH) golint $$src; \
	done

# $(VET) removed on vet for issue (https://github.com/golang/go/issues/11659)
.PHONY: vet
vet: check_gopath
	@for src in $(SOURCES); do \
      GOPATH=$(GOPATH) $(GOALIAS) vet $$src; \
	done


$(BIN): $(SOURCES) # env
	@echo "Building $(BIN) $(BUILD)..."
	@GOPATH=$(GOPATH) $(GOALIAS) build -o $(BIN) $(REPO)

$(MOCKGEN): check_gobin check_gopath
	@$(GOALIAS) get code.google.com/p/gomock/mockgen

$(GIN): check_gopath check_gobin
	@$(GOALIAS) get github.com/codegangsta/gin

$(COVER): check_gopath check_gobin
	@$(GOALIAS) get golang.org/x/tools/cmd/cover || exit 0

$(VET): check_gopath check_gobin
	@$(GOALIAS) get golang.org/x/tools/cmd/vet || exit 0

$(LINT): check_gopath check_gobin
	@$(GOALIAS) get golang.org/x/lint
