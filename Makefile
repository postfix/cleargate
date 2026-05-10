.PHONY: all build-web build-backend run clean

# Build tags required for macOS / Podman bindings without gpgme and btrfs
BUILD_TAGS = -tags exclude_graphdriver_btrfs,btrfs_noversion,containers_image_openpgp

all: build-web build-backend

build-web:
	cd web && npm install && npm run build

build-backend:
	go build $(BUILD_TAGS) -o bin/cleargate ./cmd/cleargate

run:
	go run $(BUILD_TAGS) ./cmd/cleargate --tools-dir=./tools

clean:
	rm -f bin/cleargate
	rm -f cleargate.db cleargate.db.wal
