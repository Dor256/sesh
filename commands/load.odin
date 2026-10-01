package commands

import "../git"
import "../tmux"
import "core:flags"
import "core:fmt"
import "core:os"
import "core:path/filepath"
import "core:strings"

Pane :: struct {
	name, path: string,
}

Load_Cmd :: struct {}

load_run :: proc(args: []string) {
	cmd: Load_Cmd
	flags.parse_or_exit(&cmd, args, .Unix)

	is_tmux_active := os.get_env("TMUX", context.temp_allocator) != ""
	if is_tmux_active {
		fmt.println("Can't load tmux state from within tmux!")
		os.exit(1)
	}

	config_dir, ok := user_config_dir()
	if !ok do return

	sesh_state_path, allocator_error := filepath.join(
		{config_dir, "sesh", "state"},
		context.temp_allocator,
	)
	if allocator_error != nil do return
	raw_tmux_state, file_read_error := os.read_entire_file(sesh_state_path, context.temp_allocator)
	if file_read_error != nil do return

	tmux_state := parse_tmux_state(string(raw_tmux_state))

    worktree_path: string
    session_id: string
    for session, panes in tmux_state {
        for pane, idx in panes {
            worktree_path = pane.path
            if idx == 0 {
                session_id = tmux.create_session(session, pane.path)
                tmux.rename_window(session, "1", pane.name)
            } else {
                tmux.new_window(session, pane.name, pane.path)
            }
        }
    }

    git.save_session_id(worktree_path, session_id)
    tmux.attach_default()
}

@(private = "file")
parse_tmux_state :: proc(raw_state: string) -> map[string][dynamic]Pane {
	raw_sessions := strings.split_lines(raw_state, context.temp_allocator)
	state := make(map[string][dynamic]Pane, context.temp_allocator)
	for raw_session in raw_sessions {
		session_tuple := strings.split(raw_session, " ", context.temp_allocator)
		if len(session_tuple) == 0 {
			return nil
		}
		session_name := session_tuple[0]
		pane := Pane {
			name = session_tuple[1],
			path = session_tuple[2],
		}
        if state[session_name] == nil {
            state[session_name] = make([dynamic]Pane, context.temp_allocator)
        }
		append(&state[session_name], pane)
	}
    return state
}
