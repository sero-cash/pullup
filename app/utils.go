package app

import (
	"encoding/binary"
	"fmt"
	"github.com/robfig/cron"
	"math/big"
	"sync/atomic"
)

func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

func decodeNumber(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func AddJob(spec string, run RunFunc) *cron.Cron {
	c := cron.New()
	c.AddJob(spec, &RunJob{run: run})
	c.Start()
	return c
}

type (
	RunFunc func()
)

type RunJob struct {
	runing int32
	run    RunFunc
}

func (r *RunJob) Run() {
	x := atomic.LoadInt32(&r.runing)
	if x == 1 {
		return
	}

	atomic.StoreInt32(&r.runing, 1)
	defer func() {
		atomic.StoreInt32(&r.runing, 0)
	}()

	r.run()
}

func NewBigIntFromString(s string, base int) (*big.Int, error) {
	val, flag := big.NewInt(0).SetString(s, base)
	if !flag {
		return nil, fmt.Errorf("can't convert %s to BigInt", s)
	}
	return val, nil
}
