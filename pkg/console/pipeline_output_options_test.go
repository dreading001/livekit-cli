package console

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWASAPIOutputAutoConvertPreservesConsoleVoiceBus(t *testing.T) {
	pipeline, err := NewPipeline(PipelineConfig{WASAPIOutputAutoConvert: true})
	require.NoError(t, err)
	require.True(t, pipeline.wasapiOutputAutoConvert)

	require.Equal(t, 48000, SampleRate)
	require.Equal(t, 1, Channels)
	require.Equal(t, 1440, SamplesPerFrame)
	require.Equal(t, 480, APMFrameSamples)
}
