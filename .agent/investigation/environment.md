# Development environment (this Windows PC)

Checked 2026-10-18.

| Tool | Status | Notes |
|---|---|---|
| TinyGo | ✅ 0.43.0 (scoop) | LLVM 22 |
| Go | ✅ 1.27.2 (`C:\Program Files\Go`) | |
| OpenOCD | ✅ 0.12.0 (scoop) | Only useful with an SWD probe (ST-Link) |
| 7zip | ✅ (scoop) | |
| GNU make | ❌ | `Makefile` is unusable as is. It also uses `cmd /C date`, `if [ ! -d ]` (POSIX) and `-size full` |
| STM32CubeProgrammer | ❌ | Needed for the UART bootloader path, or use `stm32flash` instead |
| stm32flash | ❌ | Open-source alternative for the UART bootloader |
| USB-UART driver | ? | Depends on the adapter (CH340 / CP210x / FTDI). The only serial port now is Intel AMT SOL `COM3` |
| Ghidra + JDK 21 | ❌ | Needed to reverse-engineer the vendor app |
| Python | ❌ | Not needed. Helper tools will be written in Go or PowerShell |

## Build

```powershell
tinygo build -o ./bin/main.hex -target bluepill-clone -opt s -panic print -gc conservative -size full ./cmd/8x8_rgb_cube/
```

Current result: flash 50 580 B, RAM 3 796 B.

### Agent-terminal-only build failure

In terminals spawned by the agent, `tinygo build` fails with `package sync/atomic is not in std (...tinygo\goroot-<hash>\src\sync\atomic)`.

Cause:
- `.vscode/settings.json` sets `go.toolsEnvVars.GOROOT` to a TinyGo-generated cached GOROOT (`%LOCALAPPDATA%\tinygo\goroot-<hash>`).
- The Go extension prepends `<GOROOT>\bin` to `PATH` for the terminals it injects into.
- TinyGo then uses that `go.exe` and builds a GOROOT-of-a-GOROOT, which is missing packages.

Your own terminal doesn't get the injection, so it works.

- Agent workaround: strip `tinygo\goroot-` entries from `PATH` before calling `tinygo`. Verified to build.
- Optional cleanup: add `"go.terminal.activateEnvironment": false` to `.vscode/settings.json`.
- The hashed GOROOT path is machine- and version-specific. Re-run *TinyGo: Select target → bluepill-clone* after each TinyGo upgrade.

## Target definition problem

`bluepill-clone` inherits `bluepill` → `targets/stm32.ld`:

```text
FLASH_TEXT : ORIGIN = 0x08000000, LENGTH = 64K
RAM        : ORIGIN = 0x20000000, LENGTH = 20K
```

The GD32F103**RE**T6 has **512 KB flash / 64 KB SRAM**. The port will need far more than 64 KB (animation data alone is tens of KB), so a custom target is required:

```text
targets/gd32f103re.json  -> inherits bluepill-clone, linkerscript with 512K / 64K
(optional) origin 0x08006000 variant, to keep the vendor IAP bootloader in place
```

Clock: TinyGo runs the PLL at 72 MHz from the 8 MHz HSE. That matches the board crystal (Y2 = 8.000 MHz). TinyGo uses **TIM4** for its tick timer, so don't use TIM4 for the display.
