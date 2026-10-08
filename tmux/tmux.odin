package tmux

import "../exec"
import "core:fmt"

attach_default :: proc() {
	exec.run_and_wait({"tmux", "attach"})
}

attach_session :: proc(session_name: string) {
	exec.run_command({"tmux", "attach", "-t", session_name})
}

switch_session :: proc(session_name: string) {
	exec.run_command({"tmux", "switch", "-t", session_name})
}

create_session :: proc(session_name, path: string) -> string {
	stdout, _ := exec.run_command(
		{"tmux", "new-session", "-d", "-P", "-F", "#{session_id}", "-s", session_name, "-c", path},
	)
	return stdout
}

new_window :: proc(session_name, window_name, path: string, tools: ..string) {
	if len(tools) == 0 {
		// Plain terminal
		exec.run_command({"tmux", "new-window", "-t", session_name, "-c", path, "-n", window_name})
        return
	}
	tool_cmd := fmt.tprintf("%s; zsh", tools[0])
	exec.run_command(
		{"tmux", "new-window", "-t", session_name, "-c", path, "-n", window_name, tool_cmd},
	)
}

rename_window :: proc(session_name, old_name, new_name: string) {
	absolute_name := fmt.tprintf("%s:%s", session_name, old_name)
	exec.run_command({"tmux", "rename-window", "-t", absolute_name, new_name})
}

kill_session :: proc(session_name_or_id: string) {
	exec.run_command({"tmux", "kill-session", "-t", session_name_or_id})
}

has_session :: proc(session_name: string) -> bool {
	_, err := exec.run_command({"tmux", "has-session", "-t", session_name})
	return err == nil
}

get_session_id :: proc(session_name: string) -> string {
	stdout, _ := exec.run_command(
		{"tmux", "display-message", "-t", session_name, "-p", "#{session_id}"},
	)
	return stdout
}

list_panes :: proc() -> string {
	stdout, _ := exec.run_command(
		{"tmux", "list-panes", "-a", "-F", "#{session_name} #{window_name} #{pane_current_path}"},
	)
	return stdout
}
