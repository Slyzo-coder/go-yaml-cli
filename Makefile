# Using the "Makes" Makefile setup - https://github.com/makeplus/makes
M := $(or $(MAKES_REPO_DIR),.cache/makes)
$(shell [ -d $M ] || git clone -q https://github.com/makeplus/makes $M)
include $M/init.mk
include $M/clean.mk
MAKES-NO-RULES := true
include $M/go.mk
include $M/shell.mk

override PATH := $(ROOT):$(PATH)
export PATH

PROGRAM := go-yaml

TEST-FILES := $(wildcard test/*.yaml)

ifneq (,$(file))
TEST-FILES := $(file)
endif

GO-YAML-PATCH := go-yaml-patch
GO-YAML-REPO-URL := https://github.com/yaml/go-yaml
GO-YAML-REPO-DIR := .go-yaml

MAKES-CLEAN := $(PROGRAM)
MAKES-REALCLEAN := $(GO-YAML-REPO-DIR)
MAKES-DISTCLEAN := .cache

DEPS := $(GO) $(wildcard *.go) $(GO-YAML-REPO-DIR)


default::

test:: $(PROGRAM)
	go $@$(if $(v), -v,)

build:: $(PROGRAM)

fmt: $(DEPS)
	go $@

tidy:: $(DEPS)
	go mod $@

ifndef PREFIX
install:
	$(error PREFIX is not set)
else
install: $(PROGRAM)
	mkdir -p $(PREFIX)/bin
	install -m 0755 $(PROGRAM) $(PREFIX)/bin/$(PROGRAM)
endif

update-patch: $(GO-YAML-REPO-DIR)
	git -C $< diff > $(GO-YAML-PATCH)

gen-test-files: $(PROGRAM)
	bin/gen-test-files


$(PROGRAM): $(DEPS)
	go mod tidy
	go fmt
	go build

$(GO-YAML-REPO-DIR):
	git clone -q --depth 1 $(GO-YAML-REPO-URL) $@
	(cd $@ && patch -p1 < ../$(GO-YAML-PATCH))
