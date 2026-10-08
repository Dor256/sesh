package sesh

import cmds "./commands"
import "core:fmt"
import "core:mem"
import "core:os"

main :: proc() {
	when ODIN_DEBUG {
		track: mem.Tracking_Allocator
		mem.tracking_allocator_init(&track, context.allocator)
		context.allocator = mem.tracking_allocator(&track)

		defer {
			if len(track.allocation_map) > 0 {
				fmt.eprintf("=== %v allocations not freed: ===\n", len(track.allocation_map))
				for _, entry in track.allocation_map {
					fmt.eprintf("- %v bytes @ %v\n", entry.size, entry.location)
				}
			}
			mem.tracking_allocator_destroy(&track)
		}
	}

	if len(os.args) < 2 {
		fmt.println("Usage: sesh <command>")
		os.exit(1)
	}

	cmd := os.args[1]
	arguments := os.args[1:]

	switch cmd {
	case "start":
		cmds.start_run(arguments)
		fmt.println("Starting")
	case "destroy":
		cmds.destroy_run(arguments)
		fmt.println("Destroying")
	case "setup":
		cmds.setup_run(arguments)
		fmt.println("Setup")
	case "save":
		cmds.save_run(arguments)
		fmt.println("Saving")
	case "load":
		cmds.load_run(arguments)
		fmt.println("Loading")
	case:
		fmt.printfln("Unknown command: %s", cmd)
		os.exit(1)
	}
	free_all(context.temp_allocator)
}
