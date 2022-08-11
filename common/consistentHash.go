package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//申明新切片类型
type uints []uint32

// Len 返回切片长度
func (x uints) Len() int {
	return len(x)
}

// Less 比对两个数大小
func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

// Swap 切片两个值交换
func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//当hash环上没有数据时，提示错误
var errEmpty = errors.New("Hash环上没有数据")

type ConsistentHash struct {
	//hash环，key 为哈希值，值存放节点的信息
	circle map[uint32]string
	//已经排序的节点hash切片
	sortedHashes uints // uints 自定义Hash切片类型
	//虚拟节点个数，用来增加hash的平衡性
	VirtualNode int
	//map 读写锁
	sync.RWMutex
}

// NewConsistentHash 创建一致性Hash算法结构体，设置默认节点数量
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		//初始变量
		circle: make(map[uint32]string),
		//设置虚拟节点数量
		VirtualNode: 20,
	}
}

//自动生成Key值
func (c *ConsistentHash) generateKey(element string, index int) string {
	//副本Key生成
	return element + strconv.Itoa(index)
}

//获取Hash位置
func (c *ConsistentHash) hashKey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组中
		copy(srcatch[:], key)
		//使用IEEE 多项式返回数据的CRC-32校验和 以中国国际标准
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//
func (c *ConsistentHash) updateSortedHashes() {
	hashesLocation := c.sortedHashes[:0]
	//判断切片容量，是否过大，如果过大就重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashesLocation = nil
	}

	//添加hashes
	for k := range c.circle {
		hashesLocation = append(hashesLocation, k)
	}
	//对所有节点Hash值进行排序
	//方便之后进行二分法查找
	sort.Sort(hashesLocation)

}

// Add 向hash环中添加节点
func (c *ConsistentHash) Add(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.add(element)
}

//添加节点
func (c *ConsistentHash) add(element string) {
	//循环虚拟节点，设置副本
	for i := 0; i < c.VirtualNode; i++ {
		//获取key
		key := c.generateKey(element, i)
		//根据key获取hash位置
		hashLocation := c.hashKey(key)
		//将位置映射在Hash环上
		c.circle[hashLocation] = element
	}
	//对生成以后的虚拟节点 更新排序
	c.updateSortedHashes()
}

// Remove 向hash环中删除节点
func (c *ConsistentHash) Remove(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.remove(element)
}

//删除节点
func (c *ConsistentHash) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		//获取key
		key := c.generateKey(element, i)
		//根据key获取hash位置
		hashLocation := c.hashKey(key)
		//将hash环上的位置删除
		c.circle[hashLocation] = element
		delete(c.circle, hashLocation)
	}
	//对删除以后的虚拟节点 更新排序
	c.updateSortedHashes()
}

//顺时针查找最近的服务器节点
func (c *ConsistentHash) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用“二分查找”算法搜索指定切片，满足条件的最小值
	index := sort.Search(len(c.sortedHashes), f)
	//如果超出范围 设置i=0
	if index > len(c.sortedHashes) {
		index = 0
	}
	return index
}

// Get 根据数据标识获取最近的服务器节点信息
func (c *ConsistentHash) Get(name string) (string, error) {
	//加读锁
	c.RLock()
	//解锁
	defer c.Unlock()
	//如果为0返回错误
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//计算机Hash值
	key := c.hashKey(name)
	//得到节点location
	index := c.search(key)
	hashLocation := c.sortedHashes[index]
	return c.circle[hashLocation], nil
}
