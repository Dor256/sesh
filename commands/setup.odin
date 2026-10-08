package commands

import "../exec"
import "core:flags"
import "core:fmt"
import "core:os"
import "core:path/filepath"
import "core:strings"

Setup_Cmd :: struct {}

setup_run :: proc(args: []string) {
	cmd: Setup_Cmd
	flags.parse_or_exit(&cmd, args, .Unix)

	config_dir, ok := user_config_dir()
	if !ok do return

	fmt.println("Setting up Sesh...")
	sesh_dir, dir_err := filepath.join({config_dir, "sesh"}, context.temp_allocator)
	sesh_state_file, file_err := filepath.join({sesh_dir, "state"}, context.temp_allocator)
	if dir_err != nil || file_err != nil do return

	if !os.exists(sesh_dir) {
		os.make_directory_all(sesh_dir)
		_, err := os.create(sesh_state_file)
		if err != nil do return
	}

	tmux_config_path, config_ok := get_tmux_config(config_dir)
	if !config_ok do return

	setup_tmux_hooks(tmux_config_path)
}

setup_tmux_hooks :: proc(config_path: string) {
	client_detached_hook := "set-hook -g client-detached 'run-shell \"sesh save\"'"
	session_closed_hook := "set-hook -g session-closed 'run-shell \"sesh save\"'"
	session_created_hook := "set-hook -g session-created 'run-shell \"sesh save\"'"

	out, err := exec.run_command(
		{
			"grep",
			"-qxF",
			strings.concatenate(
				{client_detached_hook, "\n", session_closed_hook, "\n", session_closed_hook},
				context.temp_allocator,
			),
			config_path,
		},
	)

    // Grep error means hooks exist - exit
	if err == nil do return

	file, file_err := os.open(config_path, {.Create, .Write, .Append})
	if file_err != nil do return
	defer os.close(file)

	string_builder: strings.Builder
	strings.builder_init(&string_builder)
	defer strings.builder_destroy(&string_builder)

    strings.write_string(&string_builder, "\n# Sesh owned tmux hooks for session state recovery\n")
    strings.write_string(&string_builder, client_detached_hook)
    strings.write_string(&string_builder, "\n")
    strings.write_string(&string_builder, session_closed_hook)
    strings.write_string(&string_builder, "\n")
    strings.write_string(&string_builder, session_created_hook)

    _, write_err := os.write_string(file, strings.to_string(string_builder))
    if write_err != nil {
        fmt.println("Failed to append to tmux config file")
    }
}

@(private = "file")
get_tmux_config :: proc(config_dir: string) -> (string, bool) {
	tmux_config_path, file_err := filepath.join(
		{config_dir, "tmux", "tmux.conf"},
		context.temp_allocator,
	)
	if file_err != nil do return "", false
	if os.exists(tmux_config_path) do return tmux_config_path, true

	home, home_err := os.user_home_dir(context.temp_allocator)
	if home_err != nil do return "", false

	legacy_config_path, legacy_err := filepath.join({home, ".tmux.conf"}, context.temp_allocator)
	if legacy_err != nil do return "", false

	if os.exists(legacy_config_path) do return legacy_config_path, true

	tmux_config_dir, config_dir_err := filepath.join({config_dir, "tmux"}, context.temp_allocator)
	if config_dir_err != nil do return "", false
	mkdir_err := os.make_directory_all(tmux_config_dir)
	if mkdir_err != nil do return "", false
	_, create_err := os.create(tmux_config_path)
	if create_err != nil do return "", false

	return tmux_config_path, true
}
