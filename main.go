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

func p(w http.ResponseWriter, s ...interface{}) {
	_, err := fmt.Fprintln(w, s...)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p(w, "YUEJI API!")
		dsn := "root:6NXPqMm9JZgj3UnK@tcp(localhost:33080)/yueji_test"
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
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		p(w, "DB连接成功")
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
		for rows.Next() {
			var id int64
			var cityCn string
			err := rows.Scan(&id, &cityCn)
			if err != nil {
				panic(err)
			}
			p(w, id, cityCn)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		p(w, "DB读取成功")
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "k7z9x*t[j=M^5){e",
			DB:       7,
		})
		defer func(rdb *redis.Client) {
			err = rdb.Close()
			if err != nil {
				panic(err)
			}
		}(rdb)
		p(w, "REDIS连接成功")
		err = rdb.Set(ctx, "go-key", time.Now().String(), 0).Err()
		if err != nil {
			panic(err)
		}
		val, err := rdb.Get(ctx, "go-key").Result()
		if err != nil {
			panic(err)
		}
		_, err = fmt.Fprintln(w, "go-key", val)
		if err != nil {
			return
		}
		p(w, "REDIS读取成功")
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

}
