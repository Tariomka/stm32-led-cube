# Carry-forward prompt (next session)

Paste this to the next agent, or read it first if you are that agent.

---

You are continuing work on **stm32-led-cube**: TinyGo firmware for an iCubeSmart 3D8RGB 8x8x8 RGB LED cube (GD32F103RET6 "yellow board": 512 KB flash, 64 KB RAM). The repo is `d:\Git\some_personal_shit\Led_Cube\stm32-led-cube`. The goal was to **port all of the original vendor firmware** (every premade show plus all hardware features) into this project and fix KEY5 (PA14). The port is **code-complete**. The Pico and the desktop app now talk to the cube with the same console + upload protocol as the PC tools (no separate Pico protocol). What remains is hardware verification and the sibling projects' alignment (instructions in each sibling's `.agent/stm32-alignment.md`).

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
  - It uses the custom target `targets/gd32f103re.json` with linker script `targets/gd32f103re.ld`, which gives FLASH 350K.
  - Above that: uploaded show slot at `0x08057800` (158 KB), console registry at `0x0807F000` (+ backup `0x0807F800`).
- Other scripts:
  - `scripts\flash.cmd [COMx] [file]`: builds first unless a file is given.
  - `scripts\monitor.cmd`: COM4, 115200.
  - `scripts\dump.cmd`
  - `scripts\test-oracle.cmd`: Windows/386 oracle tests against the vendor `Simulator.dll`.
  - `scripts\terminal.cmd [COMx] [baud] [command]`: console terminal (wraps lines in `#?…?!`).
  - `scripts\upload.cmd [COMx]`: starts the vendor Download Software at 115200 for show uploads.
- Host tests (the packages that don't import `machine`):
  `go test ./internal/engine/... ./internal/remote/... ./internal/console/... ./internal/registry/... ./internal/settings/... ./internal/cascade/... ./internal/show/... ./internal/upload/...`
- Regenerate the Go shows: `go run ./tools/showgen` (overwrites `internal/show/{dynamic,audio,test,library}.go`; `show.go` is hand-written).
- **Known false alarm:** `gofmt -l` flags every file because of CRLF line endings (git autocrlf). To check real formatting, normalize CRLF→LF into `.agent/external/temp/fmt.go` and run `gofmt -d` on that copy.
- `create_file` cannot overwrite a file. Write `x.go.new`, then `Move-Item -Force`.
- **Check new Go files.** A tool sometimes inserted a duplicate `package` line at the top of a new file. Check the first line after creating one.
- **Avoid `strings.EqualFold` / `strings.Fields` in firmware code.** They pull Unicode tables into RAM (+6 KB).
- Flash budget: the last build was about 160 KB flash and 10.5 KB static RAM.

## Toolchain

TinyGo 0.43.0 (scoop), Go 1.27.2, STM32CubeProgrammer v2.23.0, Ghidra + temurin21-jdk (scoop).

## Architecture

| Path | What |
|---|---|
| `cmd/8x8_rgb_cube` | Main firmware |
| `cmd/gpio_diag` | KEY5/GPIO diagnostic program |
| `internal/engine/blob` | Parser for the original show data, plus graph codecs |
| `internal/engine/render` | Canvas, graphs with rotation, ball, title |
| `internal/engine/player` | Cartoon player, key actions, script VM. Plays a `player.Show` (list of `player.Cartoon`): `keyframe.Show` (Go shows) or `player.NewBlobShow(blob, groups...)` (show data). Has `SetShow` (keeps speed/play model/pause/display off), `SetSpeed`, `Select`, `Count`, `SpeedText`, `AudioFrame(mode)` |
| `internal/engine/keyframe` | Go show model: `Cartoon{Frames, Loop, Delay, Elements}`, `Element` = `[]Prop` (`Graph` with From/To `Pose`, `Ball`), `Library` of `Shape` (ASCII art, one line per layer from the top, z rows separated by spaces) and `Palette`. Reproduces the vendor per-element frame quirk |
| `internal/engine/frame` | Voxels → MBI5024 bit planes |
| `internal/engine/audio` | Microphone → 0..63 level, using the original algorithm |
| `internal/engine/oracle` | Frame-exact tests against `Simulator.dll` (windows && 386) |
| `internal/show` | Built-in shows: the original cartoons as Go functions (`dynamic.go` 41, `audio.go` 10, `test.go` 4, `library.go` shapes/palettes/`Keys`, generated by `tools/showgen`). Exported lists `show.Dynamic`, `show.Audio`, `show.Test` (`[]func() *keyframe.Cartoon`), `show.New(list)` builds a `*keyframe.Show`. `show_test.go` proves every frame equals the original data. The original Text cartoon was dropped on purpose |
| `internal/show/original` | Original show data (`original.bin`, 32100 B), only for showgen and tests (not linked into the firmware) |
| `internal/upload` | Vendor Download Software protocol (`A5 05 len payload check 55` frames: connect 17 → 16, records 32 → 33, end 37); keeps only the show region `0x0800B000..0x08032100` of the uploaded `Update.hex` |
| `internal/remote` | NEC IR decoder and the original button → action table |
| `internal/console` | Vendor console commands plus control commands `Key 1-7`, `Action 1-19`, `Cartoon i`, `Status` (through the `console.Remote` interface, implemented by `CubeRunner`); `Demux` splits UART1 between `#?cmd?!` console lines and vendor binary frames (`A5 05 … 55`, to `upload`); other bytes are ignored |
| `internal/registry` | Key/value registry in the vendor flash format |
| `internal/settings` | 8-byte EEPROM record at 0x00: speed, flags (play mode, pause, display off), mode, cartoon |
| `internal/cascade` | Original multi-cube protocol: frames `A5 55 len payload xor`; Hello / Number / Sync (show clock + key mask); `KeyMask` |
| `internal/component` | Hardware drivers (table below) |
| `internal/controller/board.go`, `console.go` | `YellowBoard` implements the `Board` interface and is the console's `Device` |
| `internal/controller/state.go` | `Mode` (Dynamic, Audio, Test, Text) and `StateTracker`: DIP switches (debounced 3 reads), `Show()` for the mode (Text = show slot), `Uploaded` blob, `SetUploading`, `Config2` |
| `internal/runner/runner.go`, `input.go` | Main loop (below) |
| `tools/fwtool` | Extracts the show data from a vendor hex or flash dump |
| `tools/showgen` | Show data → Go shows in `internal/show` |

`internal/component` drivers:

| File / type | Hardware |
|---|---|
| `pin.go` | `NewInputPin` (pull-down config + `High()`: the KEY5 workaround); `Button` (EXTI falling edge, 30 ms debounce, priority lowered to 0xE0); `DisableDebugPort` |
| `scanner.go` | TIM2 display refresh: 2 brightness planes, about 104 Hz |
| `microphone.go` | ADC1 continuous, sampled by a TIM3 1 kHz interrupt |
| `infrared.go` | PC6 EXTI on both edges, timed with `time.Now` |
| `serial.go` | UART1: bytes go through the console `Demux` (commands and upload frames) |
| `eeprom.go` | 24C02 through TinyGo's `machine.I2C0` (I2C1, PB6/PB7, 100 kHz, `NewOnBoardI2C` in `pin.go`) |
| `flash.go` | Page erase/program: registry (`WriteFlashPage`), show slot (`ShowSlotStorage`, `ShowSlotData`) |
| `cascade.go` | UART2 on PA2/PA3, 115200 |

Runner loop:

- **Boot.** `HandleStartup` prints `T='type'`; the runner prints `Config1/Config2`. Then cascade discovery runs (about 1.1 s, prints `DeviceType/DeviceNum`). `placeDisplay` positions the canvas from the cube number (uploaded show data has its own display table; Go shows are one cube at the origin).
- **Show by DIP mode** (`Tracker.Show()`): built-in lists for Dynamic/Audio/Test; Text plays the show slot (all 4 groups of the uploaded data) or nothing if the slot is empty. `handleInput` polls `Tracker.Update()`; a change prints `Mode=<name>` and switches the show live (`changeMode` → `setShow`, keeps settings, a paused player draws one frame first). An upload (`board.OnUpload`) empties the Text player first and reloads it when done.
- **Console `ReBoot`** is deferred 300 ms and cancelled by an upload connect frame (the vendor tool sends `reboot` then connects; the original rebooted into its bootloader).
- **Loop.** Only `runOnboard`: `Player.Step(clock)`, then wait on the shared show clock; `handleInput` returns true when the show was replaced to stop waiting.
- **Input (`handleInput`).** Every tick it:
  - polls the console and serial port,
  - saves settings 2 s after the last key action,
  - polls the DIP switches,
  - polls the cascade link,
  - turns key presses (and console `Key n`) into a `KeyMask`; the master sends Sync every 8 ms and runs actions on rising edges, slaves run actions from the master's mask,
  - handles IR on the master only.
- **KEY6** is a normal key (play model). The old long-press runner mode cycle and the Serial/Debug placeholder modes are removed.

## Status

### Confirmed on hardware (2026-10-19, logs in `.agent/live_testing/`)

- KEY5 fix (`gpio_diag`)
- Display, orientation and colors
- All 4 DIP groups
- All keys, including through the cascade key mask; no freezes after the EXTI priority change
- Audio group, IR remote (all buttons), boot log lines
- EEPROM settings (incl. pause/display off) survive a power cycle; original-style log messages
- Console: Help, ReBoot, Test, RegList, SetReg, GetReg match the original (`.agent/live_testing/*console.txt`)

Hardware facts: the base's Type-C is power only (no USB serial); there's no SD slot. The EEPROM type is `8x32x8b`, so the original runs as type 7 (`8x32x8B`). The original console is at 921600 baud and its boot log is in `.agent/live_testing/original output 921600.txt`.

### Implemented but NOT yet checked on hardware

The checklist is in `.agent/documentation/what-i-need-from-you.md`:

- Console leftovers: registry surviving a power cycle (flash write), ReadMem, DelReg, ClrReg. Use `scripts\terminal.cmd` (`tinygo monitor` can't send).
- Go shows (2026-10-19): Dynamic/Audio/Test now run from Go functions instead of the embedded data; should look identical.
- Text mode + upload: `scripts\upload.cmd`, download an `Update.hex` from the 3D Software; expect `Show upload started`, `Show upload OK`, then `Uploaded show: N cartoons` in Text mode. About 13000 records, ~2-3 min.
- Cascade (needs a second cube).
- Live DIP mode switching (no reset), KEY6 as a plain key, console `Key`/`Action`/`Cartoon`/`Status` (2026-10-19 refactor).

### Not ported on purpose

- The vendor IAP bootloader (uploads go straight to the show slot instead; firmware updates use CubeProgrammer).
- The serial-number copy protection.
- The IWDG watchdog (about 4 s in the original).
- Display `dir = 1`.
- Vendor GBK help texts (ours are English).
- The original console runs at 921600 baud; ours stays at 115200 (the Pico and the PC tools use 115200).

## Next work (in order)

1. **Act on the user's next hardware report** (Go shows, Text mode upload, console leftovers). If uploads stall: the 10 ms input tick limits the ack rate; the PC tool waits up to 200 ms per record.
2. **Sibling projects** (in `..\`). The user decided: no Pico-specific protocol; the Pico uses the console (`#?cmd?!`, `Key`/`Action`/`Cartoon`/`Status`) and the upload frames like the PC tools. `led-common-lib`'s `pkg/led` and `pkg/network` are obsolete; the engine (`internal/engine/*` except `oracle`) is to move into the library, plus a show data writer, a cube protocol client and a desktop ↔ Pico link protocol. Detailed instructions: `..\led-common-lib\.agent\stm32-alignment.md`, `..\rpi-led-communication\.agent\stm32-alignment.md` (touchscreen control + log view), `..\desktop-led-controller\.agent\stm32-alignment.md` (iCubeSmart app features, upload via Pico). Once the library has the engine, delete `internal/engine/{blob,render,player,keyframe,frame,audio}` here and import them (update `tools/showgen`).
3. **Optional extras:**
   - IWDG watchdog (feed it in every loop; keep the handshake time in mind)
   - display `dir = 1`
   - tuning the brightness planes (the original uses 3 planes weighted 1:3:6 out of 10)
4. **Upstream TinyGo issue.** The user has an open issue about `Pin.Configure`'s bit-4 leak. The fix is `uint32(config.Mode)&0xf` in `machine_stm32f103.go`. A PR or issue update can be drafted on request.
5. **Housekeeping.** `tmp/cmd/test` and `tmp/cmd/deez_nuts` (gitignored) import the deleted `internal/component/registers` package and no longer build.

## Where the knowledge lives

- `docs/documentation/`: the clean, committed reference (hardware, display, serial protocols, show data, engine, storage, firmware). It must not mention `.agent/` or how the findings were obtained. Keep it in sync when behavior changes.
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
