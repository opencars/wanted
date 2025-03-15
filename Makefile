.PHONY: default all clean
APPS     := server worker pull
BLDDIR   := bin

.EXPORT_ALL_VARIABLES:
GO111MODULE  = on

default: clean all

all: $(APPS)

$(BLDDIR)/%:
	go build -o $@ ./cmd/$*

$(APPS): %: $(BLDDIR)/%

clean:
	@mkdir -p $(BLDDIR)
	@for app in $(APPS) ; do \
		rm -f $(BLDDIR)/$$app ; \
	done
