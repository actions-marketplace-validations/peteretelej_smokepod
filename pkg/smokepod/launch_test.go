package smokepod

import (
	"context"
	"testing"
)

func TestTargetLaunchSpec_ShellArgs(t *testing.T) {
	tests := []struct {
		name    string
		spec    targetLaunchSpec
		command string
		want    []string
	}{
		{
			name:    "no fixed args",
			spec:    targetLaunchSpec{path: "/bin/bash", args: nil},
			command: "echo hello",
			want:    []string{"-c", "echo hello"},
		},
		{
			name:    "with fixed args",
			spec:    targetLaunchSpec{path: "/bin/bash", args: []string{"--norc"}},
			command: "echo hello",
			want:    []string{"--norc", "-c", "echo hello"},
		},
		{
			name:    "multiple fixed args",
			spec:    targetLaunchSpec{path: "/bin/bash", args: []string{"--norc", "--noprofile"}},
			command: "echo hello",
			want:    []string{"--norc", "--noprofile", "-c", "echo hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.shellArgs(tt.command)
			if len(got) != len(tt.want) {
				t.Fatalf("shellArgs() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("shellArgs()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestTargetLaunchSpec_ProcessArgs(t *testing.T) {
	tests := []struct {
		name string
		spec targetLaunchSpec
		want []string
	}{
		{
			name: "no args",
			spec: targetLaunchSpec{path: "/usr/bin/myapp", args: nil},
			want: nil,
		},
		{
			name: "with args",
			spec: targetLaunchSpec{path: "/usr/bin/myapp", args: []string{"--verbose", "--port=8080"}},
			want: []string{"--verbose", "--port=8080"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.processArgs()
			if len(got) != len(tt.want) {
				t.Fatalf("processArgs() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("processArgs()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestTargetLaunchSpec_Cmd(t *testing.T) {
	tests := []struct {
		name      string
		spec      targetLaunchSpec
		extraArgs []string
		wantPath  string
		wantArgs  []string
	}{
		{
			name:      "no extra args",
			spec:      targetLaunchSpec{path: "/bin/bash", args: []string{"--norc"}},
			extraArgs: nil,
			wantPath:  "/bin/bash",
			wantArgs:  []string{"/bin/bash", "--norc"},
		},
		{
			name:      "with extra args",
			spec:      targetLaunchSpec{path: "/bin/bash", args: []string{"--norc"}},
			extraArgs: []string{"-c", "echo hello"},
			wantPath:  "/bin/bash",
			wantArgs:  []string{"/bin/bash", "--norc", "-c", "echo hello"},
		},
		{
			name:      "no fixed or extra args",
			spec:      targetLaunchSpec{path: "/usr/bin/myapp"},
			extraArgs: nil,
			wantPath:  "/usr/bin/myapp",
			wantArgs:  []string{"/usr/bin/myapp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.spec.cmd(context.Background(), tt.extraArgs...)
			if cmd.Path != tt.wantPath {
				// exec.CommandContext resolves the path, so check Args[0] instead
				if cmd.Args[0] != tt.wantPath {
					t.Errorf("cmd.Args[0] = %q, want %q", cmd.Args[0], tt.wantPath)
				}
			}
			if len(cmd.Args) != len(tt.wantArgs) {
				t.Fatalf("cmd.Args = %v, want %v", cmd.Args, tt.wantArgs)
			}
			for i := range cmd.Args {
				if cmd.Args[i] != tt.wantArgs[i] {
					t.Errorf("cmd.Args[%d] = %q, want %q", i, cmd.Args[i], tt.wantArgs[i])
				}
			}
		})
	}
}

func TestTargetLaunchSpec_Cmd_DoesNotMutateArgs(t *testing.T) {
	spec := targetLaunchSpec{path: "/bin/bash", args: []string{"--norc"}}

	// Call cmd with extra args
	_ = spec.cmd(context.Background(), "-c", "echo hello")

	// Verify original args were not mutated
	if len(spec.args) != 1 || spec.args[0] != "--norc" {
		t.Errorf("spec.args mutated: got %v, want [--norc]", spec.args)
	}
}
