package views

import (
	"testing"
)

func BenchmarkSeq(b *testing.B) {
	reportTPS := func(b *testing.B) {
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	}

	b.Run("seq", func(b *testing.B) {
		v := &NoopView{}
		for b.Loop() {
			_, _ = v.Call(nil)
		}
		reportTPS(b)
	})

	b.Run("par", func(b *testing.B) {
		v := &NoopView{}
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = v.Call(nil)
			}
		})
		reportTPS(b)
	})

}
