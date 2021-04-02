package util

import "strings"

func ExtractCommand(args []string) string {
	commands := args[1:]
	return strings.Join(commands, " ")
}
