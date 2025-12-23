GOFLAGS := -trimpath -buildmode=pie -mod=readonly -modcacherw -buildvcs=false
TARGET  := dist
VERSION ?= "$(shell git describe --abbrev=0 --tags)-$(shell git log -n 1 --format=%h)"
OBJECT  := $(TARGET)/mayhem

.PHONY: $(OBJECT)

all: setup $(OBJECT)

setup:
	@mkdir -p $(TARGET)

$(OBJECT):
	go build $(GOFLAGS) -ldflags "$(LDFLAGS) -X main.version=$(VERSION)" -o "$(OBJECT)" main.go

clean:
	rm -f "$(OBJECT)"
