package main

import (
	"log"
	"net/http"
	"sync"
)

/*
	【-------------------数量控制接口------替代redis---------------】
*/
var sum int64 = 0

//预存商品数量
var productNum int64 = 10000

//互斥锁
var mutex sync.Mutex

// GetOneProduct 获取秒杀商品 解决超卖问题
func GetOneProduct() bool {
	//加锁
	mutex.Lock()
	defer mutex.Unlock()
	//判断数据是否超限
	if sum < productNum {
		sum += 1
		return true
	}
	return false
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
	}
	w.Write([]byte("false"))
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("------Err:", err)
	}
}
