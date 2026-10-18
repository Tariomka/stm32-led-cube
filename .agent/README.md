# Agent workspace: LED cube

Notes, plans and reference material for:
1. porting the original iCubeSmart firmware to TinyGo
2. fixing KEY5 (PA14)

No project code has been changed so far.

## Status (2026-10-18)

| Item | State |
|---|---|
| KEY5 / PA14 root cause | ✅ Found (TinyGo `Pin.Configure` bit-4 leak). Hardware confirmation pending |
| Original firmware analysis | ⚠ The dump is **incomplete** (only about 23 % of the app). A full 512 KB dump is needed |
| Hardware / pin map | ✅ From schematics + photos |
| Flashing on this PC | ⏳ Needs a USB-UART adapter + stm32flash or CubeProgrammer |
| Porting | ⏳ Planned, waiting on the dump |

## Index

- **Plans**
  - [plan/key5-fix-plan.md](plan/key5-fix-plan.md)
  - [plan/porting-plan.md](plan/porting-plan.md)
- **Investigation**
  - [investigation/key5-pa14.md](investigation/key5-pa14.md): root cause, evidence, workaround, experiments
  - [investigation/original-firmware.md](investigation/original-firmware.md): what's in `firmware/`, what's missing, what the app does
  - [investigation/hardware.md](investigation/hardware.md): your exact boards (from the photos)
  - [investigation/environment.md](investigation/environment.md): toolchain on this PC, build quirks, target memory limits
- **Documentation**
  - [documentation/what-i-need-from-you.md](documentation/what-i-need-from-you.md): **checklist for tomorrow**
  - [documentation/flashing-and-dumping.md](documentation/flashing-and-dumping.md): UART bootloader / ST-Link / console how-to
  - [documentation/mcu-reference.md](documentation/mcu-reference.md): GD32/STM32F1 register cheat sheet, display chain timing
  - [documentation/schematics/README.md](documentation/schematics/README.md): pin maps read from the schematics (+ rendered PNGs)
- **Tools**
  - [tools/render-pdf.ps1](tools/render-pdf.ps1): PDF to PNG (+ tiles) using the built-in Windows renderer
- **Photos:** [physical/](physical/)
