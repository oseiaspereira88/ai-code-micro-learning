package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunNoArgsPrintsUsageToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), nil, &stdout, &stderr)

	if code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should stay empty on usage error, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("stderr missing usage text: %q", stderr.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"bogus"}, &stdout, &stderr)

	if code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr.String(), `unknown command "bogus"`) {
		t.Fatalf("stderr missing error text: %q", stderr.String())
	}
}

func TestRunVersionIsDeterministic(t *testing.T) {
	var out1, out2, stderr bytes.Buffer
	Run(context.Background(), []string{"version"}, &out1, &stderr)
	Run(context.Background(), []string{"version"}, &out2, &stderr)

	if out1.String() != out2.String() {
		t.Fatalf("version output not deterministic: %q vs %q", out1.String(), out2.String())
	}
	if out1.Len() == 0 {
		t.Fatal("version produced no output")
	}
}

func TestRunDoctorReportsOK(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"doctor"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d", code, exitOK)
	}
	if !strings.Contains(stdout.String(), "status: ok") {
		t.Fatalf("doctor output missing status: %q", stdout.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d", code, exitOK)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("help output missing usage: %q", stdout.String())
	}
}

func TestRunRespectsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stdout, stderr bytes.Buffer
	code := Run(ctx, []string{"doctor"}, &stdout, &stderr)

	if code != exitError {
		t.Fatalf("exit code = %d, want %d", code, exitError)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should stay empty when cancelled, got %q", stdout.String())
	}
}
