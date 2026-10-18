# Agent workspace: LED cube

Notes, plans and reference material for:
1. porting the original iCubeSmart firmware to TinyGo
2. fixing KEY5 (PA14)

No project code has been changed so far.

## Status (2026-10-18)

| Item | State |
|---|---|
| KEY5 / PA14 root cause | ✅ Found (TinyGo `Pin.Configure` bit-4 leak). Hardware confirmation pending |
| Original firmware | ✅ Full dump secured. It matches the vendor `1-3D8RGB-Demo-20250128.hex` exactly |
| Firmware architecture | ✅ Generic show engine + compiled data blob made by the PC "3D Software" (script VM, graph codecs) |
| Hardware / pin map | ✅ From schematics + photos |
| Flashing on this PC | ✅ CubeProgrammer CLI over COM4 (BOOT0 jumper) |
| Porting | 🔄 Phase 1/2 in progress: blob format ✅ specified, engine spec drafted from `Simulator.dll`, firmware hardware layer indexed. See [re/](re/README.md) |

## Index

- **Plans**
  - [plan/key5-fix-plan.md](plan/key5-fix-plan.md)
  - [plan/porting-plan.md](plan/porting-plan.md)
- **Investigation**
  - [investigation/key5-pa14.md](investigation/key5-pa14.md): root cause, evidence, workaround, experiments
  - [investigation/original-firmware.md](investigation/original-firmware.md): what's in `firmware/`, what's missing, what the app does
  - [investigation/hardware.md](investigation/hardware.md): your exact boards (from the photos)
  - [investigation/environment.md](investigation/environment.md): toolchain on this PC, build quirks, target memory limits
- **Documentation**
  - [documentation/what-i-need-from-you.md](documentation/what-i-need-from-you.md): **open requests**
  - [documentation/flashing-and-dumping.md](documentation/flashing-and-dumping.md): CubeProgrammer over UART, restore, console
  - [documentation/mcu-reference.md](documentation/mcu-reference.md): GD32/STM32F1 register cheat sheet, display chain timing
  - [documentation/schematics/README.md](documentation/schematics/README.md): pin maps read from the schematics (+ rendered PNGs)
- **Reverse engineering**
  - [re/README.md](re/README.md): how to regenerate decompiles
  - [re/spec-blob-format.md](re/spec-blob-format.md): show-data blob ✅
  - [re/spec-engine.md](re/spec-engine.md): player, renderers, keys, script VM
  - [re/firmware-map.md](re/firmware-map.md): original firmware hardware layer
- **Tools**
  - [tools/blobdump/](tools/blobdump/main.go): blob parser / disassembler (host Go)
  - [tools/ghidra/](tools/ghidra/): headless Ghidra scripts
  - [tools/annotate-firmware.ps1](tools/annotate-firmware.ps1): names peripheral addresses in the firmware decompile
  - [tools/render-pdf.ps1](tools/render-pdf.ps1): PDF to PNG (+ tiles) using the built-in Windows renderer
- **Photos:** `photos/` (gitignored)
- **Vendor PC software:** [external/](external/) (3D Software, Simulator.dll, demo `.gc`/`.hex`)

Rules: the agent doesn't install software, stage or commit. Scripts are `.cmd`, not make.
