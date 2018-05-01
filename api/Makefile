BIN = api
BIN_DIR = ./bin
PKG_DIR = ./pkg

all: fmt vet test
	@mkdir -p ${BIN_DIR}
	@rm -fr ${BIN_DIR}/*
	@env CGO_ENABLED=0 go build -v -o ${BIN_DIR}/${BIN}

fmt:
	@go fmt ./...

help:
	@echo "default\t\ttest, format, and build the code"
	@echo "fmt\t\tformat the code"
	@echo "package\t\tbuild/package the code for platforms"
	@echo "test\t\ttest the code"
	@echo "vet\t\trun go vet"

package:
	@echo "\n==> Packaging\n"
	@rm -fr ${PKG_DIR}
	@mkdir -p ${PKG_DIR}
	@echo "\n==> darwin"
	@env CGO_ENABLED=0 GOOS=darwin go build -v -o ${PKG_DIR}/darwin/${BIN}
	@zip -9 --junk-paths ${PKG_DIR}/darwin/${BIN}.zip ${PKG_DIR}/darwin/${BIN}
	@echo "\n==> linux"
	@env CGO_ENABLED=0 GOOS=linux go build -v -o ${PKG_DIR}/linux/${BIN}
	@zip -9 --junk-paths ${PKG_DIR}/linux/${BIN}.zip ${PKG_DIR}/linux/${BIN}

test:
	@env go test -v ./...

vet:
	@go vet -v ./...

.PNONY: all fmt help package test vet
