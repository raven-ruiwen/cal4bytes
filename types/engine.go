package types

import (
	"github.com/judwhite/go-svc"
	"sync"
)

type Engine struct {
	ProcessesNum     uint
	processes        []*Process
	ProcessWorkRange uint64
	workRangeMax     uint64
	locker           *sync.Mutex
	Regular          string
	Target           string
}

func NewEngine(regular string, target string, num uint, workRange uint64) *Engine {
	var m sync.Mutex
	return &Engine{
		Regular:          regular,
		Target:           target,
		ProcessesNum:     num,
		processes:        nil,
		ProcessWorkRange: workRange,
		locker:           &m,
	}
}

func (e *Engine) Init(svc.Environment) error {
	return nil
}
func (e *Engine) Stop() error {
	return nil
}

func (e *Engine) Start() error {
	return nil
}

func (e *Engine) ApplyRange() (uint64, uint64) {
	e.locker.Lock()
	defer e.locker.Unlock()

	start := e.workRangeMax
	end := start + e.ProcessWorkRange

	e.workRangeMax = e.workRangeMax + e.ProcessWorkRange

	return start, end
}
