package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "YUEJI API!")
		if err != nil {
			return
		}
		// 使用dsn格式定义数据库连接信息
		dsn := "root:123456@tcp(localhost:3306)/yueji_test"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)

		// 测试连接
		err = db.Ping()
		if err != nil {
			panic(err)
		}

		// 执行查询
		rows, err := db.Query("SELECT id,city_cn FROM dict_city LIMIT 10")
		if err != nil {
			panic(err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				panic(err)
			}
		}(rows)
		_, err = fmt.Fprintln(w, "DB 连接成功")
		if err != nil {
			return
		}
		// 处理查询结果
		for rows.Next() {
			// ... 假设你的表有一些列，你需要定义相应的变量来接收
			// 例如 var id int, name string
			var id int64
			var cityCn string
			err := rows.Scan(&id, &cityCn)
			if err != nil {
				panic(err)
			}
			_, err = fmt.Fprintln(w, id, cityCn)
			if err != nil {
				return
			}
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprintln(w, "DB读取成功")
		if err != nil {
			return
		}
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",   // Redis地址
			Password: "k7z9x*t[j=M^5){e", // Redis密码，如果没有则为空字符串
			DB:       7,                  // 使用默认DB
		})
		_, err = fmt.Fprintln(w, "REDIS连接成功")
		if err != nil {
			return
		}
		// 设置键值
		err = rdb.Set(ctx, "go-key", time.Now().String(), 0).Err()
		if err != nil {
			panic(err)
		}
		// 获取键值
		val, err := rdb.Get(ctx, "go-key").Result()
		if err != nil {
			panic(err)
		}
		_, err = fmt.Fprintln(w, "go-key", val)
		if err != nil {
			return
		}

		//// 删除键
		//err = rdb.Del(ctx, "key").Err()
		//if err != nil {
		//	panic(err)
		//}
		_, err = fmt.Fprintln(w, "DB读取成功")
		if err != nil {
			return
		}
		// 关闭连接
		err = rdb.Close()
		if err != nil {
			panic(err)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

}
