.PHONY: all build

all: build

build:
	go build -tags exclude_graphdriver_btrfs,btrfs_noversion,containers_image_openpgp -o bin/cleargate ./cmd/cleargate
