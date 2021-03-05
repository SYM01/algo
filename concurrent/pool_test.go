package concurrent_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sym01/algo/concurrent"
)

func TestPoolWithOneSlot(t *testing.T) {
	pool := concurrent.NewActionPool(1)
	buf := new(bytes.Buffer)
	expected := new(bytes.Buffer)

	for i := 0; i < 100; i++ {
		c := strconv.Itoa(i) + "."
		expected.WriteString(c)
		pool.Do(concurrent.Action{
			Runner: func() []concurrent.Action {
				buf.WriteString(c)
				return nil
			},
		})
	}
	pool.WaitAndClose()

	if bytes.Compare(buf.Bytes(), expected.Bytes()) != 0 {
		t.Fatalf("unexpected result, expect \n%s , got \n%s", expected, buf)
	}
}

func TestPoolWithOneSlotAndDeps(t *testing.T) {
	pool := concurrent.NewActionPool(1)
	buf := new(bytes.Buffer)
	expected := new(bytes.Buffer)

	for i := 0; i < 100; i++ {
		c := strconv.Itoa(i) + "."
		p := new(int32)

		expected.WriteString("2@")
		expected.WriteString(c)

		pool.Do(concurrent.Action{
			Runner: func() []concurrent.Action {
				buf.WriteString(strconv.Itoa(int(atomic.LoadInt32(p))))
				buf.WriteString("@")
				buf.WriteString(c)
				return nil
			},
			Dependencies: []concurrent.Action{
				{
					Runner: func() []concurrent.Action {
						atomic.AddInt32(p, 1)
						return nil
					},
				},
				{
					Runner: func() []concurrent.Action {
						atomic.AddInt32(p, 1)
						return nil
					},
				},
			},
		})
	}
	pool.WaitAndClose()

	if bytes.Compare(buf.Bytes(), expected.Bytes()) != 0 {
		t.Fatalf("unexpected result, expect \n%s , got \n%s", expected, buf)
	}
}

func TestPoolConcurrency(t *testing.T) {
	pool := concurrent.NewActionPool(15)
	time2sleep := time.Second * 2
	tStart := time.Now()

	for i := 0; i < 5; i++ {
		pool.Do(concurrent.Action{
			Runner: func() []concurrent.Action {
				// 2s will elapse
				time.Sleep(time2sleep)

				// 2s will elapse
				return []concurrent.Action{
					{
						Runner: func() []concurrent.Action {
							time.Sleep(time2sleep)
							return nil
						},
					},
					{
						Runner: func() []concurrent.Action {
							time.Sleep(time2sleep)
							return nil
						},
					},
					{
						Runner: func() []concurrent.Action {
							time.Sleep(time2sleep)
							return nil
						},
					},
				}
			},
		})
	}
	pool.WaitAndClose()
	timeElapsed := time.Since(tStart)
	if timeElapsed >= 3*time2sleep {
		t.Fatalf("unexpected result, expect finished in about %s, got %s elapsed", time2sleep*2, timeElapsed)
	}
}

func TestPanic(t *testing.T) {
	pool := concurrent.NewActionPool(0)
	pool.Map(nil, nil)
	pool.Do(concurrent.Action{})
	pool.WaitAndClose()
	expectPanic(func() {
		pool.Do(concurrent.Action{})
	}, t)

	pool = concurrent.NewActionPool(0)
	buf := new(bytes.Buffer)
	pool.SetLogger(buf)
	pool.Do(concurrent.Action{
		Runner: func() []concurrent.Action {
			panic("panic test")
		},
	})
	pool.Wait()
	if !strings.HasPrefix(buf.String(), "panic test") {
		t.Fatalf("unexpected result, expect logger with panic msg")
	}
}

func expectPanic(f func(), t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("unexpected result, expect panic, but nothing happened")
		}
	}()

	f()
}

func ExampleActionPool_Map() {
	pool := concurrent.NewActionPool(10)
	input := []interface{}{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	}
	transFunc := func(i interface{}) interface{} {
		return 12 - i.(int)
	}

	ret := pool.Map(input, transFunc)
	for _, item := range ret {
		fmt.Println(item)
	}

	// Output:
	// 11
	// 10
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
}
