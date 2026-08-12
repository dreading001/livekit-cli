package portaudio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWASAPISharedAutoConvert(t *testing.T) {
	t.Run("not requested leaves every host API unchanged", func(t *testing.T) {
		assert.NoError(t, validateWASAPISharedAutoConvert(false, false, "ALSA"))
	})

	t.Run("requested accepts WASAPI", func(t *testing.T) {
		assert.NoError(t, validateWASAPISharedAutoConvert(true, true, "Windows WASAPI"))
	})

	t.Run("requested rejects a non-WASAPI device clearly", func(t *testing.T) {
		err := validateWASAPISharedAutoConvert(true, false, "Windows DirectSound")
		assert.EqualError(t, err, "portaudio: WASAPI shared-mode output auto-conversion requested, but selected output device uses Windows DirectSound")
	})
}
