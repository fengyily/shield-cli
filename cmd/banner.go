package cmd

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

const banner = `
   _____ __    _       __    __   ________    ____
  / ___// /_  (_)__   / /___/ /  / ____/ /   /  _/
  \__ \/ __ \/ // _ \/ // __  / / /   / /    / /
 ___/ / / / / //  __/ // /_/ / / /___/ /____/ /
/____/_/ /_/_/ \___/_/ \__,_/  \____/_____/___/
`

func PrintBanner() {
	fmt.Print("\033[36m") // Cyan
	fmt.Print(banner)
	fmt.Print("\033[0m")

	fmt.Printf("\033[1;33m  Shield CLI\033[0m - \033[90mSecure Tunnel Connector\033[0m\n")
	fmt.Println()

	fmt.Printf("  \033[90m├─\033[0m Version:    \033[32m%s\033[0m\n", Version)
	fmt.Printf("  \033[90m├─\033[0m Go:         \033[90m%s\033[0m\n", runtime.Version())
	fmt.Printf("  \033[90m└─\033[0m Platform:   \033[90m%s/%s\033[0m\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	fmt.Println("  \033[90m──────────────────────────────────────────────────\033[0m")
	fmt.Println()
}
