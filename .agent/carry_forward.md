# Carry-forward prompt (next session)

Paste this to the next agent, or read it first if you are that agent.

---

You are continuing work on **stm32-led-cube**: TinyGo firmware for an iCubeSmart 3D8RGB 8x8x8 RGB LED cube (GD32F103RET6 "yellow board": 512 KB flash, 64 KB RAM). The repo is `d:\Git\some_personal_shit\Led_Cube\stm32-led-cube`. The goal was to **port all of the original vendor firmware** (every premade show plus all hardware features) into this project and fix KEY5 (PA14). The port is **code-complete**. The Pico and the desktop app now talk to the cube with the same console + upload protocol as the PC tools (no separate Pico protocol). What remains is hardware verification and the sibling projects' alignment (instructions in each sibling's `.agent/stm32-alignment.md`).

## Hard rules from the user

- **Never install software.** **Never `git add`/commit/push.** The user reviews and commits.
- Use the Go task runner `go run . <task>` (root `main.go`, Windows only), not make (the `Makefile` is kept for Linux only). The old `scripts\*.cmd` are gone.
- Tests live in `tests/` at the root (package `tests`, external tests of the `internal/...` packages), never next to the implementation files.
- **Never write to `%TEMP%`.** One-off scripts and files go in `.agent/external/temp` (gitignored).
- `.agent/` is the agent's memory and docs folder. `.agent/external` and `.agent/photos` are gitignored.
- Don't run ILSpy. The C# decompile is already in `.agent/external/decompiled`. Ghidra is allowed.
- Flashing: STM32CubeProgrammer CLI over **COM4** with the BOOT0 jumper (`go run . flash`). ST-Link and the logic analyzer are not used.
- Don't edit files via terminal commands (no `gofmt -w`, `sed`, etc.). Use the edit tools.
- The user may edit files between sessions. Always re-read a file before editing it.

## Environment quirks

- In agent terminals, strip TinyGo's goroot from PATH before `go`/`tinygo`, otherwise you get "sync/atomic is not in std":
  `$env:Path = ($env:Path -split ';' | Where-Object { $_ -notmatch 'tinygo\\goroot-' }) -join ';'`
- Build from the **repo root**: `go run . build [pkg]`. The default pkg is `./cmd/8x8_rgb_cube` and the output is `bin\main.hex`.
  - It uses the custom target `targets/gd32f103re.json` with linker script `targets/gd32f103re.ld`, which gives FLASH 350K.
  - Above that: uploaded show slot at `0x08057800` (158 KB), console registry at `0x0807F000` (+ backup `0x0807F800`).
- Other tasks (`go run .` lists them; defaults COM4, 115200):
  - `flash [COMx] [file]`: builds first unless a file is given.
  - `dump [COMx]`, `monitor [COMx] [baud]`, `upload [COMx]` (vendor Download Software at 115200).
  - `terminal [-raw] [COMx] [baud] [command]`: console terminal (wraps lines in `#?…?!`); serial port via kernel32 (`GetCommState`/`SetCommState`/`SetCommTimeouts`), no extra modules.
  - `test` (`go test ./tests/...`), `oracle` (`GOARCH=386`, adds the `Simulator.dll` comparison in `tests/oracle_test.go`).
- Host tests: `go test ./tests/...` (only `lightshow` builds on the PC).
- Regenerate the Go shows: `go run ./tools/showgen` (overwrites `internal/lightshow/builtin/{dynamic,audio,test,library}.go`; `show.go` is hand-written).
- **Known false alarm:** `gofmt -l` flags every file because of CRLF line endings (git autocrlf). To check real formatting, normalize CRLF→LF into `.agent/external/temp/fmt.go` and run `gofmt -d` on that copy.
- `create_file` cannot overwrite a file. Write `x.go.new`, then `Move-Item -Force`.
- **Check new Go files.** A tool sometimes inserted a duplicate `package` line at the top of a new file. Check the first line after creating one.
- **Avoid `strings.EqualFold` / `strings.Fields` in firmware code.** They pull Unicode tables into RAM (+6 KB).
- Flash budget: the last build was about 134 KB flash and 10.3 KB static RAM.

## Toolchain

TinyGo 0.43.0 (scoop), Go 1.27.2, STM32CubeProgrammer v2.23.0, Ghidra + temurin21-jdk (scoop).

## Architecture

Layers (2026-10-20 restructure, user's design): `runner` → `controller` → {`lightshow`, `component`} → `common`. `runner` only creates controllers and runs the loop; each controller owns its components; `lightshow` knows nothing about hardware; `component` = drivers only; `common` imports nothing from the project. The old `YellowBoard`/`Board` interface, `global`, `utils`, `RunnerConfig` and the packages `console`, `upload`, `registry`, `settings`, `cascade`, `remote`, `engine`, `show` are gone (their tests were deleted on purpose; only `lightshow` tests remain). No `led-common-lib` dependency anymore.

| Path | What |
|---|---|
| `cmd/8x8_rgb_cube` | `runner.New().Start()` |
| `cmd/gpio_diag` | KEY5/GPIO diagnostic program |
| `internal/runner/runner.go` | Creates the controllers (System first), `Start`: banner, `Config1/2`, cascade discovery, `show.Place`, loop `step()` |
| `internal/controller/system.go` | `SystemController`: debug port off, watchdog (start/feed, `Watchdog reset!`), LEDs, `Startup` banner, `Reset`, `ReadMemory` |
| `internal/controller/settings.go` | `SettingsController`: EEPROM settings record (0x00), cube type (cached, `CubeTypeIndex`, `IgnoresDisplayDir`), serial |
| `internal/controller/display.go` | `DisplayController`: scanner, mapping, `SetUpsideDown`, `SetBlank`, `Show(voxels)` |
| `internal/controller/keys.go` | `KeyController`: 7 buttons, `Take`, `Press` (console Key), `Held` |
| `internal/controller/remote.go` | `RemoteController`: NEC decoder (ISR via `component.InfraredReceiver` pulse callback), button → action, `Take` prints `IR:` |
| `internal/controller/mode.go` | `ModeController`: DIP switches (3-read debounce), `Show()` (builtin or show slot, logs Text-mode lines), `Uploaded()`, `SetUploading`, `Config()` |
| `internal/controller/cascade.go` | `CascadeController`: UART2 frames, `Discover(idle)`, `Clock`, `Keys(pressed, held)` → edges (master Sync 8 ms / slave), key mask |
| `internal/controller/show.go` | `ShowController`: player + canvas + microphone, `Step`, `Poll` (mode change, settings save), `Changed`, `Place`, `RunKeys`, `Action`, `Select`, `Status`, `UploadStarted/Finished` |
| `internal/controller/terminal.go`, `terminal_commands.go`, `upload.go`, `registry.go` | `TerminalController`: UART1 byte splitter (`#?…?!` / A5 05 frames), console commands, deferred ReBoot, upload receiver into the show slot, registry |
| `internal/lightshow` | Player (package `lightshow`), `Show`/`Cartoon`, `BlobShow`, VM; `SetShow` sets the missing-resource log on any show with `SetMissing` |
| `internal/lightshow/blob`, `render`, `keyframe` | Show data parser, renderers, Go show model (`keyframe.Show{Cartoons, Keys, Library}`) |
| `internal/lightshow/builtin` | Generated built-in cartoons (`Dynamic`, `Audio`, `Test`, `New(list)`), `original/` factory data (tests/showgen only) |
| `tests/` | Host tests (package `tests`): `blob_test.go`, `builtin_test.go`, `oracle_test.go` (frame-exact against `Simulator.dll`, windows && 386) |
| `internal/component` | Drivers (table below) |
| `internal/common` | `EqualFold` (ASCII), `Flag` ("1"/"0") |
| `tools/fwtool` | Extracts the show data from a vendor hex or flash dump |
| `tools/showgen` | Show data → `internal/lightshow/builtin` |

`internal/component` drivers:

| File / type | Hardware |
|---|---|
| `pin.go` | `NewInputPin` (pull-down config + `High()`: the KEY5 workaround), `NewUart`, `NewSpiOutput`, `NewOnBoardI2C`, `DisableDebugPort` |
| `button.go` | `Button` (EXTI falling edge, 30 ms debounce, priority lowered to 0xE0) |
| `scanner.go`, `frame.go` | TIM2 display refresh: 3 thermometer-coded planes weighted 1:3:6 (120/360/720 µs), about 104 Hz; `Show(voxels, Mapping)` encodes and swaps |
| `demultiplexer.go` | 74HC154 layer select |
| `microphone.go` | ADC1 continuous, TIM3 1 kHz ISR, loudness meter (`Level(mode)`) |
| `infrared.go` | PC6 EXTI on both edges, timed with `time.Now`, passes mark/space to a callback |
| `serial.go` | `SerialPort`: `NewConsolePort` (UART1), `NewCascadePort` (UART2), `Read`/`Write` |
| `eeprom.go` | 24C02 through TinyGo's `machine.I2C0` (100 kHz), one retry after re-configuring the bus |
| `flash.go` | Page erase/program: registry (`WriteFlashPage`), show slot (`ShowSlotStorage`, `ShowSlotData`) |
| `watchdog.go`, `onboardLeds.go` | IWDG via `machine.Watchdog`; status LEDs |

Main loop (`runner.step`): `show.Step(clock)`, then every 10 ms until the frame time is over: `system.Feed`, `terminal.Poll`, `show.Poll`, `show.RunKeys(cascade.Keys(keys.Take(), keys.Held()))`, master: `show.Action(remote.Take())`; stop early when `show.Changed()`. KEY6 is a normal key; console `ReBoot` is deferred 300 ms and cancelled by an upload connect frame.

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

The checklist is in `.agent/documentation/what-i-need-from-you.md`. Everything else, including the 2026-10-20 controller restructure, was confirmed on hardware by the user.

- Cascade (needs a second cube; the user assumes it works until someone reports an issue).
- The `go run . terminal` / `monitor` serial code (Win32 API) replacing the PowerShell terminal (2026-10-20): not run against the cube yet.

### Not ported on purpose

- The vendor IAP bootloader (uploads go straight to the show slot instead; firmware updates use CubeProgrammer).
- The serial-number copy protection.
- The IWDG watchdog is now ported (`component/watchdog.go`, 4 s requested ≈ 3.2 s real, fed by `SystemController.Feed` from the runner's idle step and the discovery loop; boot prints `Watchdog reset!` after a watchdog reset).
- Display `dir = 1` is ported: upside down (x and layers mirrored) except for the types whose factory dir-1 conversion is a branch to dir 0 (161616B, 161616C, 8x32x8B = this kit). Decoded from the type tables at `0x080380D8–0x08038148` (entry +8 = dir 0, +0xC = dir 1).
- Vendor GBK help texts (ours are English).
- The original console runs at 921600 baud; ours stays at 115200 (the Pico and the PC tools use 115200).

## Next work (in order)

1. **Act on the user's next hardware report** (Go shows, Text mode upload, console leftovers). If uploads stall: the 10 ms input tick limits the ack rate; the PC tool waits up to 200 ms per record.
2. **Sibling projects** (in `..\`). The user decided: no Pico-specific protocol; the Pico uses the console (`#?cmd?!`, `Key`/`Action`/`Cartoon`/`Status`) and the upload frames like the PC tools. `led-common-lib`'s `pkg/led` and `pkg/network` are obsolete; the engine (`internal/lightshow` minus `builtin`/`oracle`) is to move into the library, plus a show data writer, a cube protocol client and a desktop ↔ Pico link protocol. Detailed instructions: `..\led-common-lib\.agent\stm32-alignment.md`, `..\rpi-led-communication\.agent\stm32-alignment.md` (touchscreen control + log view), `..\desktop-led-controller\.agent\stm32-alignment.md` (iCubeSmart app features, upload via Pico). Once the library has the engine, delete `internal/engine/{blob,render,player,keyframe,frame,audio}` here and import them (update `tools/showgen`).
3. **Optional extras:** done and confirmed on hardware (watchdog, display `dir = 1`, 3-plane 1:3:6 brightness like the factory thermometer coding, EEPROM retry after an I2C reset).
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
