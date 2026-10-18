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

### 1. Hardware session (when the cube is back on the desk)

- [ ] Plug **only** the base's Type-C into the PC. Does a new COM port appear? The vendor ships a CH340 driver, so the base may have USB-serial.
- [ ] Is there an SD-card slot on the back of the core board?
- [ ] Record the original behavior once. Flash `FullDump_512K.bin` back with the restore command in [flashing-and-dumping.md](flashing-and-dumping.md), then:
  - Save the UART console log at boot (`tinygo monitor -port COM4 -baudrate 115200`; if it's garbage, try 9600/57600/921600).
  - Send `help`, or just Enter, and save the output.
  - Take a short video per DIP group (11 / 01 / 10 / 00).
  - Note what each key does.
- [x] KEY5 check and our firmware (orientation, colors, DIP groups, keys): confirmed working.
- [ ] **Audio group** (DIP SW1/SW2 = 01): play music / clap near the cube. The cartoons should react to loudness. Report if they are always flat (level 0) or always full.
- [ ] **IR remote** (if you have the kit's remote): press each button. The monitor prints `IR <code>`. Check CH- (display off), CH+ (play mode), VOL-/EQ (previous/next), VOL+ (pause), 0..9 (jump to cartoon). Send me any codes that don't match.
- [ ] Keys after the interrupt priority change: press keys quickly and often for a while; the cube must never freeze.
- [ ] Keys now act through the cascade key mask (up to ~10 ms later): check short taps and long presses still work, and KEY6 held 1 s still cycles the mode.
- [ ] **Boot log** (`scripts\monitor.cmd`): expect `T='8x8x8A'` (or similar), `Config1=..`, `DeviceType=1,DeviceNum=0`, `Run x=0,...`. Boot now takes ~1.1 s longer (cascade discovery, like the original).
- [ ] **Console**: in the monitor type `#?Help?!` then `#?SetReg name cube?!`, `#?RegList?!`, power-cycle, `#?GetReg name?!` (must survive), `#?ReadMem 0x08000000 0 32?!`, `#?ReBoot?!`. Don't run `SetType`/`SetSerial` unless you want to overwrite the vendor records.
- [ ] **EEPROM settings**: change speed / play mode / cartoon with the keys, wait 3 s, power-cycle: they must be restored.
- [ ] **Cascade** (only if you get a second cube): connect cube A's jack TX to cube B's RX, power both together. B must print `DeviceType=0,DeviceNum=1` and play in step with A, following A's keys.

### 2. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).
