package smokepod

import (
	"context"
	"runtime"
)

func DetectPlatform(ctx context.Context, target Target) PlatformInfo {
	info := PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	if localTarget, ok := target.(*LocalTarget); ok {
		info.ShellVersion = localTarget.GetVersion(ctx)
	}

	return info
}
