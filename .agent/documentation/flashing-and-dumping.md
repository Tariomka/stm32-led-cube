# Flashing, dumping and console: how-to on this PC

There are three ways to talk to the chip. **Method A needs no extra hardware besides a USB-UART adapter**, and it is how you worked before.

| | A. UART bootloader | B. ST-Link + OpenOCD (SWD) | C. Software jump to bootloader (future) |
|---|---|---|---|
| Hardware | USB-UART adapter (3.3 V logic) | ST-Link V2 (or clone) + 4 wires | Same as A, no jumper |
| Software | STM32CubeProgrammer **or** stm32flash | OpenOCD (installed) | Our own firmware command |
| Needs BOOT0 jumper | Yes | No (with *connect under reset*) | No |
| Speed (512 KB) | ~1 min @115200 | seconds | ~1 min |
| Debugging | No | Yes, until firmware disables SWJ | No |

---

## A. UART system bootloader

### Wiring

Core-board 4-pin header (`GND / TXD / RXD / +5V`) to the adapter:

| Core board | USB-UART adapter |
|---|---|
| GND | GND |
| TXD | TXD (worked for you before. If there's no response, swap with RXD) |
| RXD | RXD |
| +5V | 5V **only if the cube isn't powered by Type-C** |

- Set the adapter's logic-level jumper to **3.3 V**.
- **Unplug the Pico W** from that header first.
- Recommended: power the cube from the base's Type-C, and connect only GND/TX/RX from the adapter. LED scanning can draw more than the adapter's 5 V can supply.

### Enter / leave bootloader mode

1. Fit the **BOOT0 jumper** (2-pin header next to R7).
2. Press **Reset** (core board or base). The chip now runs the ROM bootloader. The LEDs stay dark.
3. Program or read (below).
4. Remove the BOOT0 jumper, then press Reset to run the user firmware.

The bootloader protocol is 8 data bits, **even parity**, 1 stop bit. Both tools below handle this.

### Option A1: STM32CubeProgrammer (free from st.com, needs an ST account)

GUI: select **UART**, the port (`COMx`), baud 115200, parity **Even**, then **Connect**.

CLI (`STM32_Programmer_CLI.exe` is in the install's `bin` folder):

```powershell
$cp = "C:\Program Files\STMicroelectronics\STM32Cube\STM32CubeProgrammer\bin\STM32_Programmer_CLI.exe"
& $cp -c port=COM4 br=115200                                                       # connect + chip info
& $cp -c port=COM4 br=115200 -u 0x08000000 0x80000 firmware\FullDump_512K.bin      # READ whole flash
& $cp -c port=COM4 br=115200 -ob displ                                             # show option bytes (read protection)
& $cp -c port=COM4 br=115200 -w bin\main.hex -v                                    # write + verify (erases needed sectors only)
& $cp -c port=COM4 br=115200 -g 0x08000000                                         # start the app without removing BOOT0
```

### Option A2: stm32flash (open source, single `.exe`)

Get the Windows binary from the stm32flash project on SourceForge (`stm32flash-0.7-binaries.zip`). Put `stm32flash.exe` on `PATH` or in `.agent/tools/`.

```powershell
stm32flash -b 115200 COM4                                                   # probe: should report device 0x0414 (high-density)
stm32flash -b 115200 -r firmware\FullDump_512K.bin -S 0x08000000:524288 COM5   # READ whole flash
stm32flash -b 115200 -w bin\main.hex -v -g 0x08000000 COM4                     # write + verify + run
```

> The chip reports `DBGMCU_IDCODE = 0x13090414` (DEV_ID 0x414 = high-density, 512 KB), so both tools should allow the full 512 KB range.

---

## ★ Step zero for tomorrow: full flash dump (before anything else)

1. Enter bootloader mode (above).
2. Read `0x08000000`, length `0x80000`, into `firmware/FullDump_512K_<date>.bin`. **Read it twice** and compare:
   ```powershell
   (Get-FileHash firmware\FullDump_512K_a.bin).Hash -eq (Get-FileHash firmware\FullDump_512K_b.bin).Hash
   ```
3. Show the option bytes (to rule out read protection, write protection etc.).
4. Send me the file. I'll check whether `0x08032000–0x0803D998` still holds the original app's second half: vector targets such as `SystemInit @ 0x08036768` must contain code, not `0xFF`.

---

## B. ST-Link + OpenOCD (optional, but makes debugging much nicer)

SWD needs SWDIO = **PA13**, SWCLK = **PA14**, NRST and GND. There's no SWD header, so tap the base:

| ST-Link | Cube |
|---|---|
| SWDIO | KEY6 button pad (PA13 side) |
| SWCLK | KEY5 button pad (PA14 side) |
| NRST | RESET button pad (non-GND side) |
| GND | any GND |
| 3.3V | leave unconnected (cube self-powered) or use as sense only |

Our firmware disables SWJ at boot, so either keep BOOT0 fitted (the ROM bootloader leaves SWD enabled), or use *connect under reset*:

```powershell
openocd -f interface/stlink.cfg -c "set CPUTAPID 0x2ba01477" -f target/stm32f1x.cfg `
  -c "reset_config srst_only srst_nogate connect_assert_srst" `
  -c "init; reset halt; dump_image firmware/FullDump_512K.bin 0x08000000 0x80000; reset run; shutdown"
```

TinyGo can flash through it directly with `tinygo flash -target bluepill-clone ./cmd/8x8_rgb_cube/` (it uses OpenOCD + stlink-v2).

---

## C. No-jumper workflow (planned)

Two upgrades, both optional:

1. **Software entry.** A UART console command in our firmware (`boot` or similar) jumps to the ROM bootloader: disable interrupts and SysTick, set MSP from the system-memory vector table, jump to its reset vector. Then the flash tool runs, and its `-g 0x08000000` restarts the app. There is no BOOT0 jumper and no button press.
   - The system-memory base address must be checked for GD32F103 high density (GD32F10x user manual, "Boot configuration"). On STM32F1 it is `0x1FFFF000`.
2. **Adapter-controlled reset/BOOT0.** Wire adapter **DTR → NRST** and **RTS → BOOT0** (through about 1 kΩ). stm32flash can then toggle them itself (`-i` option). This needs an adapter that brings out DTR/RTS, and two extra wires.

A `tools/flash.ps1` wrapper (build → enter bootloader → write → run → open monitor) will be added once the method is chosen.

---

## Serial console

```powershell
tinygo monitor -port COM5 -baudrate 115200
```

The VS Code "Serial Monitor" extension also works. TinyGo's `println` goes to UART1 (PA9/PA10) at 115200 8N1. Note that the bootloader uses **8E1**, the app uses **8N1**.
