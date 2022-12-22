package main

import (
	"4byte_calculator/types"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"os"
	"sync"
	"time"
)

var config types.Config
var done bool
var result string

const CMDCalculate4bytes = "cal4bytes"

func main() {
	cmd := Cmd(CMDCalculate4bytes)
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}

func Cmd(programName string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   programName,
		Short: fmt.Sprintf("%s %s", CMDCalculate4bytes, programName),
		Long:  fmt.Sprintf("%s %s", CMDCalculate4bytes, programName),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Check(); err != nil {
				return err
			}

			caEngine := types.NewEngine(config.Regular, config.Target, config.ThreadNum, config.ThreadWorkRange)
			var wg sync.WaitGroup

			p := mpb.New(mpb.WithWaitGroup(&wg))
			numBars := int(caEngine.ProcessesNum)

			for i := 0; i < numBars; i++ {
				createBar(caEngine, &wg, p)
			}

			go func() {
				for true {
					if done {
						p.Shutdown()
						return
					}
				}
			}()
			p.Wait()

			fmt.Println("complete! result is: " + result)
			return nil
		},
	}

	cmd.Flags().StringVarP(&config.Regular, "regular", "", "", "regular, example: buy_%d(bytes32,uint256)")
	cmd.Flags().StringVarP(&config.Target, "target", "", "", "target 4bytes value, example: 0x00000000")
	cmd.Flags().UintVarP(&config.ThreadNum, "thread", "", 5, "thread number")
	cmd.Flags().Uint64VarP(&config.ThreadWorkRange, "thread.range", "", 1000000000, "number of computations per thread")

	return cmd
}

func createBar(caEngine *types.Engine, wg *sync.WaitGroup, p *mpb.Progress) {
	wg.Add(1)
	start, end := caEngine.ApplyRange()
	process := types.NewProcess(caEngine.Regular, caEngine.Target, start, end)

	name := fmt.Sprintf("%s#[%d-%d]:", caEngine.Regular, start, end)
	bar := p.AddBar(
		int64(caEngine.ProcessWorkRange),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.Name("calculating", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncWidth), "done",
			),
		),
	)

	// simulating some work
	go func() {
		defer wg.Done()
		//rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		//max := 100 * time.Millisecond
		for i := start; i < end; i++ {
			// start variable is solely for EWMA calculation
			// EWMA's unit of measure is an iteration's duration
			start := time.Now()
			//time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
			n, isDone, function := process.CalOnce()
			// we need to call EwmaIncrement to fulfill ewma decorator's contract
			//bar.EwmaIncrement(time.Since(start))
			if isDone {
				done = true
				result = function
			}
			if done {
				i = end
				bar.EwmaIncrInt64(int64(end-i-1), time.Since(start))
			}
			if end-i == 1 {
				createBar(caEngine, wg, p)
			}
			bar.EwmaIncrInt64(int64(n), time.Since(start))
		}
	}()
}
