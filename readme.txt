* 抓取数据存在mysql
    需要mysql 数据库，对应的配置在 defs/fund_defs.go中
    var DBConfig = dt.DBConfig{Host: "127.0.0.1", User: "root", Pass: "root!!@@__))", Name: "fund", CharSet: "utf8", Port: "3306"}
    这个过程会持续几个小时，

    main/datainit/main.go
    cd datainit go run main.go
* 计算定投排行
  main/rank_period/main.go
