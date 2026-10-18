# MCU cheat sheet: GD32F103RET6 (STM32F103-compatible)

Sources: `docs/manuals/` (RM0008 for STM32F1, GD32F10x User Manual v2.3, GD32F103xx datasheet v2.7). This is a quick reference only. Use the manuals for anything timing-critical.

## Memory map

| Region | Address | Notes |
|---|---|---|
| Flash | `0x08000000`, 512 KB | GD32 high density: **2 KB pages** |
| SRAM | `0x20000000`, 64 KB | TinyGo `bluepill` target only uses 20 KB |
| System memory (ROM bootloader) | STM32F1: `0x1FFFF000`, **GD32: verify** | |
| Flash size register | `0x1FFFF7E0` (KB, 16-bit) | |
| Unique ID | `0x1FFFF7E8` | TinyGo `deviceIDAddr` |
| DBGMCU_IDCODE | `0xE0042000` | Reads `0x13090414` (DEV_ID 0x414 = high density) |

Peripheral bases: `AFIO 0x40010000`, `EXTI 0x40010400`, `GPIOA 0x40010800`, `GPIOB 0x40010C00`, `GPIOC 0x40011000`, `SPI1 0x40013000`, `USART1 0x40013800`, `RCC 0x40021000` (`APB2ENR` at +0x18), `FLASH 0x40022000`.

## GPIO configuration (CRL = pins 0–7, CRH = pins 8–15; 4 bits per pin = `CNF[1:0] MODE[1:0]`)

| MODE | Meaning | | CNF (input, MODE=00) | CNF (output, MODE≠00) |
|---|---|---|---|---|
| 00 | Input | | 00 analog | 00 GP push-pull |
| 01 | Output 10 MHz | | 01 floating (reset) | 01 GP open-drain |
| 10 | Output 2 MHz | | 10 pull-up/down (ODR bit: 1 = up, 0 = down) | 10 AF push-pull |
| 11 | Output 50 MHz | | 11 reserved | 11 AF open-drain |

Common nibbles: `0100` floating input (reset default), `1000` input pull, `0010` GP out 2 MHz, `1011` AF out 50 MHz (SPI/UART TX).

TinyGo `PinMode` values for F1: `PinOutput = 2`, `PinInput = 0` (**analog!**), `PinInputModeFloating = 4`, `PinInputPulldown = 0x08`, `PinInputPullup = 0x18` (⚠ bit 4 leaks, see [../investigation/key5-pa14.md](../investigation/key5-pa14.md)).

`BSRR` (+0x10) sets/resets ODR bits atomically. `IDR` (+0x08) is the input level.

## AFIO_MAPR (`0x40010004`): SWJ_CFG [26:24]

| SWJ_CFG | PA13 | PA14 | PA15 | PB3 | PB4 |
|---|---|---|---|---|---|
| 000 (reset) | SWDIO | SWCLK | JTDI | JTDO | NJTRST |
| 001 | SWDIO | SWCLK | JTDI | JTDO | GPIO |
| 010 | SWDIO | SWCLK | GPIO | GPIO | GPIO |
| **100** | **GPIO** | **GPIO** | GPIO | GPIO | GPIO |

- The AFIO clock (`RCC_APB2ENR.AFIOEN`) must be on to write it.
- STM32: the field is write-only, so any read-modify-write of MAPR can rewrite it. GD32: reads back. Write `SWJ_CFG` last.
- After `100`, SWD is dead until reset. Recover with BOOT0 or connect-under-reset.

## EXTI

- Line *n* = pin number *n* of **one** port, selected in `AFIO_EXTICR1..4`.
- IRQs: `EXTI0..4` have their own vectors. `EXTI9_5` and `EXTI15_10` are shared. TinyGo handles that via `Pin.SetInterrupt`.
- Line conflicts on this board: PC1 (KEY2) / PA1 (SW1) on line 1, and PB3 (KEY4) / PC3 (SW2) on line 3. Only the KEYs use EXTI. The switches are polled.

## Display chain

### MBI5024 (12 in a chain, constant-current sinks for the LED cathodes)

- `SDI` is shifted on the **CLK rising** edge and passed on from `SDO` to the next chip. Max 25 MHz. The current code sends LSB-first (`LSBFirst: true`).
- `LE` high copies the shift register into the output latch (transparent while high, holds on the falling edge).
- `OE` **low** enables outputs. Pulse `OE` high while switching layers (blanking) to avoid ghosting.
- 192 bits per layer = 24 bytes. The first byte shifted ends up in the **last** chip.
- Driven through a 74HC245 (3.3 V to 5 V): `PC4 → LE`, `PC5 → OE`, `PA5 → CLK`, `PA7 → SDI`. `PA6` (SPI1_MISO) reads back the chain's `SDO` (self-test possible).

### 74HC154 (4-to-16 decoder, anode/layer select)

- `A..D = PB0, PB1, PB10, PB11`. `G1 = G2 = PA8` (154EN, **active low**). Second decoder `154EN2 = PC7` for layers 17–32.
- The selected output goes low, turning on the P-MOSFET, which powers layer *n* at 5 V.

### Scan sequence per layer (target)

```text
OE=1 (blank) -> shift 24 bytes -> LE pulse -> set A..D to layer -> 154EN=0 -> OE=0 -> wait (layer time)
```

At 8 layers × 100 Hz × 2^k BAM bit-planes:
- k = 4 gives a 1.25 ms frame slice per layer.
- SPI at 9 MHz (TinyGo F1 max divider = /8 → 9 MHz) takes about 21 µs per 24 bytes.

This is easily feasible in a timer ISR. Use TIM2 or TIM3, because **TIM4 is TinyGo's tick timer**.

## Clocks (TinyGo runtime)

HSE 8 MHz → PLL ×9 → 72 MHz SYSCLK. APB1 = 36 MHz (timers ×2 = 72 MHz), APB2 = 72 MHz, ADC = 12 MHz.

## Other peripherals on the board

| Function | Pins | Peripheral |
|---|---|---|
| Console / bootloader | PA9 TX, PA10 RX | USART1 |
| Cascade jacks | PA2 TX, PA3 RX | USART2 |
| EEPROM 24C02 | PB6 SCL, PB7 SDA | I2C1 @ 0x50 |
| Microphone | PA0 | ADC1 ch0 |
| IR receiver | PC6 | GPIO/EXTI6 or TIM3_CH1 (full remap) or TIM8_CH1 |
| SD card (if fitted) | PC8–PC12, PD2, CD = PB5 | SDIO |
| Status LEDs | PB9 (D1), PB8 (D2) | GPIO, active low |
