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

## Still needed

### 1. Install ILSpy (for the .NET PC software)

Either the GUI (ILSpy from its GitHub releases, or `winget install icsharpcode.ILSpy`) or the CLI (`dotnet tool install -g ilspycmd`, which needs the .NET SDK). Then decompile:
- `3D Software - RGB.exe`
- `Download Software.exe`

Output goes to `.agent/external/decompiled/<name>/`:
- GUI: *File → Open*, select the assembly, then *File → Save Code…* (saves as a C# project).
- CLI: `ilspycmd -p -o ".agent\external\decompiled\3D_Guan" ".agent\external\3D Software - RGB Cube\3D Software - RGB.exe"`

This gives the exact data-blob format, the graph codecs and the script compiler. It's the biggest shortcut for the port.

### 2. Permission to run Ghidra headless myself

Running installed tools isn't installing. I'd run `analyzeHeadless` from the scoop install, with a small export script under `.agent/tools/ghidra/`, on `Simulator.dll` and the firmware. Output goes to `.agent/re/`. Is that OK?

### 3. Quick hardware answers (next time the cube is on the desk)

- [ ] Plug **only** the base's Type-C into the PC. Does a new COM port appear? The vendor ships a CH340 driver, so the base may have USB-serial.
- [ ] Is there an SD-card slot on the back of the core board?
- [ ] Record the original behavior once. Flash `FullDump_512K.bin` back with the restore command in [flashing-and-dumping.md](flashing-and-dumping.md), then:
  - Save the UART console log at boot (`tinygo monitor -port COM4 -baudrate 115200`; if it's garbage, try 9600/57600/921600).
  - Send `help`, or just Enter, and save the output.
  - Take a short video per DIP group (11 / 01 / 10 / 00).
  - Note what each key does.
- [ ] Then flash the KEY5 experiment builds (E1–E4). I'll prepare them as `scripts\*.cmd` targets.

### 4. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).
