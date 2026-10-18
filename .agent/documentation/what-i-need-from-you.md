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
- [ ] Then flash the KEY5 experiment builds (E1–E4). I'll prepare them as `scripts\*.cmd` targets.

### 2. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).

### 3. Decision: may I start writing code?

The next steps are code, not docs. Nothing gets staged or committed by me.
- Host-side Go packages for the engine (`internal/engine/...`, testable with `go test`)
- The custom TinyGo target
- `scripts\*.cmd`
- The KEY5 fix
- An oracle test harness that loads `Simulator.dll` (windows/386)
