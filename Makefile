SHELL=/bin/bash

EXE=db-proxy

all:
	@echo "building $(EXE) ..."
	@$(MAKE) -s -f make.inc s=static

clean:
	rm -f $(EXE)
