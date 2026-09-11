"""The command loop. Reads a command word, dispatches it, and repeats until the user quits."""

import commands
import store


def print_menu() -> None:
    """Print the application title and list of available commands."""
    print("Expense Tracker")
    print("Commands: add, list, filter, summary, edit, delete, quit")


def main() -> None:
    """Run the main command loop."""
    store.seed_data()
    print_menu()

    while True:
        try:
            line = input("> ")
        except (EOFError, KeyboardInterrupt):
            print("Goodbye.")
            break

        cmd = line.strip().lower()
        if not cmd:
            continue

        if cmd == "add":
            commands.add()
        elif cmd == "list":
            commands.list_expenses()
        elif cmd == "filter":
            commands.filter_expenses()
        elif cmd == "summary":
            commands.summary()
        elif cmd == "edit":
            commands.edit()
        elif cmd == "delete":
            commands.delete()
        elif cmd == "quit":
            print("Goodbye.")
            break
        else:
            print(f"Unknown command: {line.strip()}.")
            print_menu()


if __name__ == "__main__":
    main()
