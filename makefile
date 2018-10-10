.PHONY : list-deps install-deps
define GOPACKAGES 
github.com/360EntSecGroup-Skylar/excelize \
github.com/PuerkitoBio/goquery \
github.com/andybalholm/cascadia \
github.com/google/skylark \
github.com/mohae/deepcopy \
golang.org/x/net/html 
endef

default: install-deps

list-deps:
	go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'

install-deps:
	go get -v $(GOPACKAGES)
