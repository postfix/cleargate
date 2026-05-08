.PHONY: all build-web build-backend

all: build-web build-backend

build-web:
	cd web && npm install && npm run build

build-backend:
	go build -tags exclude_graphdriver_btrfs,btrfs_noversion,containers_image_openpgp -o bin/cleargate ./cmd/cleargate
