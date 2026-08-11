package console

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/livekit/livekit-cli/v2/pkg/apm"
)

func TestAECStatsReturnsNilBeforeSpeakerLoopSnapshot(t *testing.T) {
	pipeline := &AudioPipeline{}

	require.Nil(t, pipeline.AECStats())
}

func TestAECStatsReturnsIndependentSpeakerLoopSnapshot(t *testing.T) {
	pipeline := &AudioPipeline{}
	want := apm.Stats{
		EchoReturnLossEnhancement: 12.5,
		DelayMs:                   43,
		HasERLE:                   true,
		HasDelay:                  true,
	}
	pipeline.apmStats = &want

	got := pipeline.AECStats()
	require.Equal(t, want, *got)
	got.DelayMs = 999

	require.Equal(t, 43, pipeline.AECStats().DelayMs)
}

func TestAECStatsSnapshotIsSafeForConcurrentTUIReads(t *testing.T) {
	pipeline := &AudioPipeline{}
	const iterations = 1_000
	var wait sync.WaitGroup
	var invalidSnapshot atomic.Bool
	wait.Add(2)

	go func() {
		defer wait.Done()
		for i := range iterations {
			stats := apm.Stats{DelayMs: i, HasDelay: true}
			pipeline.mu.Lock()
			pipeline.apmStats = &stats
			pipeline.mu.Unlock()
		}
	}()
	go func() {
		defer wait.Done()
		for range iterations {
			if stats := pipeline.AECStats(); stats != nil {
				if !stats.HasDelay {
					invalidSnapshot.Store(true)
				}
			}
		}
	}()

	wait.Wait()
	require.False(t, invalidSnapshot.Load())
	require.NotNil(t, pipeline.AECStats())
}
