SHELL=/bin/bash

EXE = chatGPT-gateway

all: $(EXE)

chatGPT-gateway:
	@echo "building $@ ..."
	$(MAKE) -s -f make.inc s=static

clean:
	rm -f $(EXE)

