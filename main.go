package main

import (
	"VMAF-GUI/ui"
)

func main() {
	// Create and run UI
	u := ui.Ui{}
	u.NewUI()
	u.Run()
}
