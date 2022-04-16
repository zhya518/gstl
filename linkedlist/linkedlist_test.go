package linkedlist

import (
	"testing"

	"github.com/guonaihong/gstl/must"
	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_RPush(t *testing.T) {
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().RPush(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_Len(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").Len(), 4)
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).Len(), 4)
}

func Test_RPop(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(3), []string{"2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).RPop(3), []int{2, 3, 4})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(10), []string{"1", "2", "3", "4"})
}

func Test_LPop(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(3), []string{"1", "2", "3"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).LPop(3), []int{1, 2, 3})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(10), []string{"1", "2", "3", "4"})
}

func Test_RangeSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangeSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"1", "2", "3", "4"})
}

func Test_RangePrevSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangePrevSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"4", "3", "2", "1"})
}

func Test_First(t *testing.T) {
	// 没有值
	assert.Error(t, must.TakeOneErr(New[string]().First()))
	// 有值
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").First()), "1")
	assert.Equal(t, must.TakeOne(New[int]().RPush(1, 2).First()), 1)
}

func Test_Last(t *testing.T) {
	// 没有值
	assert.Error(t, must.TakeOneErr(New[string]().Last()))
	//有值
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Last()), "4")
	assert.Equal(t, must.TakeOne(New[int]().RPush(1, 2).Last()), 2)
}

func Test_Get(t *testing.T) {
	// 正索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(0)), "1")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(1)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(3)), "4")

	// 负索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(-1)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(-2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(-3)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Get(-4)), "1")
}

func Test_Set(t *testing.T) {
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(11, "1991").ToSlice(), []string{"1", "2", "3", "4", "5", "6"})

	// 正索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(0, "1991").ToSlice(), []string{"1991", "2", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(1, "1991").ToSlice(), []string{"1", "1991", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(2, "1991").ToSlice(), []string{"1", "2", "1991", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(3, "1991").ToSlice(), []string{"1", "2", "3", "1991", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(4, "1991").ToSlice(), []string{"1", "2", "3", "4", "1991", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(5, "1991").ToSlice(), []string{"1", "2", "3", "4", "5", "1991"})

	// 负索引
}

func Test_Index(t *testing.T) {
	assert.Error(t, must.TakeOneErr(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(11)))

	// 正索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(0)), "1")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(1)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(3)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(4)), "5")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(5)), "6")

	// 负索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-1)), "6")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-2)), "5")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-3)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-4)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-5)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-6)), "1")
}

func Test_Remove(t *testing.T) {
	// 正索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(0).ToSlice(), []string{"2", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(1).ToSlice(), []string{"1", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(2).ToSlice(), []string{"1", "2", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(3).ToSlice(), []string{"1", "2", "3", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(4).ToSlice(), []string{"1", "2", "3", "4", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(5).ToSlice(), []string{"1", "2", "3", "4", "5"})

	// 负索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-1).ToSlice(), []string{"1", "2", "3", "4", "5"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-2).ToSlice(), []string{"1", "2", "3", "4", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-3).ToSlice(), []string{"1", "2", "3", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-4).ToSlice(), []string{"1", "2", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-5).ToSlice(), []string{"1", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-6).ToSlice(), []string{"2", "3", "4", "5", "6"})
}

func Test_LPush(t *testing.T) {
	assert.Equal(t, New[string]().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().LPush("1").ToSlice(), []string{"1"})
	assert.Equal(t, New[string]().LPush("2", "1").ToSlice(), []string{"1", "2"})
	assert.Equal(t, New[string]().LPush("3", "2", "1").ToSlice(), []string{"1", "2", "3"})
	assert.Equal(t, New[string]().LPush("4", "3", "2", "1").ToSlice(), []string{"1", "2", "3", "4"})
}

func Test_PushFront(t *testing.T) {
	assert.Equal(t, New[string]().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("1").ToSlice(), []string{"1"})
	assert.Equal(t, New[string]().PushFront("2", "1").ToSlice(), []string{"1", "2"})
	assert.Equal(t, New[string]().PushFront("3", "2", "1").ToSlice(), []string{"1", "2", "3"})
	assert.Equal(t, New[string]().PushFront("4", "3", "2", "1").ToSlice(), []string{"1", "2", "3", "4"})
}

func Test_Clear(t *testing.T) {
	assert.Equal(t, New[string]().PushFront("1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("2", "1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("3", "2", "1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("4", "3", "2", "1").Clear().ToSlice(), []string(nil))
}

func Test_IsEmpty(t *testing.T) {
	assert.True(t, New[string]().IsEmpty())
	assert.True(t, New[string]().PushFront("1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("2", "1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("3", "2", "1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("4", "3", "2", "1").Clear().IsEmpty())
}