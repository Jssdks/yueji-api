package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var ctx = context.Background()

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "YUEJI API!")
		if err != nil {
			return
		}
	})

	// 使用dsn格式定义数据库连接信息
	dsn := "root:123456@tcp(localhost:3306)/yueji_test"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 执行查询
	rows, err := db.Query("SELECT * FROM dict_city LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 处理查询结果
	for rows.Next() {
		// ... 假设你的表有一些列，你需要定义相应的变量来接收
		// 例如 var id int, name string
		var id int64
		var cityCn string
		err := rows.Scan(&id, &cityCn)
		if err != nil {
			return
		}
		fmt.Println(id, cityCn)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected and queried the database!")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // Redis密码，如果没有则为空字符串
		DB:       0,                // 使用默认DB
	})

	// 设置键值
	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// 获取键值
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val) // 输出: key value

	//// 删除键
	//err = rdb.Del(ctx, "key").Err()
	//if err != nil {
	//	panic(err)
	//}

	// 关闭连接
	err = rdb.Close()
	if err != nil {
		panic(err)
	}
}
