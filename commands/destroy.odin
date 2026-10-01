package commands

import "../picker"
import "../tmux"
import "../git"
import "core:fmt"
import "core:flags"

Destroy_Cmd :: struct {}

destroy_run :: proc(args: []string) {
	cmd: Destroy_Cmd
	flags.parse_or_exit(&cmd, args, .Unix)

    worktree, err := picker.picker()
    if err != nil do return

    fmt.println("Killing worktree and session...")
    session_id := git.read_session_id(worktree.path)

    git.destroy_worktree(worktree.path)
    tmux.kill_session(session_id)
}
