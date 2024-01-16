package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	incrNum       uint64
	defaultPrefix string = "trace-id"
	currTime             = time.Now().UnixMilli()
)

// NewTraceID New trace id
func NewTraceID(prefix ...string) string {
	tmp := defaultPrefix
	if 0 < len(prefix) {
		tmp = prefix[0]
	}
	now := time.Now().UnixMilli()
	return fmt.Sprintf("%s-%d-%d",
		tmp,
		now,
		getSeqence(now))
}

func getSeqence(now int64) uint64 {
	old := currTime
	if old == now {
		return atomic.AddUint64(&incrNum, 1)
	}

	if atomic.CompareAndSwapInt64(&currTime, old, now) {
		atomic.StoreUint64(&incrNum, 0)
	}
	return atomic.AddUint64(&incrNum, 1)
}

func GoID() int64 {
	//非常慢，非debug不要使用
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine")
	)
	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}
	return int64(id)
}
