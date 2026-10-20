# What I need from you

## Done ✅

- Full 512 KB dump (`firmware/FullDump_512K.bin`). It is complete and matches the vendor hex.
- Installed: STM32CubeProgrammer, Ghidra, JDK 21. Port: **COM4**.
- PC "3D Software" in `.agent/external/`.
- Decisions:
  - our firmware owns the whole flash
  - vendor data may be committed (you commit)
  - port **all** original features
  - keep the Pico/desktop integration
  - no make (use `.cmd`)
  - no installs by the agent

- Decompiled C# of the PC software in `.agent/external/decompiled/`.
- Ghidra permission (used: decompiles in `.agent/external/re/`).
- `.agent/external` and `.agent/photos` are gitignored.

## Still needed

### 1. Hardware session results (2026-10-19) and what's left

Logs: `.agent/live_testing/` (original at 921600 baud, ours at 115200).

- [x] Base Type-C: power only, no COM port. No SD-card slot on either board.
- [x] Original boot log captured at 921600 baud (`original output 921600.txt`). EEPROM type is `8x32x8b`, so the original runs as type 7 (`8x32x8B`).
- [x] KEY5 check and our firmware (orientation, colors, DIP groups, keys): confirmed working.
- [x] Audio group reacts to loudness.
- [x] IR remote: all buttons work.
- [x] Keys: no freezes; short taps, long presses and KEY6 mode cycling work through the key mask.
- [x] Boot log matches the expected lines.
- [x] EEPROM settings survive a power cycle. Pause and display off are now saved too (they weren't before); a saved pause shows the first frame, then pauses.
- [x] **Console** (`console.txt` logs, ours vs original): Help (same layout, English texts), ReBoot, Test, RegList, SetReg, GetReg all behave like the original.
- [x] Console leftovers: registry survives power cycles, `DelReg`, `ClrReg` work (`upload.txt`, `before folding start.txt`).
- [x] `ReadMem` works (`readmem.txt`). Once, right after power-up, a correct line printed only the Help line; not reproducible since. Probably garbage bytes from the boot (or a paste during boot) ended up in the line. Not pursued.
- [x] **Go shows**: Dynamic, Audio and Test look the same as before.
- [x] **Text mode (00)**: with nothing uploaded the cube stays dark and the log says `No uploaded show`.
- [x] **Upload**: `Show upload OK`, `Uploaded show: 18 cartoons`, survives a power cycle (`upload.txt`, `registry after upload.txt`). The reboot during the first attempt was the user resetting the cube.
- [x] `Uploaded show: 0 cartoons` in `logs after latest version.txt`: the user removed all cartoons on purpose; works as expected.
- [x] Original-style log messages and pause/display off restored after a power cycle.
- [~] **Cascade**: no second cube available; assumed to work until someone reports an issue.
- [x] **Live DIP modes** (`logs after latest version.txt`).
- [x] **KEY6**: only the play model key.
- [x] **EEPROM on hardware I2C**: works. Before the user folded the `Start()` calls into the constructors the boot printed `T=''` (`before folding start.txt`), after it `T='8x32x8b'`. Suspected: TinyGo's F1 I2C driver isn't interrupt-safe in its read sequence and the boot timing of the TIM2/TIM3 interrupts decides it. Watch for `T=''` or settings not restored.
- [x] **Control commands**.

### 2. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).
