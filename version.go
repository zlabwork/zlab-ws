package app

import "fmt"

const (
	VersionName = "messenger"

	VersionNumber = "v0.2.0"

	Website = "https://zlab.dev"

	// http://patorjk.com/software/taag/#p=display&h=0&f=Small%20Slant&t=RPC
	banner = `
   __  ___
  /  |/  /  ___  ___ _
 / /|_/ /  (_-< / _ '/
/_/  /_/  /___/ \_, /
               /___/   %s %s

High performance, App framework
Support by %s
%s
____________________________________O/_______
                                    O\
`
)

func Banner(message string) {
	fmt.Printf(banner, VersionName, VersionNumber, Website, message)
}
