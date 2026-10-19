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
- [ ] Console leftovers: power-cycle after `SetReg name LedCube`, then `GetReg name` (must survive: proves the flash write); `ReadMem 0x08000000 0 32`; `DelReg name`; `ClrReg`. Note: the flash layout changed (registry stays at `0x0807F000`), so the registry should still be there after reflashing.
- [ ] **Go shows**: Dynamic (11), Audio (01) and Test (10) now run from Go functions instead of the embedded show data. They should look exactly the same.
- [ ] **Text mode (00)**: with nothing uploaded the cube stays dark and the log says `No uploaded show`.
- [ ] **Upload**: in "3D Software - RGB" export `Update.hex` (or use a demo hex from `.agent/external/3D Software - RGB Cube/`), run `scripts\upload.cmd COM4` (close the terminal/monitor first), pick the file and download. Expect `Show upload started` … `Show upload OK` on the console (~2-3 min), then in Text mode `Uploaded show: N cartoons` and the show playing; it must survive a power cycle.
- [x] Original-style log messages and pause/display off restored after a power cycle.
- [ ] **Cascade**: needs a second cube. Connect cube A's jack TX to cube B's RX, power both together. B must print `DeviceType=0,DeviceNum=1` and play in step with A, following A's keys.

### 2. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).
