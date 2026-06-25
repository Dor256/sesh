package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

type Worktree struct {
	Path   string
	Branch string
}

func Picker() (*Worktree, error) {
	script := `
			git worktree list --porcelain \
				| awk '
					/^worktree / { path = $2 }
					/^branch / { sub("refs/heads/", "", $2); print path, $2 }
				'| fzf --height=40% \
					--reverse \
					--border \
					--with-nth=2 \
					--print-query
	`

	cmd := exec.Command("bash", "-c", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &Worktree{
			Path:   "",
			Branch: strings.TrimSpace(string(output)),
		}, fmt.Errorf("%s", "User dismissed picker")
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no output received from picker")
	}
	selectionLine := lines[len(lines)-1]

	worktreeData := strings.Fields(selectionLine)
	if len(worktreeData) < 2 {
		return nil, fmt.Errorf("invalid picker selection: %s", selectionLine)
	}
	worktree := &Worktree{
		Path:   worktreeData[0],
		Branch: worktreeData[1],
	}

	return worktree, nil
}
