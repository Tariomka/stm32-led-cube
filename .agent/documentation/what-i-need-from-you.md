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
- [ ] Then the KEY5 check: `scripts\build.cmd ./cmd/gpio_diag`, `scripts\flash.cmd COM4 bin\main.hex`, `scripts\monitor.cmd`. Follow the prompts and save the log (expected values are in each step's title).
- [ ] Then our firmware: `scripts\flash.cmd` (builds + flashes COM4). Check and report:
  - **Orientation:** DIP group Test (SW1/SW2 = 10). Do the test patterns sweep the same way as the original? If the cube is mirrored/rotated, `frame.Mapping` (`FlipX`, `FlipZ`, `FlipLayer`, `SwapXZ` in `controller.NewYellowBoard`) fixes it.
  - **Layer order:** does "up" in the shows point up? If not, `FlipLayer`.
  - **Colors:** red/green/blue in the right places (Test group has single-color fills).
  - **DIP polarity:** which DIP position plays which group? (Code assumes 1 = pin high = switch off.)
  - **Brightness/flicker:** 2 BAM planes, 400 µs base, ≈ 104 Hz full refresh. Any visible flicker or ghosting between layers?
  - **Keys:** each key's function vs. the original. KEY6 held ≥ 1 s cycles the runner mode.
  - **Speed:** compare a show's speed with the original video.

### 2. Later / optional

- ST-Link (we'll debug it together some day).
- Logic analyzer (exact scan timing of the original).
