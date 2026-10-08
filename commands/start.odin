package commands

import "../git"
import "../picker"
import "../tmux"
import "core:flags"
import "core:fmt"
import "core:os"
import "core:path/filepath"
import "core:reflect"
import "core:slice"
import "core:strings"

Start_Cmd :: struct {
	session: string `short:"s"`,
}

start_run :: proc(args: []string) {
	cmd: Start_Cmd
	parsed_arguments := start_parse_args(args)
	flags.parse_or_exit(&cmd, parsed_arguments, .Unix)

	session_name := cmd.session
	cwd, wd_err := os.get_working_directory(context.temp_allocator)
	if wd_err != nil do return

	home, home_err := os.user_home_dir(context.temp_allocator)
	if home_err != nil do return

	worktree_dir := fmt.tprintf("%s/dev/.worktrees", home)
	worktree, picker_err := picker.picker()
	cwd_base := filepath.base(cwd)

	worktree_path := worktree.path
	// New worktree
	if picker_err != nil {
		worktree_path = fmt.tprintf(
			"%s/%s-%s",
			worktree_dir,
			cwd_base,
			strings.to_lower(session_name, context.temp_allocator),
		)
		ok := create_worktree(&worktree, worktree_path)
		if !ok {
			fmt.println("Aborting, branch was not provided")
			os.exit(1)
		}
	}
	session_id: string
	if tmux.has_session(session_name) {
		session_id = tmux.get_session_id(session_name)
	} else {
		fmt.println("Creating session...")
		session_id = tmux.create_session(session_name, worktree_path)

		//Rename the default window
		tmux.rename_window(session_name, "1", "Terminal")
		tmux.new_window(session_name, "Opencode", worktree_path, "opencode")
		tmux.new_window(session_name, "Neovim", worktree_path, "nvim .")
		tmux.new_window(session_name, "Claude", worktree_path, "claude")
	}

    git.save_session_id(worktree_path, session_id)

    is_tmux_active := os.get_env("TMUX", context.temp_allocator) != ""
    if is_tmux_active {
        tmux.switch_session(session_name)
    } else {
        tmux.attach_session(session_name)
    }
}

@(private = "file")
create_worktree :: proc(worktree: ^git.Worktree, worktree_path: string) -> bool {
	if worktree.branch == "" do return false
	git.create_worktree(worktree.branch, worktree_path)
	return true
}

@(private = "file")
start_parse_args :: proc(args: []string) -> []string {
	parsed := slice.clone(args, context.temp_allocator)
	cmd_typeid := typeid_of(Start_Cmd)
	struct_fields := reflect.struct_field_names(cmd_typeid)
	fields := reflect.struct_fields_zipped(cmd_typeid)

	shorthand: string
	for arg, i in args {
		for field in fields {
			shorthand = reflect.struct_tag_get(field.tag, "short")
			if strings.has_prefix(arg, "-") && arg[1:] == shorthand {
				parsed[i] = fmt.tprintf("--%s", field.name)
				break
			}
		}
	}
	return parsed
}
