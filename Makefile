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
	$(OBJECT) version --config tests/settings.toml

clean:
	rm -f "$(OBJECT)"
	find internal/ tests/ -type f -wholename "*testdata*" -delete
	find internal/ tests/ -type d -empty -delete
