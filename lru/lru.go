package lru

import "container/list"

// Cache 包含字典和双向链表的结构体类型 Cache，方便实现后续的增删查改操作。
// lru 缓存淘汰策略
type Cache struct {
	// 允许使用的最大内存
	maxBytes int64
	// 当前已使用的内存
	nbytes int64
	// Go 语言标准库实现的双向链表list.List
	ll *list.List
	// 键是字符串，值是双向链表中节点型指针。
	cache map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为 nil。
	OnEvicted func(key string, value Value)
}

// 键值对 entry 是双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

/*
Value 接口
为了通用性，我们允许值是实现了 Value 接口的任意类型。
该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
*/
type Value interface {
	Len() int
}

// Len 方法, Cache 类实现 Len 方法，返回双向链表中节点的 len
func (c *Cache) Len() int {
	return c.ll.Len()
}

// New 方便实例化 Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add 新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 获取 value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除 “最近最少使用的值”
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
