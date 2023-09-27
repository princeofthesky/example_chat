package main

import (
	"github.com/liyue201/gostl/utils/comparator"
	"github.com/princeofthesky/example_chat/skiplist"
	"time"
)

type t struct {
	a int
}

func testarray() {
	a := []t{}
	i := 0
	go func() {
		for {
			i++
			if i%10 != 0 {
				a = append(a, t{a: 1})
			} else {
				println("a", len(a))
				c := make([]t, len(a)-7)
				copy(c, a[7:])
				a = c
				println("c", len(a))
			}
		}
	}()
	t := 0
	go func() {
		for {
			c := 0
			println(len(a))
			for i := 0; i < len(a); i++ {
				c = c + a[i].a
			}
			t = c
		}
	}()

	time.Sleep(time.Minute)
	println(t)
}
func main() {

	max := 1000000
	a := map[int]int{}
	b := skiplist.New[int, int](comparator.IntComparator)
	start := time.Now().UnixMilli()
	for i := 0; i < max; i++ {
		a[i] = i
	}
	end := time.Now().UnixMilli() - start
	start = time.Now().UnixMilli()
	for i := 0; i < max; i++ {
		b.Insert(i, i)
	}
	done := time.Now().UnixMilli() - start

	println("done ", end, done)
}
