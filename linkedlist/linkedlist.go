package linkedlist

import (
	"errors"
	"github.com/guonaihong/gstl/cmp"
)

var ErrListElemEmpty = errors.New("list is empty")
var ErrNotFound = errors.New("element not found")

// https://cs.opensource.google/go/go/+/go1.18.1:src/container/list/list.go
// https://github.com/torvalds/linux/blob/master/tools/include/linux/list.h
// https://redis.io/commands/?group=list
// LTRIM
// LINDEX -->done
// LINSERT -->done
// LLEN -->done
// LMOVE
// LPOP -->done
// LPOS
// LPUSH LPUSHX -->done
// LRANGE -->done
// LREM -->done
// LSET -->done
// RPOP -->done
// RPOPLPUSH
// RPUSH RPUSHX -->done
type LinkedList[T any] struct {
	root   Node[T]
	length int
}

// 每个Node节点, 包含前向和后向两个指针和数据域
type Node[T any] struct {
	next    *Node[T]
	prev    *Node[T]
	element T
}

// 返回一个双向循环链表
func New[T any]() *LinkedList[T] {
	return new(LinkedList[T]).Init()
}

// 指向自己, 组成一个环
func (l *LinkedList[T]) Init() *LinkedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.length = 0
	return l
}

// 返回长度
func (l *LinkedList[T]) Len() int {
	return l.length
}

// 延迟初始化
func (l *LinkedList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// at e at.next
// at <- e
//       e -> at.next
// at -> e
//       e <- at.next
func (l *LinkedList[T]) insert(at, e *Node[T]) {
	e.prev = at
	e.next = at.next
	at.next = e
	at.next.prev = e
	l.length++
}

func (l *LinkedList[T]) Append(other LinkedList[T]) {
	l.lazyInit()

}

// 类似redis lpop命令
// O(n)
func (l LinkedList[T]) LPop(count int) []T {
	if count <= 0 {
		return nil
	}

	count = cmp.Min(count, l.length)
	all := make([]T, count)
	i := 0

	l.RangeSafe(func(n *Node[T]) bool {
		if i == count {
			return true
		}

		all[i] = n.element
		l.remove(n)
		i++
		return false
	})
	return all
}

// 类似redis rpop命令
// O(n)
func (l *LinkedList[T]) RPop(count int) []T {
	if count <= 0 {
		return nil
	}

	count = cmp.Min(count, l.length)
	all := make([]T, count)
	l.RangePrevSafe(func(n *Node[T]) bool {
		if count <= 0 {
			return true
		}

		all[count-1] = n.element
		l.remove(n)
		count--
		return false
	})
	return all
}

// 从后向前遍历
func (l *LinkedList[T]) RangePrevSafe(callback func(n *Node[T]) bool) {

	var pos *Node[T]
	var n *Node[T]

	for pos, n = l.root.prev, pos.next; pos != &l.root; pos, n = n, pos.next {
		if callback(pos) {
			break
		}
	}
}

// 从前向后遍历
// callback 返回truek就退出遍历
func (l *LinkedList[T]) RangeSafe(callback func(n *Node[T]) (exit bool)) {

	var pos *Node[T]
	var n *Node[T]

	for pos, n = l.root.next, pos.next; pos != &l.root; pos, n = n, pos.next {
		if callback(pos) {
			break
		}
	}
}

// 类似redis lpush命令
func (l *LinkedList[T]) LPush(elems ...T) {
	l.PushFront(elems...)
}

// 往头位置插入
func (l *LinkedList[T]) PushFront(elems ...T) {
	l.lazyInit()
	for _, e := range elems {
		l.insert(&l.root, &Node[T]{element: e})
	}
}

// RPush是PushBack的同义词, 类似redis的RPush命令
func (l *LinkedList[T]) RPush(elems ...T) {
	l.PushBack(elems...)
}

// 往尾部的位置插入
func (l *LinkedList[T]) PushBack(elems ...T) {
	l.lazyInit()
	for _, e := range elems {
		l.insert(l.root.prev, &Node[T]{element: e})
	}
}

// 返回第1个元素
func (l *LinkedList[T]) First() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.next.element, nil
}

// 返回最后1个元素
func (l *LinkedList[T]) Last() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.prev.element, nil
}

// 链表是否为空
func (l *LinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

// 清空链表 O(n)
func (l *LinkedList[T]) Clear() {

	l.RangeSafe(func(n *Node[T]) bool {
		l.remove(n)
		return false
	})
}

// 类似于redis linsert after 命令
func (l *LinkedList[T]) InsertAfter(value T, equal func(value T) bool) {
	l.RangeSafe(func(n *Node[T]) bool {
		if equal(n.element) {
			l.insert(n, &Node[T]{element: value})
			return true
		}
		return false
	})
}

//  类似于redis linsert before 命令
func (l *LinkedList[T]) InsertBefore(value T, equal func(value T) bool) {
	l.RangeSafe(func(n *Node[T]) bool {
		if equal(n.element) {
			l.insert(n.prev, &Node[T]{element: value})
			return true
		}
		return false
	})
}

// 查找是否包含这个value
func (l *LinkedList[T]) ContainsFunc(value T, cb func(value T) bool) bool {
	for pos := l.root.next; pos != &l.root; pos = pos.next {
		if cb(pos.element) {
			return true
		}
	}
	return false
}

// 通过索引查找是否包含这个value
func (l *LinkedList[T]) Get(idx int) (e T, err error) {
	if idx >= 0 {
		return l.Index(idx)
	}
	err = ErrNotFound
	return
}

// 删除这个元素
// n.prev n n.next
// n.prev --> n.next
// n.prev <-- n.next
func (l *LinkedList[T]) remove(n *Node[T]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	l.length--
}

// 类似redis lrem命令
// count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
// count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
// count = 0 : 移除表中所有与 VALUE 相等的值。

// 返回值
// 被删除元素个数
func (l *LinkedList[T]) RemFunc(value T, count int, cb func(value T) bool) (ndel int) {
	var (
		pos *Node[T]
		n   *Node[T]
		i   int
	)

	if count >= 0 {
		for pos, n = l.root.next, pos.next; pos != &l.root; pos, n = n, pos.next {
			if count == 0 || i <= count {
				if cb(pos.element) {
					l.remove(pos)
					ndel++
				}
				i++
			}
		}

		return
	}

	count = -count
	for pos, n = l.root.prev, pos.prev; pos != &l.root; pos, n = n, pos.prev {
		if count == 0 || i <= count {
			if cb(pos.element) {
				l.remove(pos)
				ndel++
			}
			i++
		}
	}

	return
}

// Index 通过索引查找是否包含这个value
// 和redis lindex命令类似,
// idx >= 0 形为和Get(idx int) 一样, 获取指向索引的元素
// idx < 0 获取倒数第几个元素
// O(min(index, length - index))
func (l *LinkedList[T]) Index(idx int) (e T, err error) {
	var n *Node[T]
	n, err = l.indexInner(idx)
	if err != nil {
		return
	}
	return n.element, nil
}

// 类型redis lset命令
// index >= 0 正着数
// index < 0 倒着数
// On(min(index, length - index))
func (l *LinkedList[T]) Set(index int, value T) {
	n, err := l.indexInner(index)
	if err != nil {
		return
	}
	n.element = value
}

func (l *LinkedList[T]) indexInner(idx int) (*Node[T], error) {
	idx, front := l.index(idx)

	if front {
		for pos, i := l.root.next, 0; pos != &l.root; pos, i = pos.next, i+1 {
			if i == idx {
				return pos, nil
			}
		}
	} else {
		for pos, i := l.root.prev, idx; pos != &l.root; pos, i = pos.prev, i-1 {
			if i == 0 {
				return pos, nil
			}
		}

	}

	return nil, ErrNotFound

}

// 计算索引
func (l *LinkedList[T]) index(idx int) (newIdx int, front bool) {
	length := l.length

	// 转正索引
	if idx < 0 {
		idx = idx + length
	}

	// 如果倒序遍历元素更少
	other := length - idx
	if idx > other {
		return other, false
	}

	// 如果正序遍历元素更少
	return idx, true
}

// 删除指定索引的元素, 效率 min(O(index), O(len - index))
func (l *LinkedList[T]) Remove(index int) {
	l.removeInner(index)
}

// list 转成slice , 效率O(n)
func (l *LinkedList[T]) ToSlice() []T {
	if l.length == 0 {
		return nil
	}

	rv := make([]T, 0, l.length)
	for pos := l.root.next; pos != &l.root; pos = pos.next {
		rv = append(rv, pos.element)
	}

	return rv
}

func (l *LinkedList[T]) removeInner(index int) {
	var (
		pos *Node[T]
		n   *Node[T]
		i   int
	)

	idx, front := l.index(index)

	if front {
		for pos, n = l.root.next, pos.next; pos != &l.root; pos, n, i = n, pos.next, i+1 {
			if i == idx {
				l.remove(pos)
				return
			}
		}
	}

	i = idx
	for pos, n = l.root.prev, pos.prev; pos != &l.root; pos, n, i = n, pos.prev, i-1 {
		if i == 0 {
			l.remove(pos)
			return
		}
	}

}

// 打印
// range 类似redis lrange命令
func (l *LinkedList[T]) Range(pr func(value T), startAndEnd ...int) {

	start := 0
	end := 0

	if len(startAndEnd) > 0 {
		start = startAndEnd[0]
	}

	if len(startAndEnd) > 1 {
		start = startAndEnd[1]
	}

	i := 0
	l.rangeStartEndSafe(start, end, func(n *Node[T]) bool {

		if len(startAndEnd) != 0 {
			if i >= start && i <= end {
				pr(n.element)
			}

			if i > end {
				return true
			}
			return false
		}

		pr(n.element)
		i++
		return false
	})
}

func (l *LinkedList[T]) rangeStartEndSafe(start, end int, callback func(n *Node[T]) (exit bool)) {

	if start < 0 {
		start += l.length
		if start < 0 {
			start = 0
		}
	}

	if end < 0 {
		end += l.length
	}

	if start > end || start >= l.length {
		return
	}
	var pos *Node[T]
	var n *Node[T]

	for pos, n = l.root.next, pos.next; pos != &l.root; pos, n = n, pos.next {
		if callback(pos) {
			break
		}
	}
}

// 类似于redis ltrim命令, 对列表进行裁剪
func (l *LinkedList[T]) Trim(start, end int) *LinkedList[T] {

	i := 0
	l.rangeStartEndSafe(start, end, func(n *Node[T]) bool {
		if i < start || i > end {
			l.remove(n)
		}
		i++
		return true
	})
	return l
}
