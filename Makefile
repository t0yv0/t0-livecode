.PHONY: start
start:
	nix build .#assets
	(cd www && git clean -fxd)
	cp result/www/* www/
	nix build
	STATE_DIRECTORY=$(shell pwd)/data result/bin/t0-livecode
