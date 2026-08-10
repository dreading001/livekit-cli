//go:build windows

package portaudio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWASAPIAutoConvertStreamInfo(t *testing.T) {
	spec, err := newWASAPIAutoConvertSpec()
	require.NoError(t, err)
	require.NotZero(t, spec.Size)
	require.Equal(t, 13, spec.HostAPIType) // paWASAPI
	require.Equal(t, uint(1), spec.Version)
	require.Equal(t, uint(1<<6), spec.Flags) // paWinWasapiAutoConvert
}

func TestWASAPIAutoConvertSelectedDevice(t *testing.T) {
	query := os.Getenv("LIVEKIT_TEST_WASAPI_OUTPUT_DEVICE")
	if query == "" {
		t.Skip("set LIVEKIT_TEST_WASAPI_OUTPUT_DEVICE to run the physical output open/write proof")
	}

	require.NoError(t, Initialize())
	t.Cleanup(func() { require.NoError(t, Terminate()) })

	device, err := FindDevice(query, false)
	require.NoError(t, err)
	stream, err := OpenOutputStreamWithOptions(device, 48000, 1, 1440, OutputStreamOptions{
		WASAPISharedAutoConvert: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, stream.Close()) })

	require.NoError(t, stream.Start())
	require.NoError(t, stream.Write(make([]int16, 1440)))
	require.NoError(t, stream.Stop())
}
