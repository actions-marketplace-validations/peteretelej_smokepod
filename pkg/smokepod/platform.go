package smokepod

import (
	"context"
	"runtime"
	"strings"
)

func DetectPlatform(ctx context.Context, target Target) PlatformInfo {
	info := PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	shellVersion := detectShellVersion(ctx, target)
	info.ShellVersion = shellVersion

	return info
}

func detectShellVersion(ctx context.Context, target Target) string {
	commands := []string{
		"--version",
		"version",
		"-version",
	}

	for _, cmd := range commands {
		result, err := target.Exec(ctx, cmd)
		if err != nil {
			continue
		}

		output := strings.TrimSpace(result.Stdout)
		if output != "" && result.ExitCode == 0 {
			firstLine := strings.Split(output, "\n")[0]
			if len(firstLine) > 100 {
				firstLine = firstLine[:100]
			}
			return firstLine
		}
	}

	return ""
}
