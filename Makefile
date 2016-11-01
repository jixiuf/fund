# -*- coding:utf-8 -*-
.PHONY: get-deps
get-deps:

	mkdir -p $$GOPATH/src/golang.org/x
	if [ ! -d $$GOPATH/src/golang.org/x/net ]; then \
		cd $$GOPATH/src/golang.org/x;git clone  https://github.com/golang/net.git ;cd -; \
	else \
		cd $$GOPATH/src/golang.org/x/net;git pull;cd -; \
	fi

	mkdir -p $$GOPATH/src/golang.org/x
	if [ ! -d $$GOPATH/src/golang.org/x/text ]; then \
		cd $$GOPATH/src/golang.org/x;git clone  https://github.com/golang/text.git ;cd -; \
	else \
		cd $$GOPATH/src/golang.org/x/text;git pull;cd -; \
	fi

	go get github.com/PuerkitoBio/goquery
