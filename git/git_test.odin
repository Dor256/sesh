package git

import "core:strings"
import "core:os"
import "../exec"
import "core:testing"

setup :: proc(t: ^testing.T) -> string {
    tmp_path, err := os.make_directory_temp("", "odintest", allocator = context.temp_allocator)
    if err != nil do testing.fail(t)
    dir := strings.trim_space(tmp_path)

    exec.run_command({"git", "init", "-b", "main", dir})
    exec.run_command({"git", "-C", dir, "-c", "user.name=test", "-c", "user.email=test@test", "commit", "--allow-empty", "-m", "init"})
    return dir
}

teardown :: proc(git_repo: string) {
    err := os.remove_all(git_repo)
    if err != nil do panic("Could not perform teardown!")
    free_all(context.temp_allocator)
}

@(test)
test_create_destroy_worktree :: proc(t: ^testing.T) {
    repo := setup(t)
    defer teardown(repo)

    expected := "feature-test"

    wt := strings.concatenate({repo, "-wt"}, context.temp_allocator)

    create_worktree(expected, wt, repo)

    output, err := exec.run_command({"git", "-C", wt, "branch", "--show-current"})
    if err != nil do testing.expectf(t, false, "Error was thrown %v", err)
    testing.expect_value(t, output, expected)

    destroy_worktree(wt)
    worktree_list, list_err := exec.run_command({"git", "-C", repo, "worktree", "list"})
    if list_err != nil do testing.expectf(t, false, "Error was thrown %v", list_err)

    testing.expect(t, !strings.contains(worktree_list, wt))
}

@(test)
test_session_id_read_write :: proc(t: ^testing.T) {
    repo := setup(t)
    defer teardown(repo)

    expected := "sesh-id"

    session_id := read_session_id(repo)

    testing.expect(t, expected != session_id, "Session ID was prepopulated, testing leak!")

    save_session_id(repo, expected)

    session_id = read_session_id(repo)

    testing.expect_value(t, session_id, expected)
}

