//go:build !windows

package portaudio

import (
	"runtime"
	"unsafe"
)

func outputStreamHostInfo(_ *DeviceInfo, options OutputStreamOptions) (unsafe.Pointer, func(), error) {
	noop := func() {}
	err := validateWASAPISharedAutoConvert(options.WASAPISharedAutoConvert, false, runtime.GOOS+" (not WASAPI)")
	return nil, noop, err
}
