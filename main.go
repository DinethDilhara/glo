/*
Copyright © 2025 Dineth De Alwis <dinethdinlhara66@gmail.com>
*/
package main

import "github.com/DinethDilhara/glo/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
