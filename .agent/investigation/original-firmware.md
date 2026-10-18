# Original iCubeSmart firmware: what we have

## ✅ Complete original image secured (2026-10-18)

| Source | Content |
|---|---|
| `firmware/FullDump_512K.bin` (SHA-256 `057B70BE…6804`) | Full flash read **before** the TinyGo write: IAP + complete app |
| `.agent/external/3D Software - RGB Cube/1-3D8RGB-Demo-20250128.hex` | Vendor `update.hex` for this cube. **Byte-identical to the app in the dump** |
| `firmware/Start.hex` | IAP bootloader only |
| `firmware/Update.hex`, `firmware/ICubeSmart_Dump.hex` | ⚠ Truncated copies (first app segment only). Superseded |

Proof that the vendor hex is exactly what was on the chip:
- The IAP descriptor at `0x08004800` = `89762314` (magic), `0x08006000` (app address), MD5 `2D8075470D9E5FF90862F60100FFDD9C`, size `0x00037999`.
- The MD5 and size are those of the **`.hex` text file** that the IAP reads from SD. `1-3D8RGB-Demo-20250128.hex` is 227 737 B (`0x37999`) and its MD5 is `2D8075…DD9C`. Match.
- The app regions of that hex and of the dump: 0 byte differences.

> Correction to an earlier note: `0x37999` is the hex-file size, not the binary size. The app does **not** extend to `0x0803D999`.

## Flash layout of the original product

```text
0x08000000–0x080045FF  IAP bootloader "IAP 2.10, Create By Fan" (SD update.hex + UART update)
0x08004800–0x0800481B  IAP descriptor (magic, app addr, MD5 + size of update.hex)
0x08006000–0x08012D63  App vector table + app code (~20 KB) + COMPILED SHOW DATA (from ~0x0800B000)
0x08032100–0x08038FC3  Firmware library part: startup (SystemInit @0x08036768), drivers, IRQ handlers, ...
RAM: initial SP 0x2000F3C0, i.e. uses ~62 KB of the 64 KB SRAM
```

The hex has two parts: the fixed firmware **head and tail**, with the **compiled show data** in between. That matches how the PC software builds it (next section).

## Key discovery: the PC "3D Software" generates the firmware

`.agent/external/3D Software - RGB Cube/` contains:

| File | What it is |
|---|---|
| `3D Software - RGB.exe` | .NET 4 WinForms app (assembly `3D_Guan` v2.2.1.0, x86). Show editor, script compiler, **firmware generator** |
| `Simulator.dll` | **Native x86 DLL. Most likely the firmware's show engine compiled for the PC.** Exports: `SetCartoon`, `RunCartoon`, `EndCartoon`, `DoKeyFunction`, `DoKey1..7`, `GetCube`, `AddCube`, `CleanCube`, `SetGraphFun`, `DrawBallFun`, `DrawTitleFun`, `GetSourcePath`, `GetKeyChar`, `GetBuff`, `Test` |
| `*.gc` | Show projects (.NET `BinaryFormatter` of `_3D_Guan.Cartoon…` objects). `1-3D8RGB-Demo-20250128.gc` is the source of our cube's shows |
| `*.hex` | Ready firmware per cube type (3D8RGB, 8X8X32RGB, 3D16RGB, 3D16S/MINI, 3D8P) |
| `Download Software.exe` | .NET "STM32 Debug Tool v2.2". Serial console + uploads `Update.hex` to the IAP (sends `reboot`, then streams the hex) |
| `Downloader Driver/ch341ser.exe` | CH340 USB-serial driver |
| `3D Software Instructions -RGB Cube-English.pdf` | User manual (not yet read) |

### What reflection tells us about the PC app (metadata only, nothing executed)

- **Show model:**
  - `Cartoon` (types `Normal`, `Script`, `Audio`) is made of `CartoonFrames`, which hold `FrameElement`s.
  - Element types: `SetGraph` (a graph/bitmap at x,y,z with `RotateAxis` + angle), `DrawBall` (center, radius, in/out color, hollow), `TitleStyle` (scrolling text with `PlayStyle`).
  - Resources: `GraphData` (3-D bitmaps), `PaletteData` (colors), `IndexData` (index tables).
  - `CartoonConfig` has `m_key1..m_key7` (configurable key functions) and `m_disp`.
- **Cube types:** `T_8_8_8`, `T_8_16_8`, `T_16_16_16`, `T_8_32_8`.
- **Graph compression** (`Compile.Compress`): `NoCompress`, `Normal_2/4/8`, `Continue_2/4/8` (palette + RLE-style codecs with 2/4/8-bit indices).
- **Script language** compiled to a **stack VM bytecode** (`CompileScript.CommandElement.ElementType`):
  `LDR_SYSTEM/GLOBAL/LOCAL(_ADD)_A/B`, `LDR_STACK_BEGIN/END`, `STR_*`, `SWAP_A_B`, `PUSH/POP(_A)`, `GOTO`, `CALL_BEFORE/SYSTEM/GLOBAL/LOCAL`, `RET_FUNCTION/EVENT`, `CREATE_REF/DESTORY_REF`, arithmetic/logic ops (`minus … small_e_then` = 32–54), `COPY`, `COPY_STR`, `CAT`, `JZ`, `JNZ`.
- **Firmware build:**
  - `Compile.CompileAllToArray` produces the data blob.
  - `PaddingHex(bytes, head, tail)` + `ParseHexHeadAndTail` merge it into the firmware hex template.
  - `CompileFirmware` writes the `.hex`.

### Consequences for the port

1. The firmware is a **generic show engine + interpreter**, and all "premade light shows" are **data** in one blob.
   - Porting = re-implementing the engine (the data decoder, the frame/element renderer, the script VM, audio mode, keys) in Go, then embedding the **original blob unchanged**.
   - Every original show then runs as is, and new shows made in the PC app can be dropped in too.
2. `Simulator.dll` is the best RE target. It is x86 (much nicer decompilation than Thumb), has named exports, and is likely built from the same C source. It also gives an **executable oracle**: a Windows/386 Go test harness can load the DLL, feed it the blob, step it, read frames with `GetCube`, and compare them frame by frame with our Go engine.
3. The PC app (.NET) can be decompiled to C# with ILSpy. That gives the exact **blob format writer**, the **compression codecs** and the **bytecode compiler**, so the formats don't have to be guessed.

## Strings in the app (engine features)

| Area | Strings |
|---|---|
| Show engine | `Cartoon End`, `Change Cartoon %d`, `Replay Cartoon`, `Run x=%d,y=%d,z=%d,dir=%d,type=%d`, `Warn: Pos > Cubes`, `SetData error!` |
| Resource names (data blob) | `cubeamp1..4`, `Circle1..5`, `waterfall1/2`, `rain1/2`, `square1..6`, `Fly_1..4`, `Wave1/2`, `AudioG1/2`, `HELLO`, `HelloiCube`, `Demo-Text`, … plus color series `red-/green-/blue-/yellow-/cyan-/RGB-series` |
| Keys | `Previous Cartoon`, `Next Cartoon`, `Add/Sub speed to %.2f`, `Speed reset to %.2f`, `Play Model=%d`, `Display off=%d` |
| Cube config | `8x8x8A/B`, `8x16x8A`, `161616A/B/C`, `8x32x8A/B`, `Config1=%d,Config2=%d`, `DeviceType:%d,DeviceNum:%d` |
| Settings (24C02) | `Regedit …`, `ReadReg/WriteReg/DelReg`, `Commands explain:` (UART console) |
| Misc | `tickTimer`, `DisplayAdd=0x%.8X`, `IRQ:%s`, `Error Bus!`, `Error: Stack Out!` (the VM's stack check) |

## Existing decompilations in `tmp/decompiled/`

| File | What it actually is |
|---|---|
| `simplified.c` | Wrong architecture (RISC-V). ❌ Discard |
| `iCubeSmart.c` | IAP bootloader (ARM) |
| `stm32_cube.c` | **Your own TinyGo build.** It proved the KEY5 bug |

The vendor app and `Simulator.dll` have not been decompiled yet. That is Phase 1 of the porting plan.
