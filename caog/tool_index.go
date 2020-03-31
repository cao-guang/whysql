package main

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "caog",
		Alias:     "caog",
		BuildTime: time.Date(2020, 3, 30, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/cao-guang/whysql",
		Summary:   "caog工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "caog",
		Hidden:    true,
	},
	{
		Name:      "genrc",
		Alias:     "caog-gen-rc",
		BuildTime: time.Date(2020, 3, 30, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/cao-guang/whysql",
		Summary:   "rc缓存代码生成",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "caog",
	},

}
