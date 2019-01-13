package visitercontrol

// import (
// 	"errors"
// )

// type CircleQueue struct {
// 	maxSize int           //比实际队列长度大1
// 	slice   []interface{} //切片会被实际队列长度大1
// 	head    int           //头
// 	tail    int           //尾
// }

// //初始化环形队列
// func NewCircleQueue(size int) *CircleQueue {
// 	var c CircleQueue
// 	c.maxSize = size + 1
// 	c.slice = make([]interface{}, c.maxSize)
// 	return &c
// }

// //入对列
// func (this *CircleQueue) Push(val interface{}) (err error) {
// 	if this.IsFull() {
// 		return errors.New("queue is full")
// 	}
// 	this.slice[this.tail] = val
// 	this.tail = (this.tail + 1) % this.maxSize
// 	return
// }

// //出对列
// func (this *CircleQueue) Pop() (val interface{}, err error) {
// 	if this.IsEmpty() {
// 		return 0, errors.New("queue is empty")
// 	}
// 	val = this.slice[this.head]
// 	this.head = (this.head + 1) % this.maxSize
// 	return
// }

// //判断队列是否已满
// func (this *CircleQueue) IsFull() bool {
// 	return (this.tail+1)%this.maxSize == this.head
// }

// //判断队列是否为空
// func (this *CircleQueue) IsEmpty() bool {
// 	return this.tail == this.head
// }

// //判断有多少个元素
// func (this *CircleQueue) Size() int {
// 	return (this.tail + this.maxSize - this.head) % this.maxSize
// }

// //转化为slice输出，方便range
// func (this *CircleQueue) ToSlice() []interface{} {
// 	size := this.Size()
// 	ret := make([]interface{}, 0, size)
// 	if size == 0 {
// 		return ret
// 	}
// 	temphead := this.head
// 	for i := 0; i < size; i++ {
// 		ret = append(ret, this.slice[temphead])
// 		temphead = (temphead + 1) % this.maxSize
// 	}
// 	return ret
// }
