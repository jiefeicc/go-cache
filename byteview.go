package go_cache

// ByteView 缓存值的抽象与封装
type ByteView struct {
	b []byte
}

// Len 返回数据视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 以切片的形式返回数据的副本
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 以字符串形式返回数据，必要时进行复制
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
