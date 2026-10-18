// Package original holds the show data of the original iCubeSmart 3D8RGB firmware
// (extracted with tools/fwtool from 1-3D8RGB-Demo-20250128.hex / the flash dump).
package original

import _ "embed"

//go:embed original.bin
var Blob string
