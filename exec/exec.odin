package exec

import "core:os"

Command_Error :: struct {
	reason: string,
}

Exec_Error :: union {
	os.Error,
	Command_Error,
}

run_command :: proc(
	command: []string,
	allocator := context.temp_allocator,
) -> (
	res: string,
	error: Exec_Error,
) {
	state, stdout, stderr, _ := os.process_exec(os.Process_Desc{command = command}, allocator)
	out := string(stdout)
	err := string(stderr)

	if !state.success do return "", Command_Error{reason = err}
	return out, nil
}

run_and_wait :: proc(
	command: []string,
	allocator := context.temp_allocator,
) -> (
	error: Exec_Error,
) {
	process := os.process_start(
		os.Process_Desc {
			command = command,
			stdin = os.stdin,
			stdout = os.stdout,
			stderr = os.stderr,
		},
	) or_return
	state := os.process_wait(process) or_return
	if !state.success do return Command_Error{reason = "Command failed"}
	return nil
}
