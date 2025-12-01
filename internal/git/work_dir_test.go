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
	from    string
}

func createWorkDir(
	t *testing.T,
	commits []*commit,
	branches []*branch,
) (*WorkDir, func()) {
	ctx := t.Context()

	tmp, err := os.MkdirTemp("", "gz-test-work-dir")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// init
	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	wd, err := WorkDirAt(tmp)
	if err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}

	// create main branch commits
	for _, commit := range commits {
		if err := os.WriteFile(filepath.Join(tmp, "data.txt"), []byte(commit.content), 0644); err != nil {
			t.Fatalf("failed to write commit: %v", err)
		}

		if err := wd.gitCommand(ctx, []string{"add", "data.txt"}).Run(); err != nil {
			t.Fatalf("failed to add commit: %v", err)
		}

		if err := wd.gitCommand(
			ctx,
			[]string{"commit", "-m", commit.message},
			WithEnv(
				"GIT_AUTHOR_DATE="+commit.time.Format(time.RFC3339),
				"GIT_COMMITTER_DATE="+commit.time.Format(time.RFC3339),
			)).Run(); err != nil {
			t.Fatalf("failed to commit: %v", err)
		}
	}

	for _, branch := range branches {
		if err := wd.gitCommand(ctx, []string{"checkout", "-b", branch.name, branch.from}).Run(); err != nil {
			t.Fatalf("failed to checkout branch: %v", err)
		}

		for _, commit := range branch.commits {
			if err := os.WriteFile(filepath.Join(tmp, "data.txt"), []byte(commit.content), 0644); err != nil {
				t.Fatalf("failed to write commit: %v", err)
			}

			if err := wd.gitCommand(ctx, []string{"add", "data.txt"}).Run(); err != nil {
				t.Fatalf("failed to add commit: %v", err)
			}

			if err := wd.gitCommand(ctx, []string{"commit", "-m", commit.message}, WithEnv(
				"GIT_AUTHOR_DATE="+commit.time.Format(time.RFC3339),
				"GIT_COMMITTER_DATE="+commit.time.Format(time.RFC3339),
			)).Run(); err != nil {
				t.Fatalf("failed to commit: %v", err)
			}
		}
	}

	return wd, func() {
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}
}

func TestCreateBranch(t *testing.T) {
	ctx := t.Context()

	wd, cleanup := createWorkDir(t, []*commit{
		{
			message: "initial commit",
			content: "initial commit",
			time:    time.Now(),
		},
	}, []*branch{
		{
			name: "feature",
			commits: []*commit{
				{
					message: "feature commit",
					content: "feature commit",
					time:    time.Now(),
				},
			},
			from: "main",
		},
	})
	defer cleanup()

	if err := wd.CreateBranch(ctx, "foo", "main"); err != nil {
		t.Fatal(err)
	}

	head, err := wd.Repository().Head()
	if err != nil {
		t.Fatal(err)
	}

	if head.Name().Short() != "foo" {
		t.Fatal("foo branch not created")
	}

	// expected := "417470bd69d615dc5323f1f97b5c2ba47322c705"
	// if sha := head.Hash().String(); sha != "foo" {
	// 	t.Fatalf("incorrect foo head hash expected: %s, got: %s", expected, sha)
	// }

	fmt.Println(wd)
}
