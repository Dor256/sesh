package tmux

import "core:os"
import "base:runtime"
import "core:testing"
import "../exec"

tmux_tmpdir: string

@(init)
setup :: proc "contextless" () {
    context = runtime.default_context()
    dir, err := os.make_directory_temp("", "odintmux", context.temp_allocator)
    if err != nil do panic("tmux tests: could not isolate socket dir")

    tmux_tmpdir = dir
    os.unset_env("TMUX")
    os.set_env("TMUX_TMPDIR", dir)
}

@(fini)
teardown :: proc "contextless" () {
    context = runtime.default_context()

    exec.run_command({"tmux", "kill-server"})
    err := os.remove_all(tmux_tmpdir)
    if err != nil do panic("Could not perform teardown!")

    free_all(context.temp_allocator)
}

@(test)
test_tmux_session_lifecycle :: proc(t: ^testing.T) {
    expected_session_name := "test-session"
    testing.expect(t, !has_session(expected_session_name))

    session_name := create_session(expected_session_name, "path")

    testing.expect(t, has_session(session_name))

    kill_session(session_name)

    testing.expect(t, !has_session(session_name))
}

