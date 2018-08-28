.PHONY: rocket install clean re dir re_all all test release docs
.PHONY: darwin_386 darwin_amd64
.PHONY: linux_arm6 linux_arm7 linux_arm64 linux_386 linux_amd64 linux_mips linux_mips64
.PHONY: windows_386 windows_amd64
.PHONY: freebsd_386 freebsd_amd64 freebsd_arm6 freebsd_arm7

NAME = rocket
DIST_DIR = dist
REPO="github.com/astrocorp42/rocket"
VERSION := $(shell cat version/version.go| grep "\sVersion" | cut -d '"' -f2)

define checksums
	echo $$(openssl sha512 $(1) | cut -d " " -f2) $$(echo $(1) | rev | cut -d "/" -f1 | rev) >> $(2)/sha512sum$(3)
endef

define build_for_os_arch
	mkdir -p $(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION)
	GOARM=$(4) GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -a -installsuffix cgo \
		 -ldflags "-X $(REPO)/version.UTCBuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
		 -X $(REPO)/version.GitCommit=`git rev-parse HEAD` \
		 -X $(REPO)/version.GoVersion=`go version | cut -d' ' -f 3 | cut -c3-`" \
		 -o $(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION)/$(NAME)$(3)
	@# binary checksums
	$(call checksums,$(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION)/$(NAME)$(3),$(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION),.txt)

	zip -j $(DIST_DIR)/$(NAME)_$(VERSION)_$(1)_$(2)$(4).zip $(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION)/$(NAME)$(3) \
		$(DIST_DIR)/$(1)_$(2)$(4)/$(VERSION)/sha512sum.txt

	@#archive checksums
	$(call checksums,$(DIST_DIR)/$(NAME)_$(VERSION)_$(1)_$(2)$(4).zip,dist,s.txt)
endef


$(NAME): dir
	CGO_ENABLED=0 go build -a -installsuffix cgo \
		 -ldflags "-X $(REPO)/version.UTCBuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
		 -X $(REPO)/version.GitCommit=`git rev-parse HEAD` \
		 -X $(REPO)/version.GoVersion=`go version | cut -d' ' -f 3 | cut -c3-`" \
		 -o $(DIST_DIR)/$(NAME)

test:
	go vet $(go list ./... | grep -v /vendor/)
	go test -v -race ./...

install:
	CGO_ENABLED=0 go install -a -installsuffix cgo \
		 -ldflags "-X $(REPO)/version.UTCBuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
		 -X $(REPO)/version.GitCommit=`git rev-parse HEAD` \
		 -X $(REPO)/version.GoVersion=`go version | cut -d' ' -f 3 | cut -c3-`"

clean:
	rm -rf $(NAME) $(DIST_DIR)

re: clean $(NAME)

dir:
	mkdir -p $(DIST_DIR)

re_all: clean all

all: darwin_386 darwin_amd64 linux_arm6 linux_arm7 linux_arm64 linux_386 linux_amd64 linux_mips linux_mips64 windows_386 windows_amd64 freebsd_386 freebsd_amd64 freebsd_arm6 freebsd_arm7

release: clean
	git tag v$(VERSION)
	git push origin v$(VERSION)

docs:
	cd _docs && mkdocs build

darwin_386:
	$(call build_for_os_arch,darwin,386,)

darwin_amd64:
	$(call build_for_os_arch,darwin,amd64,)

linux_arm6:
	$(call build_for_os_arch,linux,arm,,6)

linux_arm7:
	$(call build_for_os_arch,linux,arm,,7)

linux_arm64:
	$(call build_for_os_arch,linux,arm64,)

linux_386:
	$(call build_for_os_arch,linux,386,)

linux_amd64:
	$(call build_for_os_arch,linux,amd64,)

linux_mips:
	$(call build_for_os_arch,linux,mips,)

linux_mips64:
	$(call build_for_os_arch,linux,mips64,)

windows_386:
	$(call build_for_os_arch,windows,386,.exe)

windows_amd64:
	$(call build_for_os_arch,windows,amd64,.exe)

freebsd_386:
	$(call build_for_os_arch,freebsd,386,)

freebsd_amd64:
	$(call build_for_os_arch,freebsd,amd64,)

freebsd_arm6:
	$(call build_for_os_arch,freebsd,arm,,6)

freebsd_arm7:
	$(call build_for_os_arch,freebsd,arm,,7)
