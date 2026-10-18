# Carry-forward prompt (next session)

Paste this to the next agent, or read it first if you are that agent.

---

You are continuing work on **stm32-led-cube**: TinyGo firmware for an iCubeSmart 3D8RGB 8x8x8 RGB LED cube (GD32F103RET6 "yellow board": 512 KB flash, 64 KB RAM). The repo is `d:\Git\some_personal_shit\Led_Cube\stm32-led-cube`. The goal was to **port all of the original vendor firmware** (every premade show plus all hardware features) into this project, keep the Raspberry Pi Pico link and the desktop app working, and fix KEY5 (PA14). The port is **code-complete**. What remains is hardware verification and the later "serial/Pico" stage.

## Hard rules from the user

- **Never install software.** **Never `git add`/commit/push.** The user reviews and commits.
- Use `.cmd` scripts in `scripts\`, not make (the `Makefile` is kept for Linux only).
- **Never write to `%TEMP%`.** One-off scripts and files go in `.agent/external/temp` (gitignored).
- `.agent/` is the agent's memory and docs folder. `.agent/external` and `.agent/photos` are gitignored.
- Don't run ILSpy. The C# decompile is already in `.agent/external/decompiled`. Ghidra is allowed.
- Flashing: STM32CubeProgrammer CLI over **COM4** with the BOOT0 jumper (`scripts\flash.cmd`). ST-Link and the logic analyzer are not used.
- Don't edit files via terminal commands (no `gofmt -w`, `sed`, etc.). Use the edit tools.
- The user may edit files between sessions. Always re-read a file before editing it.

## Environment quirks

- In agent terminals, strip TinyGo's goroot from PATH before `go`/`tinygo`, otherwise you get "sync/atomic is not in std":
  `$env:Path = ($env:Path -split ';' | Where-Object { $_ -notmatch 'tinygo\\goroot-' }) -join ';'`
- Build from the **repo root**: `cmd /c scripts\build.cmd [pkg]`. The default pkg is `./cmd/8x8_rgb_cube` and the output is `bin\main.hex`.
  - It uses the custom target `targets/gd32f103re.json` with linker script `targets/gd32f103re.ld`, which gives FLASH 508K.
  - The last 4 KB of flash are reserved for the console registry.
- Other scripts:
  - `scripts\flash.cmd [COMx] [file]`: builds first unless a file is given.
  - `scripts\monitor.cmd`: COM4, 115200.
  - `scripts\dump.cmd`
  - `scripts\test-oracle.cmd`: Windows/386 oracle tests against the vendor `Simulator.dll`.
- Host tests (the packages that don't import `machine`):
  `go test ./internal/engine/... ./internal/remote/... ./internal/console/... ./internal/registry/... ./internal/settings/... ./internal/cascade/...`
- **Known false alarm:** `gofmt -l` flags every file because of CRLF line endings (git autocrlf). To check real formatting, normalize CRLF→LF into `.agent/external/temp/fmt.go` and run `gofmt -d` on that copy.
- `create_file` cannot overwrite a file. Write `x.go.new`, then `Move-Item -Force`.
- **Check new Go files.** A tool sometimes inserted a duplicate `package` line at the top of a new file. Check the first line after creating one.
- **Avoid `strings.EqualFold` / `strings.Fields` in firmware code.** They pull Unicode tables into RAM (+6 KB).
- Flash budget: the last build was about 120 KB flash and 7.9 KB static RAM.

## Toolchain

TinyGo 0.43.0 (scoop), Go 1.27.2, STM32CubeProgrammer v2.23.0, Ghidra + temurin21-jdk (scoop).

## Architecture

| Path | What |
|---|---|
| `cmd/8x8_rgb_cube` | Main firmware |
| `cmd/gpio_diag` | KEY5/GPIO diagnostic program |
| `internal/engine/blob` | Parser for the original show data, plus graph codecs |
| `internal/engine/render` | Canvas, graphs with rotation, ball, title |
| `internal/engine/player` | Cartoon player, key actions, script VM. Has `SetSpeed`, `Select`, `AudioFrame(mode)` |
| `internal/engine/frame` | Voxels → MBI5024 bit planes |
| `internal/engine/audio` | Microphone → 0..63 level, using the original algorithm |
| `internal/engine/oracle` | Frame-exact tests against `Simulator.dll` (windows && 386) |
| `internal/show/original` | Embedded original show data (`original.bin`, 32100 B) |
| `internal/remote` | NEC IR decoder and the original button → action table |
| `internal/console` | Vendor console commands; `Demux` splits UART1 between `#?cmd?!` console lines and Pico frames (`0xAA…0x55` plus the ASCII handshake) |
| `internal/registry` | Key/value registry in the vendor flash format |
| `internal/settings` | 8-byte EEPROM record at 0x00: speed, play mode, group, cartoon |
| `internal/cascade` | Original multi-cube protocol: frames `A5 55 len payload xor`; Hello / Number / Sync (show clock + key mask); `KeyMask` |
| `internal/component` | Hardware drivers (table below) |
| `internal/controller/board.go`, `console.go` | `YellowBoard` implements the `Board` interface and is the console's `Device` |
| `internal/runner/runner.go`, `input.go` | Main loop (below) |
| `tools/fwtool` | Extracts the show data from a vendor hex or flash dump |

`internal/component` drivers:

| File / type | Hardware |
|---|---|
| `pin.go` | `NewInputPin` (pull-down config + `High()`: the KEY5 workaround); `Button` (EXTI falling edge, 30 ms debounce, priority lowered to 0xE0); `DisableDebugPort` |
| `scanner.go` | TIM2 display refresh: 2 brightness planes, about 104 Hz |
| `microphone.go` | ADC1 continuous, sampled by a TIM3 1 kHz interrupt |
| `infrared.go` | PC6 EXTI on both edges, timed with `time.Now` |
| `serial.go` | UART1 split between console and Pico |
| `eeprom.go` | Bit-banged I2C on PB6/PB7 to the 24C02 |
| `flash.go` | Page erase and program for the registry at `0x0807F000` (+ backup `0x0807F800`) |
| `cascade.go` | UART2 on PA2/PA3, 115200 |

Runner loop:

- **Boot.** `HandleStartup` prints `T='type'` and `Config1/Config2`. Then cascade discovery runs (about 1.1 s, prints `DeviceType/DeviceNum`). `placeDisplay` positions the canvas from the cube number.
- **Onboard mode.** `Player.Step(clock)` runs, then the loop waits on the shared show clock.
- **Input (`handleInput`).** Every tick it:
  - polls the console and serial port,
  - saves settings 2 s after the last key action,
  - polls the cascade link,
  - turns key presses into a `KeyMask`; the master sends Sync every 8 ms and runs actions on rising edges, slaves run actions from the master's mask,
  - handles IR on the master only.
- **KEY6.** A press of 1 s or longer cycles the runner mode (Onboard / Serial / Debug). A short press runs the show key function.
- **Serial and Debug modes are placeholders.**

## Status

### Confirmed on hardware

- KEY5 fix (`gpio_diag`)
- Display, orientation and colors
- All 4 DIP groups
- All keys

### Implemented but NOT yet checked on hardware

The checklist is in `.agent/documentation/what-i-need-from-you.md`:

- Audio group reacts to sound.
- IR remote buttons.
- Keys now act through the cascade key mask (up to 10 ms later); check short taps.
- No freezes after the EXTI priority change.
- Boot log lines.
- Console: `#?Help?!`, SetReg / power-cycle / GetReg, ReadMem, ReBoot.
- EEPROM settings survive a power cycle.
- Cascade (needs a second cube).

### Not ported on purpose

- The vendor IAP bootloader and PC update tool.
- The serial-number copy protection.
- The IWDG watchdog (about 4 s in the original).
- Display `dir = 1`.
- Vendor GBK help texts (ours are English).
- The original console runs at 921600 baud; ours stays at 115200 because the Pico link shares UART1.

## Next work (in order)

1. **Act on the user's hardware report** from the checklist above and fix whatever fails. Likely suspects:
   - EEPROM bit-bang timing or missing pull-ups
   - flash page write on the GD32 (stalls the display briefly; check that the registry survives)
   - microphone level scale (the thresholds assume a 12-bit ADC idling at about 1050)
   - IR timing jitter
2. **Serial (Pico) mode** (`runner.runSerialLoop` is a placeholder). The sibling projects are in `..\` (one directory up from the repo):
   - `led-common-lib`, v0.4.3: `network.UartProcessor` (frames `AA type len16BE payload crc32BE 55`, handshake `"SGFuZHNoYWtl"`, types Message/Bytes/Setting/Ping/Pong) and `led.LedLayout` (8 layers × 24 bytes: G 0-7, B 8-15, R 16-23).
   - `rpi-led-communication`: Pico W, a TCP server that will forward packets to this cube over UART. Forwarding is still a TODO there.
   - `desktop-led-controller`: Raylib GUI that sends `network.Packet` over TCP. It is pinned to the older v0.1.0 of the lib.

   Design decisions to make with the user: live frames via `Board.ShowLayout`, uploading show data, selecting shows, settings packets.

   Note: `UartProcessor.WriteMessage` performs a blocking handshake (up to 5 s). `global.HandlePrint` uses it.
3. **Optional extras:**
   - IWDG watchdog (feed it in every loop; keep the handshake time in mind)
   - display `dir = 1`
   - tuning the brightness planes (the original uses 3 planes weighted 1:3:6 out of 10)
4. **Upstream TinyGo issue.** The user has an open issue about `Pin.Configure`'s bit-4 leak. The fix is `uint32(config.Mode)&0xf` in `machine_stm32f103.go`. A PR or issue update can be drafted on request.
5. **Housekeeping.** `tmp/cmd/test` and `tmp/cmd/deez_nuts` (gitignored) import the deleted `internal/component/registers` package and no longer build.

## Where the knowledge lives

- `.agent/README.md`: index and status table.
- `.agent/plan/porting-plan.md`: phase checklist with done/remaining items.
- `.agent/plan/key5-fix-plan.md`
- `.agent/re/`:
  - `spec-blob-format.md`
  - `spec-engine.md`: player, renderers, VM, audio level
  - `firmware-map.md`: interrupts, IR, console wire format and registry flash format, EEPROM records, cascade protocol, next reverse-engineering steps
  - `README.md`: how to regenerate the decompiles
- `.agent/investigation/`:
  - `key5-pa14.md`: root cause and hardware results
  - `hardware.md`: pin map; SW3/SW4 are unreachable and read low
  - `original-firmware.md`
  - `environment.md`
- `.agent/documentation/`:
  - `what-i-need-from-you.md`: hardware checklist
  - `flashing-and-dumping.md`
  - `mcu-reference.md`
  - `schematics/`
- **Reverse-engineering artifacts** (gitignored; grep them with `includeIgnoredFiles`):
  - `.agent/external/re/Firmware.annotated.c`: Ghidra decompile of `firmware/FullDump_512K.bin`, with peripheral and `RAM_2000xxxx` labels. The `printf` function is `FUN_08036b98`.
  - `Firmware.strings.txt`
  - `Simulator.c`
  - The Ghidra project is in `tmp\ghidra_proj` (Ghidra rejects dot-prefixed paths). Regenerate with the commands in `.agent/re/README.md`, then run `.agent/tools/annotate-firmware.ps1`.
  - `Gd32Setup.java` also seeds functions from code pointers found in literal pools.
- **Vendor material:**
  - `.agent/external/3D Software - RGB Cube/`: PC software, `Simulator.dll`, demo hex files.
  - `.agent/external/decompiled/`: C# of the PC tools; `ClassCommunication.cs` holds the console wire format.
- **Original firmware dump:** `firmware/FullDump_512K.bin`. Restore instructions are in `flashing-and-dumping.md`.

## How the user works

- Brief, precise answers.
- Do the work instead of only proposing it, but ask before big design decisions (the user already chose: full console, cascade, EEPROM persistence).
- Summarize what needs a hardware check.
- Always end by reminding the user that nothing was staged or committed.
