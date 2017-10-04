VERSION=0.1.0

.PHONY: all
all: dspool-bsd

.PHONY: clean
clean:
	rm dspool-bsd

.PHONY: distclean
distclean: maintainer-clean

.PHONY: maintainer-clean
maintainer-clean:
	git clean -fxd

.PHONY: dist
dist:
	git archive -o dspool-${VERSION}.tar.gz --prefix=dspool/ HEAD

dspool-bsd:
	GOOS=freebsd go build -o dspool-bsd felesatra.moe/dspool/cmd/dspool
