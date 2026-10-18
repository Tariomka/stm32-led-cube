# Reverse engineering workspace

Specs (committed):
- [spec-blob-format.md](spec-blob-format.md): show-data blob format ✅ verified
- [spec-engine.md](spec-engine.md): player, renderers, keys, script VM (from `Simulator.dll`)
- [firmware-map.md](firmware-map.md): hardware layer index of the original firmware

Generated, vendor-derived outputs live in `.agent/external/` (gitignored):

| Path | What |
|---|---|
| `.agent/external/decompiled/` | ILSpy C# of `3D Software - RGB.exe` and `Download Software.exe` (you made it) |
| `.agent/external/re/Simulator.c` | Ghidra decompile of `Simulator.dll` |
| `.agent/external/re/Firmware.c` | Ghidra decompile of `firmware/FullDump_512K.bin` |
| `.agent/external/re/Firmware.annotated.c` | Same, with literal-pool peripheral addresses named |

## Regenerate

```powershell
$env:JAVA_HOME = scoop prefix temurin21-jdk
$g = scoop prefix ghidra
$out = "$PWD\.agent\external\re"; New-Item -ItemType Directory -Force $out, "$env:TEMP\ghidra_proj" | Out-Null

# Simulator.dll (x86 PE, auto-detected)
& "$g\support\analyzeHeadless.bat" "$env:TEMP\ghidra_proj" sim -import ".agent\external\3D Software - RGB Cube\Simulator.dll" `
  -overwrite -scriptPath ".agent\tools\ghidra" -postScript ExportDecompiled.java "$out\Simulator.c"

# Firmware (raw ARM Cortex-M image)
& "$g\support\analyzeHeadless.bat" "$env:TEMP\ghidra_proj" fw -import "firmware\FullDump_512K.bin" -overwrite `
  -processor "ARM:LE:32:Cortex" -loader BinaryLoader -loader-baseAddr 0x08000000 `
  -scriptPath ".agent\tools\ghidra" -preScript Gd32Setup.java -postScript ExportDecompiled.java "$out\Firmware.c"

# Name peripheral addresses hidden in literal pools
.\.agent\tools\annotate-firmware.ps1
```

Each run takes about 1 minute. The Ghidra projects stay in `%TEMP%\ghidra_proj` (`sim`, `fw`) and can be opened in the GUI for interactive work.

Blob inspection: [../tools/blobdump](../tools/blobdump/main.go) (`cd .agent\tools\blobdump; go run . <dump.bin|vendor.hex> [-v]`).
