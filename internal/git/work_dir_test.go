package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type commit struct {
	message string
	content string
	time    time.Time
}

type branch struct {
	name    string
	commits []*commit
}

func createWorkDir(
	t *testing.T,
	commits []*commit,
	branches []*branch,
) (*WorkDir, func()) {
	tmp, err := os.MkdirTemp("", "gz-test-work-dir")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// init
	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// create main branch commits
	for _, commit := range commits {
		if err := os.WriteFile(filepath.Join(tmp, "data.txt"), []byte(commit.content), 0644); err != nil {
			t.Fatalf("failed to write commit: %v", err)
		}

		cmd = exec.Command("git", "add", "data.txt")
		cmd.Dir = tmp
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to add commit: %v", err)
		}

		cmd = exec.Command("git", "commit", "-m", commit.message)
		cmd.Env = append(
			os.Environ(),
			"GIT_AUTHOR_DATE="+commit.time.Format(time.RFC3339),
			"GIT_COMMITTER_DATE="+commit.time.Format(time.RFC3339),
		)
		cmd.Dir = tmp
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to commit: %v", err)
		}
	}

	wt, err := WorkDirAt(tmp)
	if err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}

	return wt, func() {
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}
}

func TestCreateBranch(t *testing.T) {
	wd, cleanup := createWorkDir(t, []*commit{
		{
			message: "initial commit",
			content: "initial commit",
			time:    time.Now(),
		},
	}, nil)
	defer cleanup()

	fmt.Println(wd)
}
