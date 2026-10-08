package commands

import "core:os"
import "core:path/filepath"

user_config_dir :: proc() -> (string, bool) {
    xdg := os.get_env("XDG_CONFIG_HOME", context.temp_allocator)
    if xdg != "" do return xdg, true

    home, err := os.user_home_dir(context.temp_allocator)
    if err != nil do return "", false

    if ODIN_OS == .Linux || ODIN_OS == .Darwin {
        joined, err := filepath.join({home, ".config"}, context.temp_allocator)
        if err != nil do return "", false
        return joined, true
    }
    return "", false
}
