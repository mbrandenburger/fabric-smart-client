package benchmark

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

type Config struct {
	NumWorker       int
	Duration        time.Duration
	WarumupDuration time.Duration
	MaxSleep        int
}

func Benchmark(cfg Config, f func()) error {
	var (
		numWorker      = cfg.NumWorker
		duration       = cfg.Duration
		warmupDuration = cfg.WarumupDuration
		maxSleep       = cfg.MaxSleep

		hists []*hdrhistogram.Histogram
		mu    sync.Mutex

		wg sync.WaitGroup
	)
	wg.Add(numWorker)

	start := time.Now()
	warmDeadline := start.Add(warmupDuration)
	endDeadline := warmDeadline.Add(duration)

	for range numWorker {
		go func() {
			defer wg.Done()
			lh := hdrhistogram.New(1, 5_000_000, 3)

			for {
				start := time.Now()
				if start.After(endDeadline) {
					mu.Lock()
					hists = append(hists, lh)
					mu.Unlock()
					return
				}
				f()
				elapsed := time.Since(start)
				if start.After(warmDeadline) {
					_ = lh.RecordValue(elapsed.Microseconds())
				}

				if maxSleep > 0 {
					time.Sleep(time.Duration(rand.IntN(maxSleep)))
				}
			}
		}()
	}
	wg.Wait()
	lh := hdrhistogram.New(1, 5_000_000, 3)
	for _, h := range hists {
		lh.Merge(h)
	}

	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	fmt.Println(printResult(lh, duration))

	return nil
}

func printResult(lh *hdrhistogram.Histogram, runtime time.Duration) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Total ops:\t%v\n", lh.TotalCount())
	fmt.Fprintf(&sb, "TPS:\t\t%v\n", int(float64(lh.TotalCount())/runtime.Seconds()))
	fmt.Fprintf(&sb, "Avg Latency:\t%v\n", toDuration(int64(lh.Mean())))
	fmt.Fprintf(&sb, "P99:\t\t%v\n", toDuration(lh.ValueAtPercentile(99)))
	fmt.Fprintf(&sb, "P95:\t\t%v\n", toDuration(lh.ValueAtPercentile(95)))
	fmt.Fprintf(&sb, "P50:\t\t%v\n", toDuration(lh.ValueAtPercentile(50)))
	fmt.Fprintf(&sb, "stddev:\t\t%v\n", toDuration(int64(lh.StdDev())))
	fmt.Fprintf(&sb, "max:\t\t%v\n", toDuration(lh.Max()))
	fmt.Fprintf(&sb, "min:\t\t%v\n", toDuration(lh.Min()))

	return sb.String()
}

func toDuration(usecs int64) time.Duration {
	return time.Duration(usecs * 1000)
}
