package rhash

import (
	"errors"
	"github.com/cespare/xxhash/v2"
	"math"
	"reflect"
	"unsafe"
)

const (
	HT_INITIAL_EXP  = 2
	HT_INITIAL_SIZE = (1 << (HT_INITIAL_EXP))
)

var ErrHashing = errors.New("rehashing...")
var ErrSize = errors.New("wrong size")

// 元素
type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

// hash 表头
type Hash[K comparable, V any] struct {
	table   [2][]*entry[K, V] //hash table
	used    [2]uint64         // 记录每个table里面存在的元素个数
	sizeExp [2]int8           //记录exp

	rehashidx int
	keySize   int //key的长度
	hashFunc  func(str string) uint64
	isKeyStr  bool //是string类型的key, 或者不是
}

// 初始化一个hashtable
func New[K comparable, V any]() *Hash[K, V] {
	return &Hash[K, V]{
		rehashidx: -1,
		hashFunc:  xxhash.Sum64String,
	}
}

// 初始化一个hashtable并且可以设置值
func NewWithHashFunc[K comparable, V any](hashFunc func(str string) uint64) *Hash[K, V] {
	h := New[K, V]()
	h.hashFunc = hashFunc
	return h
}

// 保存key的类型和key的长度
func (h *Hash[K, V]) keyTypeAndKeySize() {
	var k K
	switch (interface{})(k).(type) {
	case string:
	default:
		h.keySize = int(unsafe.Sizeof(k))
	}
}

// 计算hash值
func (h *Hash[K, V]) calHash(k K) uint64 {
	var key string

	if h.isKeyStr {
		// 直接赋值会报错, 使用unsafe绕过编译器检查
		key = *(*string)(unsafe.Pointer(&k))
	} else {
		// 因为xxhash.Sum64String 接收string, 所以要把非string类型变量当成string类型来处理
		key = *(*string)(unsafe.Pointer(&reflect.StringHeader{
			Data: uintptr(unsafe.Pointer(&k)),
			Len:  h.keySize,
		}))
	}

	return xxhash.Sum64String(key)
}

func (h *Hash[K, V]) isRehashing() bool {
	return h.rehashidx != -1
}

// TODO 这个函数可以优化下
func nextExp(size uint64) int8 {
	if size >= math.MaxUint64 {
		return 63
	}
	e := int8(HT_INITIAL_EXP)
	for {
		if 1<<e >= size {
			return e
		}
		e++
	}
	return e
}

// 手动修改hashtable的大小
func (h *Hash[K, V]) Resize(size uint64) error {
	// 如果正在扩容中, 或者需要扩容的数据小于已存在的元素, 直接返回
	if h.isRehashing() || h.used[0] > uint64(size) {
		return ErrHashing
	}

	newSizeExp := nextExp(uint64(size))
	// 新大小比需要的大小还小
	newSize := uint64(1 << newSizeExp)
	if newSize < size {
		return ErrSize
	}

	// 新扩容大小和以前的一样
	if uint64(newSizeExp) == uint64(h.sizeExp[0]) {
		return nil
	}

	newTable := make([]*entry[K, V], newSize)

	// 第一次初始化
	if h.table[0] == nil {
		h.sizeExp[0] = newSizeExp
		h.table[0] = newTable
		return nil
	}

	// 把新hash表放到table[1]里面
	h.sizeExp[1] = newSizeExp
	h.used[1] = 0
	h.table[1] = newTable
	h.rehashidx = 0
	return nil
}

// 收缩hash table
func (h *Hash[K, V]) ShrinkToFit() error {
	if h.isRehashing() {
		return ErrHashing
	}

	minimal := h.used[0]
	if minimal < HT_INITIAL_SIZE {
		minimal = HT_INITIAL_SIZE
	}

	return h.Resize(minimal)
}

// 返回索引值和entry
func (h *Hash[K, V]) findIndexAndEntry() (i int, e *entry[K, V]) {

	return
}

func (h *Hash[K, V]) rehash(n int) error {
	// 控制桶的个数
	emptyVisits := n * 10

	// 正在rehashing 就退出
	if h.isRehashing() {
		return ErrHashing
	}

	// n是控制桶数
	for ; n > 0 && h.used[0] != 0; n-- {

		for h.table[0][h.rehashidx] == nil {
			h.rehashidx++
			emptyVisits--
			if emptyVisits == 0 {
				return nil
			}
		}

		// 取出hash槽中第一个元素
		head := h.table[0][h.rehashidx]
		for head != nil {
			next := head.next
			newIdx := h.calHash(head.key) & sizeMask(h.sizeExp[1])
			head.next = h.table[1][newIdx]
			h.table[1][newIdx] = head
			h.used[0]--
			h.used[1]++
			head = next
		}

		h.table[0][h.rehashidx] = nil
		h.rehashidx++
	}

	if h.used[0] == 0 {
		h.table[0] = h.table[1]
		h.used[0] = h.used[1]
		h.sizeExp[0] = h.sizeExp[1]
		h.reset(1)
		h.rehashidx = -1
	}
	return nil
}

func (h *Hash[K, V]) reset(idx int) {
	h.table[idx] = nil
	h.sizeExp[idx] = -1
	h.used[idx] = 0
}

func sizeMask(exp int8) uint64 {
	if exp == -1 {
		return 0
	}

	return (1 << exp) - 1
}

// 获取
func (h *Hash[K, V]) Get(key K) (v V, err error) {

	return
}

// 获取
func (h *Hash[K, V]) GetOrZero(key K) (v V) {
	return
}

// 遍历
func (h *Hash[K, V]) Range() {

}

// 设置
func (h *Hash[K, V]) Set(k K, v V) {
}

// 删除
func (h *Hash[K, V]) Delete(key K) {

}
