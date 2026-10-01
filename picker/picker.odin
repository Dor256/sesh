package picker

import "core:strings"
import "../exec"
import "../git"

Error :: enum {
    None,
    Dismiss,
    Invalid_Selection,
}

picker :: proc() -> (git.Worktree, Error) {
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

    stdout, exec_err := exec.run_command({"bash", "-c", script})
    output := strings.trim_space(stdout)

    if exec_err == nil do return git.Worktree{path = "", branch = output}, .None

    lines := strings.split_lines(output, context.temp_allocator)

    if len(lines) == 0 do return {}, .Dismiss
    selection_line := lines[len(lines) - 1]
    worktree_data, fields_err := strings.fields(selection_line, context.temp_allocator)

    if fields_err != nil do return {}, .Invalid_Selection
    if len(worktree_data) < 2 do return {}, .Invalid_Selection

    return git.Worktree{path = worktree_data[0], branch = worktree_data[1]}, .None
}

