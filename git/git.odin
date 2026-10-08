package git

import "../exec"

Worktree :: struct {
	path:   string,
	branch: string,
}

create_worktree :: proc(branch_name, worktree_path: string, maybe_src_dir: Maybe(string) = nil) {
	dir, ok := maybe_src_dir.?
	if ok {
		exec.run_command({"git", "-C", dir, "worktree", "add", "-b", branch_name, worktree_path, "main"})
	} else {
		exec.run_command({"git", "worktree", "add", "-b", branch_name, worktree_path, "main"})
	}
}

destroy_worktree :: proc(path: string) {
	exec.run_command({"git", "-C", path, "worktree", "remove", "--force", "."})
}


save_session_id :: proc(path, tmux_session_id: string) {
	if err := ensure_worktree_config(path); err != nil do return
	exec.run_command(
		{"git", "-C", path, "config", "--worktree", "custom.worktree.tmuxid", tmux_session_id},
	)
}

read_session_id :: proc(path: string) -> string {
	out, _ := exec.run_command(
		{"git", "-C", path, "config", "--worktree", "custom.worktree.tmuxid"},
	)
	return out
}


@(private = "file")
ensure_worktree_config :: proc(path: string) -> exec.Exec_Error {
	_, err := exec.run_command({"git", "-C", path, "config", "extensions.worktreeConfig", "true"})
	return err
}
