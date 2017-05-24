# -*- coding:utf-8 -*-
.PHONY: get-deps
.PHONY: vendor
build:
	go install github.com/jixiuf/fund/db
	go install github.com/jixiuf/fund/defs
	go install github.com/jixiuf/fund/dt
	go install github.com/jixiuf/fund/utils
	go install github.com/jixiuf/fund/eastmoney
	go install github.com/jixiuf/fund/main/datainit
	go install github.com/jixiuf/fund/main/dailyupdate
	go install github.com/jixiuf/fund/main/rank_period
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
vendor:
	if [ -z `command -v godep` ] ; then \
		go get github.com/tools/godep;\
	fi


	if [ -d vendor ] ; then  \
		godep update ./...; \
	else \
		godep save ./...; \
	fi


