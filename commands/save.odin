package commands

import "../tmux"
import "core:path/filepath"
import "core:os"
import "core:flags"

Save_Cmd :: struct {}

save_run :: proc(args: []string) {
	cmd: Save_Cmd
	flags.parse_or_exit(&cmd, args, .Unix)
    config_dir, ok := user_config_dir()
    if !ok do return

    sesh_state_path, join_err := filepath.join({config_dir, "sesh", "state"}, context.temp_allocator)
    if join_err != nil do return
    tmux_state := tmux.list_panes()
    save_err := os.write_entire_file_from_string(sesh_state_path, tmux_state)
    if save_err != nil do return
}

