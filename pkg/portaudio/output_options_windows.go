//go:build windows

package portaudio

/*
#include <stdlib.h>
#include <portaudio.h>
#include <pa_win_wasapi.h>

static PaHostApiTypeId lk_device_host_api_type(PaDeviceIndex device) {
	const PaDeviceInfo *deviceInfo = Pa_GetDeviceInfo(device);
	if (deviceInfo == NULL) {
		return paInDevelopment;
	}
	const PaHostApiInfo *hostInfo = Pa_GetHostApiInfo(deviceInfo->hostApi);
	if (hostInfo == NULL) {
		return paInDevelopment;
	}
	return hostInfo->type;
}

static const char *lk_device_host_api_name(PaDeviceIndex device) {
	const PaDeviceInfo *deviceInfo = Pa_GetDeviceInfo(device);
	if (deviceInfo == NULL) {
		return NULL;
	}
	const PaHostApiInfo *hostInfo = Pa_GetHostApiInfo(deviceInfo->hostApi);
	if (hostInfo == NULL) {
		return NULL;
	}
	return hostInfo->name;
}

static PaWasapiStreamInfo *lk_new_wasapi_auto_convert_stream_info(void) {
	PaWasapiStreamInfo *info = (PaWasapiStreamInfo *)calloc(1, sizeof(PaWasapiStreamInfo));
	if (info == NULL) {
		return NULL;
	}
	info->size = sizeof(PaWasapiStreamInfo);
	info->hostApiType = paWASAPI;
	info->version = 1;
	info->flags = paWinWasapiAutoConvert;
	return info;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func outputStreamHostInfo(device *DeviceInfo, options OutputStreamOptions) (unsafe.Pointer, func(), error) {
	noop := func() {}
	if !options.WASAPISharedAutoConvert {
		return nil, noop, nil
	}

	hostType := C.lk_device_host_api_type(C.PaDeviceIndex(device.Index))
	hostNamePtr := C.lk_device_host_api_name(C.PaDeviceIndex(device.Index))
	if hostNamePtr == nil {
		return nil, noop, fmt.Errorf("portaudio: cannot determine host API for output device %q", device.Name)
	}
	hostName := C.GoString(hostNamePtr)
	if err := validateWASAPISharedAutoConvert(true, hostType == C.paWASAPI, hostName); err != nil {
		return nil, noop, err
	}

	info := C.lk_new_wasapi_auto_convert_stream_info()
	if info == nil {
		return nil, noop, errors.New("portaudio: allocate WASAPI output stream configuration")
	}
	return unsafe.Pointer(info), func() { C.free(unsafe.Pointer(info)) }, nil
}

type wasapiAutoConvertSpec struct {
	Size        uintptr
	HostAPIType int
	Version     uint
	Flags       uint
}

func newWASAPIAutoConvertSpec() (wasapiAutoConvertSpec, error) {
	info := C.lk_new_wasapi_auto_convert_stream_info()
	if info == nil {
		return wasapiAutoConvertSpec{}, errors.New("allocate WASAPI output stream configuration")
	}
	defer C.free(unsafe.Pointer(info))
	return wasapiAutoConvertSpec{
		Size:        uintptr(info.size),
		HostAPIType: int(info.hostApiType),
		Version:     uint(info.version),
		Flags:       uint(info.flags),
	}, nil
}
