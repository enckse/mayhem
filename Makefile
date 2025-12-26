GOFLAGS := -trimpath -buildmode=pie -mod=readonly -modcacherw -buildvcs=false
TARGET  := dist
VERSION ?= "$(shell git describe --abbrev=0 --tags)-$(shell git log -n 1 --format=%h)"
OBJECT  := $(TARGET)/mayhem

.PHONY: $(OBJECT)

all: setup $(OBJECT)

setup:
	@mkdir -p $(TARGET)

$(OBJECT):
	go build $(GOFLAGS) -ldflags "$(LDFLAGS) -X main.version=$(VERSION)" -o "$(OBJECT)" cmd/mayhem.go

unittest:
	go test ./...

check: unittest $(OBJECT)
	find tests/ -type f -name "*.db" -delete
	$(OBJECT) version
	cat tests/objects.json | XDG_CACHE_HOME=tests/testdata $(OBJECT) import --config tests/settings.toml
	cat tests/objects.json | XDG_CACHE_HOME=tests/testdata $(OBJECT) import --config tests/settings.toml 2>&1 | grep 'import not supported into existing database'
	cat tests/objects.json | XDG_CACHE_HOME=tests/testdata $(OBJECT) import --config tests/settings.toml --overwrite
	cat tests/objects.json | sed 's/Section/ZZZ/g' | XDG_CACHE_HOME=tests/testdata $(OBJECT) merge --config tests/settings.toml
	XDG_CACHE_HOME=tests/testdata $(OBJECT) export --config tests/settings.toml > tests/testdata/results.json
	diff -u tests/testdata/results.json tests/expected.json
	diff -u tests/testdata/mayhem/todo.json tests/expected.json

clean:
	rm -f "$(OBJECT)"
	find internal/ tests/ -type f -wholename "*testdata*" -delete
	find internal/ tests/ -type d -empty -delete
