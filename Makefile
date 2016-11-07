# -*- coding:utf-8 -*-
.PHONY: get-deps
build:
	go install bitbucket.org/jixiuf/fund/db
	go install bitbucket.org/jixiuf/fund/defs
	go install bitbucket.org/jixiuf/fund/dt
	go install bitbucket.org/jixiuf/fund/utils
	go install bitbucket.org/jixiuf/fund/eastmoney
	go install bitbucket.org/jixiuf/fund/main/datainit
	go install bitbucket.org/jixiuf/fund/main/dailyupdate
	go install bitbucket.org/jixiuf/fund/main/rank_period
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
