BIN_NAME=yeti
LDFLAGS = -ldflags "-s -w"

export GOARCH = amd64
export CGO_ENABLED = 0
export GO111MODULE = on

define build-yeti
go build ${LDFLAGS} -o dist/${GOOS}_${GOARCH}/${BIN_NAME} main.go
endef

define upx-yeti
upx dist/${GOOS}_${GOARCH}/${BIN_NAME}
endef

linux_build: export GOOS = linux
linux_build:
	$(build-yeti)

linux_upx: export GOOS = linux
linux_upx:
	$(upx-yeti)

linux: linux_build linux_upx
	