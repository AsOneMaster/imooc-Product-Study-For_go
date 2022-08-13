package main

import (
	"encoding/json"
	"fmt"
	"imooc-Product/common"
	"imooc-Product/datamodels"
	"imooc-Product/encrypt"
	"imooc-Product/rabbitmq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

/*
	分布式 集群
	分布式权限验证
	秒杀规则整合
*/
import (
	"errors"
)

//设置集群地址，最好内网IP
var hostArray = []string{"127.0.0.1", "127.0.0.1"} //手动指定

var localHost = "" //使用动态获取本地IP
// GetOneIP 数量控制接口服务器内网IP，或者getOne的SLB内网IP
var GetOneIP = "127.0.0.1"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.ConsistentHash

//消息队列 验证
var rabbitMqValidate *rabbitmq.RabbitMQ

// AccessControl 用来存放控制信息
type AccessControl struct {
	//用户想要存放的信息
	sourcesArray map[int]interface{}
	//高并发下 保证数据安全
	sync.RWMutex
}

//创建全局变量
var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

// GetNewRecord 获取定制的数据
func (a *AccessControl) GetNewRecord(uid int) interface{} {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	data := a.sourcesArray[uid]
	return data
}

// SetNewRecord 设置记录
func (a *AccessControl) SetNewRecord(uid int) {
	a.RWMutex.Lock()
	a.sourcesArray[uid] = "hello imooc"
	a.RWMutex.Unlock()
}

// GetDistributedRight 获取分布式权限
func (a *AccessControl) GetDistributedRight(r *http.Request) bool {
	uid, err := r.Cookie("userid")
	if err != nil {
		return false
	}
	//采用一致性算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return a.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据 返回结果
		return GetDataFromOtherMap(hostRequest, r)
	}
}

// GetDataFromMap 获取本机map，并且处理业务逻辑，返回的结果类型为bool类型
func (a *AccessControl) GetDataFromMap(uid string) (isOK bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := a.GetNewRecord(uidInt)
	fmt.Println("-------GetDataFromMap:-----", data)
	//执行判断逻辑
	//测试使用
	if data != nil {
		return true
	}
	return true
}

// GetDataFromOtherMap 获取其他节点处理结果
func GetDataFromOtherMap(host string, r *http.Request) bool {

	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, r)
	if err != nil {
		return false
	}

	//获取uid
	//uidPre, err := r.Cookie("userid")
	//if err != nil {
	//	return false
	//}
	////获取sign
	//uidSign, err := r.Cookie("sign")
	//if err != nil {
	//	return false
	//}
	////模拟接口访问
	//client := &http.Client{}
	//req, err := http.NewRequest("Get", "http://"+host+":"+"port"+"/access", nil)
	//if err != nil {
	//	return false
	//}
	////手动指定，排查多余Cookies
	//cookieUid := &http.Cookie{Name: "userid", Value: uidPre.Value, Path: "/"}
	//cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	////添加Cookie到模拟的请求
	//r.AddCookie(cookieUid)
	//r.AddCookie(cookieSign)
	//
	////获取返回结果
	//response, err := client.Do(req)
	//if err != nil {
	//	return false
	//}
	//body, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return false
	//}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

// GetCurl 模拟http请求
func GetCurl(hostUrl string, r *http.Request) (response *http.Response, body []byte, err error) {
	//获取uid
	uidPre, err := r.Cookie("userid")
	if err != nil {
		return
	}
	//获取sign
	uidSign, err := r.Cookie("sign")
	if err != nil {
		return
	}
	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("Get", hostUrl, nil)
	if err != nil {
		return
	}
	//手动指定，排查多余Cookies
	cookieUid := &http.Cookie{Name: "userid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加Cookie到模拟的请求
	r.AddCookie(cookieUid)
	r.AddCookie(cookieSign)

	//获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}

// CheckRight 分离权限验证
func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// Check  执行正常业务逻辑[--------------最重要----------]
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	fmt.Println("执行check！")
	/*
		url.ParseQuery()
		Code:
		m, err := url.ParseQuery(`x=1&y=2&y=3`) if err != nil {     log.Fatal(err) } fmt.Println(toJSON(m))
		Output:
		{"x":["1"], "y":["2", "3"]}
	*/
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["productID"]) <= 0 && len(queryForm["productID"][0]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println("Validate check:", productString)
	//获取用户cookie
	userCookie, err := r.Cookie("userid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//1. 分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}
	//2. 获取数量控制权限，防止秒杀超卖现象
	hostUrl := "https://" + GetOneIP + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态 调用消息队列
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//整合下单
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				//真实场景还要获取 用户信息等
				w.Write([]byte("false"))
				return
			}
			//获取商品ID，获取用户ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				//真实场景还要获取 用户信息等
				w.Write([]byte("false"))
				return
			}
			//消息体
			message := datamodels.Message{ProductID: productID, UserID: userID}
			//消息体类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				//真实场景还要获取 用户信息等
				w.Write([]byte("false"))
				return
			}
			//生产消息 发送下单消息 下单操作由消费端服务器接受并操作
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				//真实场景还要获取 用户信息等
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

// Auth 统一验证拦截器
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserInfo 身份校验函数 检测Cookie信息
func CheckUserInfo(r *http.Request) error {
	//获取id，cookie
	uidCookie, err := r.Cookie("userid")
	if err != nil {
		return errors.New("用户UserID Cookie 获取失败")
	}
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败")
	}
	//对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串被篡改")
	}
	//进行比对
	fmt.Println("CheckUserInfo--------", "用户ID：", uidCookie.Value, "解密后用户ID：", string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

//解密后的信息比对函数
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}
func main() {
	//负载均衡器设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistentHash()
	//采用一致性哈希节点，添加节点 服务器哈希环的添加
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	//自动获取本地IP
	localIP, err := common.GetInterfaceIP()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIP
	//注册消息队列服务
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocPoduct")
	defer rabbitMqValidate.Destory()

	fmt.Println("本地ip：", localHost)
	//1. 过滤器
	filter := common.NewFilter()
	//2. 注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	filterHandlerFunc := filter.Handle(Check)
	filterHandlerFuncRight := filter.Handle(CheckRight)
	//3. 启动服务
	http.HandleFunc("/check", filterHandlerFunc)
	http.HandleFunc("/checkRight", filterHandlerFuncRight)
	http.ListenAndServe(":8083", nil)

}
