package types

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

type Process struct {
	Regular      string
	StartIndex   uint64
	EndIndex     uint64
	NowIndex     uint64
	Target       string
	spinner      *spinner.Spinner
	Success      bool
	SuccessIndex uint64
}

func NewProcess(regular string, target string, start uint64, end uint64) *Process {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	return &Process{
		Regular:    regular,
		StartIndex: start,
		NowIndex:   start,
		EndIndex:   end,
		Target:     target,
		spinner:    s,
	}
}
func (p *Process) State() string {
	return fmt.Sprintf(" " + p.Regular + fmt.Sprintf("[%d-%d]", p.StartIndex, p.EndIndex))
}

func (p *Process) Stop() error {
	return nil
}

func (p *Process) Start() error {
	p.NowIndex = p.StartIndex
	for true {
		if p.NowIndex > p.EndIndex {
			break
		}
		function := fmt.Sprintf(p.Regular, p.NowIndex)
		d := crypto.Keccak256Hash([]byte(function))

		fourByte := d.String()[:10]

		result := fmt.Sprintf("%s: %s zero num: %d", function, fourByte)
		if fourByte == p.Target {
			p.Success = true
			p.SuccessIndex = p.NowIndex
			p.spinner.FinalMSG = "Complete!\nNew line!\nAnother one!\n" + result
			p.spinner.Stop()
			break
		}
		p.NowIndex++
	}
	return nil
}

func (p *Process) CalOnce() (uint64, bool, string) {
	if p.NowIndex > p.EndIndex {
		return 0, false, ""
	}
	function := fmt.Sprintf(p.Regular, p.NowIndex)
	d := crypto.Keccak256Hash([]byte(function))

	fourByte := d.String()[:10]

	result := fmt.Sprintf("%s: %s zero num: %d", function, fourByte)
	if fourByte == p.Target {
		p.Success = true
		p.SuccessIndex = p.NowIndex
		p.spinner.FinalMSG = "Complete!\nNew line!\nAnother one!\n" + result
		p.spinner.Stop()
		return p.EndIndex - p.NowIndex, true, function
	}

	p.NowIndex++
	return 1, false, ""
}
