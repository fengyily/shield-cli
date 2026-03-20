//go:build (darwin || windows) && withtray

package tray

import _ "embed"

//go:embed icon.png
var iconData []byte

//go:embed icon_22x2.png
var iconSmallData []byte
