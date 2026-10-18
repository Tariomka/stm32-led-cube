# Firmware hardware layer: first map

Source: Ghidra decompile of `firmware/FullDump_512K.bin` (ARM Cortex-M, base `0x08000000`, set up by [../tools/ghidra/Gd32Setup.java](../tools/ghidra/Gd32Setup.java)). Literal-pool constants are annotated with peripheral names in `.agent/external/re/Firmware.annotated.c`.

The firmware is written in **ST Standard Peripheral Library** style: peripherals are passed as pointers to small driver functions. Status: **index only.** The detailed reading is the next RE step.

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
| Scan timer | `FUN_0800ac80`, `FUN_0800af28` (reads TIM3_CNT), `FUN_08032100` | **TIM3** | Layer multiplex / BAM timing interrupt |
| SPI init | `FUN_0800af74`, `FUN_08034c14` | SPI1 | |
| Layer select | `FUN_08009df0 … FUN_0800a3b0` (9 funcs), `FUN_08035ea4/f40/fc4` | GPIOA/B/C | 74HC154 address/enable, keys, DIP |
| **Audio** | `FUN_08009ab8` | **ADC1** (DR) + **DMA1 ch1**, GPIOA | Mic on PA0, sampled by DMA |
| Console UART | `FUN_0800a598`, `FUN_0800ad08`, `FUN_08036cc0`, `FUN_08002be4` | **USART1** + DMA1 ch4 (TX) / ch5 (RX) | Debug console, PC "3D software" link |
| Cascade UART | `FUN_0800a5d8`, `FUN_0800ae14` | **USART2** + DMA1 ch7 (TX) / ch6 (RX) | 3.5 mm jacks |
| **EEPROM** | `FUN_08032dd4/e14/e3c`, `FUN_08035c80 … 08035fc4` | GPIOB IDR/BSRR (bit-bang) | Software I2C on PB6/PB7 → 24C02 "Regedit" |
| Keys / IR | `FUN_080322c8`, `FUN_080322e8` | EXTI | Key or IR (PC6) edge interrupts |
| AFIO | `FUN_08009bfc`, `FUN_08032458/251c/d84` | AFIO | SWJ disable + EXTI line mapping (it reads/writes MAPR) |
| Watchdog | `FUN_08032e78 … 08032eb0` | **IWDG** | The original runs a watchdog |
| Clocks | `FUN_08036768` (SystemInit), `FUN_08033758`, `FUN_0803404c …` | RCC, FLASH | 72 MHz setup |
| Generic timer init | `FUN_08036a20` | TIM1–8 | SPL `TIM_DeInit`-style switch |
| SDIO/FLASH in app | `FUN_08001bf8 …` (shared with IAP?), `FUN_08032388 …` | SDIO, FLASH | Probably linked-in IAP helpers / Regedit-in-flash? ❓ |

Strings that anchor app functions (`FUN_08036b98` = printf-like console output):
`Debug Usart Init OK!`, `Regedit format error!`, `Delete Regedit:%s OK!`, the key/speed messages, `Commands explain:`.

## Next RE steps

1. TIM3 IRQ handler (vector 45 at `0x08006000+4·45`) → exact scan sequence: layer order, BAM bit-planes for 2-bit color, timings, OE/LE pulse widths.
2. Frame-buffer → SPI byte layout (bit order of the 192 bits per layer, R/G/B chip order). Compare with our current `led-common-lib` layout.
3. Key scanning (polling vs EXTI, debounce), DIP read at boot (group select), `displayOff` handling.
4. ADC/DMA rate + the audio algorithm feeding Audio cartoons.
5. IR decoder (timer capture or EXTI + TIM) + remote code table.
6. Console command table + Regedit format (24C02 layout).
7. USART2 cascade protocol.
