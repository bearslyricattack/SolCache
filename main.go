package main

import (
	"OrnnCache/cache/basefunction/baseclient"
	mysql2 "OrnnCache/cache/mysql"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Item 定义gorm model
type Item struct {
	Name    string `gorm:"column:name; type:varchar(32);"`
	Content string `gorm:"column:content; type:varchar(512);"`
	Author  string `gorm:"column:author; type:varchar(512);"`
	Type    int    `gorm:"column:type; type:varchar(32);"`
}

// TableName 定义表名
func (i Item) TableName() string {
	return "Item"
}

// client 演示
func main() {
	//新建数据库
	db := mysql2.New()
	//redis
	//redisClient := redis.New()
	////内存数据库
	base := baseclient.New()
	//新建数据库表,根据结构体创建
	//m := db.Migrator()
	//err := m.CreateTable(&Item{})
	//if err != nil {
	//	panic("新建表失败!")
	//}

	//创建路由
	r := gin.Default()

	//随机获取其中一个元素
	r.GET("/poem", func(c *gin.Context) {
		//获取随机key
		randomKey := base.RandomKey()
		//先查缓存
		res, err := base.Get(context.Background(), randomKey)

		if err {
			item := res.(Item)
			c.JSON(http.StatusOK, gin.H{
				"content": item.Content,
				"name":    item.Name,
				"author":  item.Author,
			})
			return
		}
		//缓存没查到，就去查mysql
		var mysql Item
		db.First(&mysql, "name = ?", randomKey)
		//查到之后 写回内存缓存
		base.Set(context.Background(), randomKey, mysql, 0)
		//返回
		c.JSON(http.StatusOK, gin.H{
			"content": mysql.Content,
			"name":    mysql.Content,
			"author":  mysql.Author,
		})
		return
	})
	//写入
	r.GET("/set", func(c *gin.Context) {
		name := c.Query("name")
		content := c.Query("content")
		author := c.Query("author")
		//组装对象
		item := &Item{
			Name:    name,
			Content: content,
			Author:  author,
		}
		//存入内存
		base.Set(context.Background(), name, item, 0)
		//存入mysql
		db.Create(&item)
	})
	//把数据库中的值存入缓存
	r.GET("/read", func(c *gin.Context) {
		var s []Item
		var swap Item
		db.Find(&s)
		for _, v := range s {
			swap = v
			base.Set(context.Background(), v.Name, swap, 0)
		}
	})
	r.GET("/poemH", func(c *gin.Context) {
		r.LoadHTMLGlob("test/*")
		//获取随机key
		randomKey := base.RandomKey()
		//先查缓存
		res, err := base.Get(context.Background(), randomKey)

		if err {
			item := res.(Item)
			c.HTML(http.StatusOK, "test.html", gin.H{
				"title":   item.Name,
				"content": item.Content,
				"author":  item.Author,
			})
			return
		}
		//缓存没查到，就去查mysql
		var mysql Item
		db.First(&mysql, "name = ?", randomKey)
		//查到之后 写回内存缓存
		base.Set(context.Background(), randomKey, mysql, 0)
		//返回
		c.HTML(http.StatusOK, "test.html", gin.H{
			"title":   mysql.Name,
			"content": mysql.Content,
			"author":  mysql.Author,
		})
		return
	})
	r.Run(":8081")
}
