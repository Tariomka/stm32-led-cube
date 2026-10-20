# Firmware hardware layer: first map

Source: Ghidra decompile of `firmware/FullDump_512K.bin` (ARM Cortex-M, base `0x08000000`, set up by [../tools/ghidra/Gd32Setup.java](../tools/ghidra/Gd32Setup.java)). Literal-pool constants are annotated with peripheral names in `.agent/external/re/Firmware.annotated.c`.

The firmware is written in **ST Standard Peripheral Library** style: peripherals are passed as pointers to small driver functions. Status: index + findings below. Literal-pool constants pointing to RAM are annotated too (`RAM_2000xxxx`).

## Regions

| Range | Content |
|---|---|
| `0x08000000–0x080045FF` | IAP bootloader (SDIO + FAT + FLASH + USART1). Not ported |
| `0x08006000–0x0800AFFF` | App code: engine, console, drivers glue |
| `0x0800B000–0x08012D63` | Show-data blob ([spec-blob-format.md](spec-blob-format.md)) |
| `0x08032100–0x08038FC3` | SPL drivers, startup, IRQ handlers, libc, more app code |

## Peripheral usage by function

| Area | Functions | Peripherals | Meaning (first guess 🔶) |
|---|---|---|---|
| **Display scan** | `FUN_0800a68c … FUN_0800ab90` (8 variants), `FUN_08034344 … FUN_08034728` (8 variants) | SPI1 + **DMA1 ch3** (SPI1_TX), GPIOC | Shift one layer per call via DMA. Probably 8 variants = layer or cube-type specific. GPIOC = LE/OE strobes |
| Scan timer | `FUN_08032100` (init), `FUN_0800a518` (ISR) | **TIM2** | Layer multiplex timing (see Findings) |
| IR timer | `FUN_0800ac80`, `FUN_0800af28` (reads/resets TIM3_CNT), `FUN_0800a53c` (ISR) | **TIM3** | IR pulse widths |
| SPI init | `FUN_0800af74`, `FUN_08034c14` | SPI1 | |
| Layer select | `FUN_08009df0 … FUN_0800a3b0` (9 funcs), `FUN_08035ea4/f40/fc4` | GPIOA/B/C | 74HC154 address/enable, keys, DIP |
| **Audio** | `FUN_08009ab8` | **ADC1** (DR) + **DMA1 ch1**, GPIOA | Mic on PA0, sampled by DMA |
| Console UART | `FUN_0800a598`, `FUN_0800ad08`, `FUN_08036cc0`, `FUN_08002be4` | **USART1** + DMA1 ch4 (TX) / ch5 (RX) | Debug console, PC "3D software" link |
| Cascade UART | `FUN_0800a5d8`, `FUN_0800ae14` | **USART2** + DMA1 ch7 (TX) / ch6 (RX) | 3.5 mm jacks |
| **EEPROM** | `FUN_08032dd4/e14/e3c`, `FUN_08035c80 … 08035fc4` | GPIOB IDR/BSRR (bit-bang) | Software I2C on PB6/PB7 → 24C02 (type + serial only, see Findings) |
| Keys / IR | `FUN_080322c8`, `FUN_080322e8` | EXTI | Key or IR (PC6) edge interrupts |
| AFIO | `FUN_08009bfc`, `FUN_08032458/251c/d84` | AFIO | SWJ disable + EXTI line mapping (it reads/writes MAPR) |
| Watchdog | `FUN_08032e78 … 08032eb0` | **IWDG** | The original runs a watchdog |
| Clocks | `FUN_08036768` (SystemInit), `FUN_08033758`, `FUN_0803404c …` | RCC, FLASH | 72 MHz setup |
| Generic timer init | `FUN_08036a20` | TIM1–8 | SPL `TIM_DeInit`-style switch |
| SDIO/FLASH in app | `FUN_08001bf8 …` (shared with IAP?), `FUN_08032388 …` | SDIO, FLASH | Probably linked-in IAP helpers / Regedit-in-flash? ❓ |

Strings that anchor app functions (`FUN_08036b98` = printf-like console output):
`Debug Usart Init OK!`, `Regedit format error!`, `Delete Regedit:%s OK!`, the key/speed messages, `Commands explain:`.
All app strings: `.agent/external/re/Firmware.strings.txt` (regenerate with the snippet in [README.md](README.md)).

## Findings (2026-10-18)

### Interrupts

The vector stubs (`app_vec_N`) jump through a RAM table at `0x2000042C`, index `N-12`, filled at runtime by `FUN_080347b8(index, handler)`. `Gd32Setup.java` now also seeds functions from Thumb pointers in literal pools, so these handlers are decompiled.

| Index | IRQ | Handler | Use |
|---|---|---|---|
| 0x1b | EXTI9_5 | `FUN_08009bb4` | IR receiver PC6, both edges |
| 0x20 | TIM2 | `FUN_0800a518` | Display scan. PSC 72 (1 µs), period 180/90/45 µs by cube type |
| 0x21 | TIM3 | `FUN_0800a53c` | 20 ms period, 1 µs ticks. IR interval timer (overflow = "too long") |
| 3 | SysTick? | `FUN_08036b3c` | OS tick |

The firmware runs a small cooperative OS (`OS_Timer`, tasks, `OS_Mutex*`), with periodic "timer" callbacks (`FUN_08036ae0`).

### Audio ✅

See [spec-engine.md](spec-engine.md#audio-level-firmware-fun_08007048-fun_08032584--). Implemented.

### IR remote ✅ (`FUN_08009bb4`, `FUN_08033eb4`, table in `FUN_08034ddc`)

- On each PC6 edge: interval = TIM3 count since the previous edge (0 if TIM3 overflowed, > 20 ms), then TIM3 is reset. On a falling edge the decoder gets `(previous interval = mark, current interval = space)`.
- NEC timing windows (µs): leader mark 7001–10999 with space 3501–5499 (start) or 1751–2749 (repeat, ignored); bit mark 411–709 with space 411–709 = 0, 1281–2079 = 1.
- 32 bits stored MSB first per byte; valid when byte1 = ~byte0 and byte3 = ~byte2. Code = byte0 << 8 | byte2 (address is 0x00 on the kit remote).
- Code → key action: A2→9, E2→7, 02→4, C2→5, E0→2, A8→8, 90→3, 98→6, B0→1, 68/30/18/7A/10/38/5A/42/4A/52 → 10..19. 0x22 and 0x62 do nothing. The firmware prints `IR:%x`.

Implemented: `controller.RemoteController` (decoder + table), `component.InfraredReceiver` (EXTI with `time.Now`).

### Console (USART1) and Regedit

Command records (`0x80` bytes: name[24], help[100], handler) at `0x08038868` (8 commands) plus 3 more at `0x08038160`:

| Command | Handler | Does |
|---|---|---|
| `Help` | `FUN_080095f4` | `Commands explain:` + `%s\t:%s` per command |
| `ReBoot` | `FUN_08009904` | NVIC system reset (`System Reboot!`) |
| `ReadMem [address] [offset] [count]` | `FUN_08009784` | Hex dump |
| `SetReg [key] [value]` / `GetReg [key]` / `RegList` / `DelReg [key]` / `ClrReg` | `FUN_080099ec` / `08009870` / `08009580` / `080094fc` / `08009434` | Key/value registry (keys < 12 chars, ≤ 64 entries, 2 KB), saved to internal flash (`FUN_08036f5c` → `FUN_080370e4`) |
| `SetType [type]` | `FUN_08009984` | Writes the cube type string to 24C02 address 0xA8, then reboots |
| `SetSerial [16 hex]` | `FUN_08009908` | Writes 8 bytes to 24C02 address 0xA0, then reboots |
| `Test` | `FUN_080099cc` | Parses a number, nothing else visible |

- **No firmware code reads the registry** (`FUN_080339c8` is only called by `GetReg`). It is a generic store for the PC tools.
- The PC "Download Software" sends `#?reboot?!` (every other try) and connect frames until the IAP bootloader answers, then streams `Update.hex` (`ClassUpdate.cs`). Protocol, all in `A5 05 len payload check 55` frames: PC `[17, md5[16], size u32]` → cube `[16, 0, 0,0,0,0]` (byte 1 ≠ 0 = already up to date); PC `[32, n, line u16 BE, address u32 BE, data[n]]` per type-00 hex record → cube `[33, n, line u16 BE]` within 200 ms (else the upload fails); PC `[37, 0]` twice at the end. Baud: the tool tries 921600 unless started as `"Download Software.exe" COMx <baud>`. The 3D Software's `Update.hex` is its firmware template with the show data at `0x0800B000..0x08032100` (`Compile.cs`).
- Our firmware answers this protocol itself (`controller/upload.go`, part of `TerminalController`), keeps only the show region in a flash slot and plays it in Text mode; the console's ReBoot is deferred 300 ms and cancelled by the connect frame.
- **Wire format** (`FUN_0800901c`, matches `ClassCommunication.WriteStringData` in the PC tool): text commands are `#?<command line>?!` (≤ 180 chars, a 1 s gap resets the parser). Tokens are split on spaces, the command name is compared case-insensitively, the handler gets `(argv+1, argc-1)`. Unknown commands print the Help line. Binary frames `A5 05 len data xor 55` are only logged (`RecvData=%d:`); they're for the IAP bootloader.
- Every reply line ends with `\r\n`. USART1 runs at **921600** baud in the original (`FUN_0800ad08`, confirmed by the boot log in `.agent/live_testing/original output 921600.txt`); ours stays at 115200 (the Pico uses this same console and upload protocol).
- Original boot log: IAP banner (`IAP version: 2.10`, `Start App!`), then `Debug Usart Init OK!`, `Regedit format error!`, OS task list, `WDOG Init OK!`, `System init!`, `T='8x32x8b'`, `Config1=0,Config2=0`, `DeviceType:1,DeviceNum:0`, `DisplayAdd=0x200005CC`, `Total=1`, `Run x=0,y=0,z=0,dir=0,type=7`. Key presses print `DoKeyN` followed by the action message.
- Registry flash format (`FUN_08036f84`, pages `0x08005000` and backup `0x08005800`, both erased in our dump): `u32 count`, `"Registry for 2.00"` padded to 28 bytes, `count × {key[12], u32 size}`, then the values (NUL-terminated, size rounded up to 4). Limits: key 1..11 chars, size ≤ 256, ≤ 64 entries, whole image ≤ 2 KB.

Implemented: `controller.TerminalController` (`terminal.go`: `#?…?!` / upload frame splitter; `terminal_commands.go`: commands plus our `Key`/`Action`/`Cartoon`/`Status` control commands; `registry.go`: vendor format, stored at `0x0807F000` + backup `0x0807F800`, reserved in the linker script), `component.SerialPort`, `component.WriteFlashPage`. Help texts are English (the original's are GBK Chinese). `ReadMem` refuses peripheral addresses.

### 24C02 EEPROM

Only two records: serial (8 bytes at 0xA0) and type string (8 bytes at 0xA8: `8x8x8A`, `8x8x8B`, `8x16x8A`, `161616A/B/C`, `8x32x8A/B`). At boot (`FUN_08008288`) the type selects the cube geometry and scan variant, and `FUN_08008be0` checks the serial against a digest (`Error1/Error2` on mismatch = vendor copy protection). The firmware stops (`Cartoon End`) if either check fails. Our firmware doesn't need any of it; it prints `T='<type>'` at boot and `SetType`/`SetSerial` write the records like the original.

The bus is bit-banged (`FUN_08032d84` etc.): PB6 = SCL, PB7 = SDA, open-drain outputs. Our driver: `component.EEPROM` on the hardware I2C1 peripheral (`machine.I2C0`, 100 kHz) instead of bit-banging. Our player settings live at 0x00..0x07 (`controller.SettingsController`).

### Cascade (USART2) ✅

USART2 115200 8N1 (`FUN_0800ae14`), TX PA2 → next cube's RX PA3. State at `0x200003E0`: +0 pause (`DeviceSuppend`), +1 DeviceType (1 = master), +2 last key mask, +4 DeviceNum, +8 discovery state, +0x30 u64 show clock.

- **Frame** (`FUN_08034ed0` send, `FUN_080091a0` receive): `A5 55 len payload xor`, `len` 1..15 (receiver accepts 1..16), `xor = len ^ payload bytes`.
- **Discovery** (`FUN_08036558`, at boot): wait 100 ms, send `[01]` 5× every 100 ms. A cube that receives `[01]` while discovering becomes a slave. Slave: wait 500 ms for `[02 n]` → DeviceNum = n, forward `[02 n+1]`. Master: send `[02 01]` 5× every 100 ms. Prints `DeviceType=%d,DeviceNum=%d`.
- **Display**: DeviceNum indexes the blob's display table (x, y, z, dir); `DeviceNum >= count` prints `Warn: Pos > Cubes, %d,%d` and uses 0,0,0, so with the original blob every cube shows the same thing. `dir` (0/1) selects one of two frame conversion functions (`FUN_0800a618` stores type and dir at `0x20000000`, the type table pointer at `0x20000004`; table entries: +0 scan ISR, +4 layer select, +8 dir 0 conversion, +0xC dir 1 conversion). For 8x8x8A the dir 1 conversion (`FUN_08035770`) reverses the 64 columns of each 64-voxel block (index x + 8y within a z slab), i.e. x and y mirrored = upside down. For 161616B, 161616C and 8x32x8B (this kit) the dir 1 entries (`0x08035C78`, `0x08035C7C`, `0x08035C74`) are `B.W` to their dir 0 functions. Ported as `DisplayController.SetUpsideDown`.
- **Sync** (master only, 8 ms OS timer `FUN_08036abc`, registered in `FUN_080070d8` when DeviceType = 1): unless paused, clock += 8. Keys are polled; a key counts as held after 15 consecutive polls. Sends `[20, clock u64 LE, key mask]`, then runs the key functions for rising mask bits (`FUN_0803292c`). IR codes are handled on the master only and are **not** forwarded.
- **Slave** (`FUN_08032be4`): on `[20 …]` forward it, take the clock, run key functions for rising mask bits. Slaves don't poll their own keys or IR.
- Frame waits (`FUN_08037090`) compare against the shared clock, so every cube advances frames on the master's clock.

Implemented: `controller.CascadeController` (frames, discovery, Sync, key mask) on `component.NewCascadePort`. Differences: our master builds the key mask from key interrupts (a press holds its bit for 3 frames, a key held 3 frames keeps it), so short taps work.

## Next RE steps

1. Done: display `dir` = 1 (see Cascade above).
2. Done: the 3 planes (0x600 apart) are thermometer coded (`FUN_080353f4`: level ≥ 1 → plane A, ≥ 2 → B, 3 → C) and shown in 1, 3 and 6 of 10 whole refresh sweeps (`FUN_0803460c`: slot 0 → A, 1–3 → B, 4–9 → C). Ported as 3 planes per layer weighted 1:3:6.
3. Done: IWDG watchdog (`FUN_08038038`), ported with TinyGo's `machine.Watchdog`.
