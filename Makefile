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

ifndef tnum
gen-test-files:
	@echo "tnum is not set"
	@exit 1
else ifeq (,$(wildcard test/$(tnum)/in.yaml))
gen-test-files:
	@echo "test/$(tnum)/in.yaml does not exist"
	@exit 1
else
gen-test-files: $(PROGRAM)
	<test/$(tnum)/in.yaml $< -y       > test/$(tnum)/out-y.yaml
	<test/$(tnum)/in.yaml $< -y -p    > test/$(tnum)/out-y-p.yaml
	<test/$(tnum)/in.yaml $< -j       > test/$(tnum)/out-j.yaml
	<test/$(tnum)/in.yaml $< -j -p    > test/$(tnum)/out-j-p.yaml
	<test/$(tnum)/in.yaml $< -t       > test/$(tnum)/out-t.yaml
	<test/$(tnum)/in.yaml $< -t -p    > test/$(tnum)/out-t-p.yaml
	<test/$(tnum)/in.yaml $< -e       > test/$(tnum)/out-e.yaml
	<test/$(tnum)/in.yaml $< -e -p    > test/$(tnum)/out-e-p.yaml
	<test/$(tnum)/in.yaml $< -e -c    > test/$(tnum)/out-e-c.yaml
	<test/$(tnum)/in.yaml $< -e -p -c > test/$(tnum)/out-e-p-c.yaml
	<test/$(tnum)/in.yaml $< -t       > test/$(tnum)/out-t.yaml
	<test/$(tnum)/in.yaml $< -t -p    > test/$(tnum)/out-t-p.yaml
	<test/$(tnum)/in.yaml $< -t -c    > test/$(tnum)/out-t-c.yaml
	<test/$(tnum)/in.yaml $< -t -p -c > test/$(tnum)/out-t-p-c.yaml
	<test/$(tnum)/in.yaml $< -n       > test/$(tnum)/out-n.yaml
endif


$(PROGRAM): $(DEPS)
	go mod tidy
	go fmt
	go build

$(GO-YAML-REPO-DIR):
	git clone -q --depth 1 $(GO-YAML-REPO-URL) $@
	(cd $@ && patch -p1 < ../$(GO-YAML-PATCH))
