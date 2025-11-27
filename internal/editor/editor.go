package editor

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/kellegous/poop"
)

var fallbackEditors = []string{
	"nvim",
	"vim",
	"vi",
	"nano",
}

type Editor struct {
	cmd string
}

func (e *Editor) Edit(
	ctx context.Context,
	message string,
) (string, error) {
	tmp, err := os.CreateTemp("", "gz-commit-message-*.txt")
	if err != nil {
		return "", poop.Chain(err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if err := os.WriteFile(tmp.Name(), []byte(message), 0644); err != nil {
		return "", poop.Chain(err)
	}

	cmd := exec.CommandContext(ctx, e.cmd, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", poop.Chain(err)
	}

	content, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", poop.Chain(err)
	}

	return string(content), nil
}

func EditFrom(
	ctx context.Context,
	r *git.Repository,
	message string,
) (string, error) {
	e, err := From(r)
	if err != nil {
		return "", poop.Chain(err)
	}
	return e.Edit(ctx, message)
}

func From(r *git.Repository) (*Editor, error) {
	if cmd := fromEnv("GIT_EDITOR", "VISUAL", "EDITOR"); cmd != "" {
		return &Editor{cmd: cmd}, nil
	}

	if cmd := fromGitConfig(r); cmd != "" {
		return &Editor{cmd: cmd}, nil
	}

	if cmd := fromFallbackEditors(); cmd != "" {
		return &Editor{cmd: cmd}, nil
	}

	return nil, errors.New("no editor found")
}

func fromEnv(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return ""
}

func fromGitConfig(r *git.Repository) string {
	cfg, err := r.ConfigScoped(config.SystemScope)
	if err != nil {
		return ""
	}

	if s := cfg.Raw.Section("core"); s != nil {
		return s.Option("editor")
	}

	return ""
}

func fromFallbackEditors() string {
	for _, editor := range fallbackEditors {
		cmd, err := exec.LookPath(editor)
		if err == nil {
			return cmd
		}
	}
	return ""
}
