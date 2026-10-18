# Original iCubeSmart firmware: what we have

## Files in `firmware/`

| File | Range | Size | What it is |
|---|---|---|---|
| `Start.hex` | `0x08000000–0x0800481F` | 18 464 B | **IAP bootloader** "IAP version: 2.10 / Create By Fan" |
| `Update.hex` | `0x08006000–0x08012D7F` | 52 608 B | **Application** (first part only, see below) |
| `ICubeSmart_Dump.hex` | `0x08000000–0x08013FFF` | 81 920 B | Byte-identical to `Start.hex` + `Update.hex`, padded with `0xFF` |

## ⚠ The application is incomplete

- The app vector table at `0x08006000` points far beyond the end of the dump:
  - Reset handler `0x08006415`. Its literal pool calls `SystemInit` at **`0x08036769`**, then `__main` at `0x08006131`.
  - NMI `0x080336CD`, HardFault `0x08032CD1`, most IRQ handlers at `0x08032xxx–0x08037xxx`.
- The IAP keeps an image descriptor at `0x08004800`:
  - `89 76 23 14` (magic?)
  - `00 60 00 08` = app address `0x08006000`
  - 16 bytes at `0x08004808` (MD5 of the app, passed to the IAP's MD5 check)
  - `99 79 03 00` at `0x08004818`, so **size = `0x00037999`** (227 737 B). The full app therefore spans `0x08006000–0x0803D998`.
- So `Update.hex` and the dump contain only about 23 % of the application. Missing are the startup code, clock/peripheral drivers, interrupt handlers, and probably more animation data.
- The dump stopping at exactly `0x14000` (80 KB) looks like a read-size setting in the tool that made it.

### Recovering the rest

1. **Best: dump the whole 512 KB from the chip** (see [../documentation/flashing-and-dumping.md](../documentation/flashing-and-dumping.md)).
   - If you only ever erased the sectors you programmed (CubeProgrammer's default), the region from `0x08032000` up is probably still intact. Your TinyGo builds are about 50 KB, so they only overwrite `0x08000000–0x0800C7FF`.
   - **Do this before flashing anything else.**
2. Ask the vendor again for the *complete* `update.hex`. It should cover about 227 KB of binary.
3. Fallback reference: Sliicy's repo has `Firmware/RGB/Original Firmware 8x8x8 RGB Cube.hex`.
   - It is a **different (older?) build**, at `0x08006000–0x0801110F` + `0x08032100–0x08037863`, with one malformed record at `0x08011100`.
   - It has the same structure: app code + data low, library code high from `0x08032100`. It can't be stitched into ours (addresses differ), but it is a complete second sample for reverse engineering.

## Memory layout of the original product

```text
0x08000000  IAP bootloader (Start.hex)          ~18 KB
0x08004800  IAP image descriptor (magic, app address, MD5, size)
0x08006000  App vector table, app code (~20 KB), then animation data ("Cartoon" scripts, patterns)
   ...
0x08032100  Library / drivers / startup (by analogy with the Sliicy build)
0x0803D999  End of app (per IAP descriptor)
RAM: SP init 0x2000F3C0, so it uses almost all 64 KB of SRAM
```

The IAP:
- Prints a banner on UART1.
- Mounts the SD card (FAT), looks for `update.hex`, and compares its MD5 against the stored one ("Md5 the same!"). Otherwise it reflashes the app ("Start SD Update!").
- Then sets `VTOR = 0x08006000` and jumps to the app.

## What the app does (from strings in the dump)

| Area | Evidence (strings) | Notes |
|---|---|---|
| Animation engine | `Cartoon End`, `Change Cartoon %d`, `Replay Cartoon`, `Run x=%d,y=%d,z=%d,dir=%d,type=%d`, `Warn: Pos > Cubes`, `SetData error!` | Data-driven: scripts reference named **patterns** and **color series** and move them around |
| Pattern names | `cubeamp1..4`, `Circle1..5`, `Circle-2..5`, `waterfall1/2`, `rain1/2`, `square1..6`, `Arrows5`, `TurnRight1/2`, `Fly_1..4`, `Wave1/2`, `frame1..4`, `HELLO`, `HelloiCube`, `happyB`, `B/G/R-CUBE`, `4strip`, `vertical8dot`, `crosswise8dot`, `8x8test`, `meteor`, `AllAround`, `TwoSide`, `Demo-Text`, ... | Bitmaps live next to their names (`~BBBB~`-style 8×8 glyph bytes are visible) |
| Colors | `green, yellow, purple, white, red`, `red/green/blue/yellow/cyan/RGB-series` | Implies several brightness levels per color, i.e. **PWM/BAM color depth** |
| Audio | `AudioG1`, `AudioG2`, `Wave1/2` | Mic on PA0 (ADC). DIP group 2 = "Audio Effects" |
| Buttons | `Previous Cartoon`, `Next Cartoon`, `Add speed to %.2f`, `Sub speed to %.2f`, `Speed reset to %.2f`, `Play Model=%d`, `Display off=%d` | Matches the KEY1–KEY7 silkscreen |
| Cube types | `8x8x8A`, `8x8x8B`, `8x16x8A`, `161616A/B/C`, `8x32x8A/B`, `Config1=%d,Config2=%d`, `DeviceType:%d,DeviceNum:%d` | One firmware for many bases. Your base is "8x8x" with 8/16/32 layer options |
| Settings store ("Regedit") | `Clean Regedit OK!`, `Delete Regedit:%s OK!`, `ReadReg:%s=%s`, `WriteReg:%s OK!`, `Regedit Count = %d`, `Regedit format error!`, `Commands explain:` | Key/value store, most likely in the 24C02 EEPROM, editable over the UART console. "The button function can be modified in the 3D software" probably goes through this |
| Debug console | `ReadMemery End!`, `Address=0x%.8X,offset=%d,count=%d`, `IRQ:%s`, `Error Bus!`, `Error: Stack Out!` | Text command shell on UART |
| Misc | `tickTimer`, `DisplayAdd=0x%.8X` | |

## Existing decompilations in `tmp/decompiled/`

| File | What it actually is | Useful? |
|---|---|---|
| `simplified.c` | Decoded as **RISC-V** (registers `fa0`, `s0`, ...), i.e. wrong architecture | ❌ Discard |
| `iCubeSmart.c` | Correct ARM decompile of the **IAP bootloader** (SDIO, FLASH, FAT, banner strings) | Only for the IAP |
| `stm32_cube.c` | Decompile of **your own TinyGo build** (TinyGo panic strings, PA14 print code) | ✅ It's what proved the KEY5 bug. Not the vendor app |

So **the vendor app has never been decompiled yet.** Plan: Ghidra (ARM Cortex-M Thumb, LE), image base `0x08006000`, plus a peripheral map. See [../plan/porting-plan.md](../plan/porting-plan.md).
