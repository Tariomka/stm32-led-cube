# Flashing, dumping and console: how-to on this PC

**Chosen method: UART system bootloader + STM32CubeProgrammer CLI on `COM4`.**
- ST-Link (SWD) is parked for later. It's at the end, for reference.
- No `make`. Helper scripts will be `.cmd` files.

## Wiring

Core-board 4-pin header (`GND / TXD / RXD / +5V`) to the USB-UART adapter (3.3 V logic):

| Core board | Adapter |
|---|---|
| GND | GND |
| TXD | TXD |
| RXD | RXD |
| +5V | 5V, **only if the cube isn't powered from Type-C** |

Unplug the Pico W from that header while flashing.

## Bootloader mode

1. Fit the **BOOT0 jumper**.
2. Press **Reset**. The ROM bootloader is now running (8E1, 115200).
3. Run the commands below. `-g 0x08000000` starts the app even with BOOT0 still fitted.
4. To boot normally after a power cycle: remove BOOT0, then press Reset.

## Commands (verified 2026-10-18 on COM4, CubeProgrammer v2.23.0)

```powershell
$cp = "C:\Program Files\STMicroelectronics\STM32Cube\STM32CubeProgrammer\bin\STM32_Programmer_CLI.exe"
& $cp -c port=COM4 br=115200                                                     # connect + chip info (Chip ID 0x414, 512 KB)
& $cp -c port=COM4 br=115200 -u 0x08000000 0x80000 firmware\FullDump_512K.bin    # read the whole flash (~70 s)
& $cp -c port=COM4 br=115200 -ob displ                                           # option bytes
& $cp -c port=COM4 br=115200 -w bin\main.hex -v                                  # write + verify (erases only the touched 2 KB pages)
& $cp -c port=COM4 br=115200 -g 0x08000000                                       # run the app
```

Chip facts from the run:
- Chip ID 0x414 (high density, 512 KB), bootloader protocol 1.0.
- RDP Level 0 (0xA5), no write protection.
- Writing the 49.4 KB TinyGo hex erased pages 0–24 (`0x08000000–0x0800C7FF`) and took about 8 s, plus about 7 s to verify.

### Restore the original firmware (if ever needed)

```powershell
& $cp -c port=COM4 br=115200 -e all -w firmware\FullDump_512K.bin 0x08000000 -v -g 0x08000000
```

## Planned `.cmd` helpers (next stage, will live in `scripts/`)

| Script | Does |
|---|---|
| `scripts\build.cmd` | `tinygo build` with the project's flags into `bin\main.hex` |
| `scripts\flash.cmd [COMx]` | build, then CubeProgrammer `-w … -v -g 0x08000000` |
| `scripts\monitor.cmd [COMx]` | `tinygo monitor -port COMx -baudrate 115200` |
| `scripts\dump.cmd [COMx]` | full 512 KB read into `firmware\dump_<timestamp>.bin` |

Default port is `COM4`.

## Removing the BOOT0 jumper step (later)

Our firmware gets a console command that jumps into the ROM bootloader:
- Disable IRQs and SysTick.
- Load MSP and the reset vector from system memory.
- `flash.cmd` sends that command, waits, flashes, then uses `-g` to run the new app.

The GD32F103 high-density system-memory address must be checked in the GD32F10x user manual ("Boot configuration"). On STM32F1 it is `0x1FFFF000`.

The vendor did the same with its own IAP: the app understands `reboot`, and the IAP then waits for `Update.hex` over UART ("IAP Waiting!"). We're replacing the IAP, so we use the ROM bootloader instead.

## Serial console

```powershell
tinygo monitor -port COM4 -baudrate 115200
```

The app uses 8N1, the ROM bootloader 8E1. Close the monitor before flashing, because the port can only be opened once.

## Parked: ST-Link + OpenOCD

Reference only. SWDIO = PA13 (KEY6 pad), SWCLK = PA14 (KEY5 pad), NRST = RESET pad, GND. The firmware disables SWJ, so use BOOT0 or *connect under reset*:

```powershell
openocd -f interface/stlink.cfg -c "set CPUTAPID 0x2ba01477" -f target/stm32f1x.cfg `
  -c "reset_config srst_only srst_nogate connect_assert_srst" -c "init; reset halt"
```
