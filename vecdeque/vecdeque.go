package vecdeque

import (
	"fmt"
	"github.com/guonaihong/gstl/cmp"
	"math"
)

// 参考文档如下
// https://doc.rust-lang.org/std/collections/struct.VecDeque.html
// https://doc.rust-lang.org/src/alloc/collections/vec_deque/mod.rs.html
// 翻译好的中文文档
// https://rustwiki.org/zh-CN/src/alloc/collections/vec_deque/mod.rs.html

const (
	INITIAL_CAPACITY uint = 7 // 2^3 - 1
	MINIMUM_CAPACITY uint = 1 // 2 - 1
)

type VecDeque[T any] struct {
	//总是指向可以读取的第一个元素
	//head只是指向应该写入数据的位置
	// 如果tail == head, 则缓存区为空. 环形缓冲区的长度定义为两者之间的距离
	tail uint
	head uint
	buf  []T
}

// 初始化
func New[T any]() *VecDeque[T] {
	return &VecDeque[T]{}
}

// 初始VecDeque, 并设置实际需要的容量
func WithCapacity[T any](capacity int) *VecDeque[T] {
	cap := nextPowOfTwo(cmp.Max(uint(capacity)+1, MINIMUM_CAPACITY+1))
	return &VecDeque[T]{buf: make([]T, cap, cap)}
}

// 如果缓冲区满了. 就返回true
func (v *VecDeque[T]) IsFull() bool {
	return v.Cap()-v.Len() == 1
}

//TODO
func (v *VecDeque[T]) Len() int {
	return 0
}

func (v *VecDeque[T]) grow() *VecDeque[T] {

	return v
}

// 将一个元素添加到VecDeque 后面
func (v *VecDeque[T]) PushBack(value T) {
	if v.IsFull() {
		v.grow()
	}
	head := v.head
	v.head = v.wrapAdd(v.head, uint(1))

	v.buf[head] = value

	// 先检查是否满了
	// 没有满就扩容
	// 修改head值
}

func (v *VecDeque[T]) Get(i int) {

}

func (v *VecDeque[T]) cap() int {
	return len(v.buf)
}

func (v *VecDeque[T]) Cap() int {
	return v.cap() - 1
}

func (v *VecDeque[T]) wrapAdd(index uint, addend uint) uint {
	return v.wrapIndex(index + addend)
}

func (v *VecDeque[T]) wrapIndex(index uint) uint {
	return wrapIndex(index, uint(v.cap()))
}

func wrapIndex(index uint, size uint) uint {
	// 判断size是否是2的n次方
	if n := (size & (size - 1)); n != 0 {
		panic(fmt.Sprintf("size is always a power of 2, the current size is %d", size))
	}

	return index & (size - 1)
}

// TODO 优化下
func nextPowOfTwo(n uint) uint {

	for i := 1; i < 32; i++ {

		if nextPowOfTwoNum := math.Pow(2, float64(i)); nextPowOfTwoNum > float64(n) {
			return uint(nextPowOfTwoNum)
		}
	}

	return 0
}