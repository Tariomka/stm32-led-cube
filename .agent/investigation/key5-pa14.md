# KEY5 (PA14) investigation

Status: **root cause found by static analysis, workaround implemented** (`internal/component/pin.go`: `NewInputPin` = pull-down config + `High()`, `DisableDebugPort` = SWJ_CFG 100). It still needs hardware confirmation: `scripts\build.cmd ./cmd/gpio_diag`, then `scripts\flash.cmd COM4 bin\main.hex` (passing the file skips the rebuild), then `scripts\monitor.cmd`.

## TL;DR

This is **not** a GD32 quirk and **not** a JTAG/SWD problem. It is a bug in TinyGo's STM32F1 GPIO driver. It is still present in TinyGo `0.43.0` and on the `dev` branch.

`src/machine/machine_stm32f103.go`:

```go
func (p Pin) Configure(config PinConfig) {
	...
	pos := (pin % 8) * 4
	if pin < 8 {
		port.CRL.ReplaceBits(uint32(config.Mode), 0xf, pos)   // <-- value not masked
	} else {
		port.CRH.ReplaceBits(uint32(config.Mode), 0xf, pos)
	}
	...
}
```

`src/runtime/volatile/register.go`:

```go
func (r *Register32) ReplaceBits(value uint32, mask uint32, pos uint8) {
	StoreUint32(&r.Reg, LoadUint32(&r.Reg)&^(mask<<pos)|value<<pos)   // value is NOT ANDed with mask
}
```

TinyGo encodes the pull direction in an extra "spare" bit of `PinMode`:

```go
PinInputPulldown PinMode = PinInputModePullUpDown        // 0x08 = 0b0_1000
PinInputPullup   PinMode = PinInputModePullUpDown | 0x10 // 0x18 = 0b1_1000
```

So `PinInputPullup` (`0x18`) is written as `0x18 << pos`. Bit 4 spills into **bit 0 of the next pin's 4-bit config nibble**, which is `MODE[0]` of pin N+1. That turns pin N+1 into a **10 MHz output** if it was an input, or changes its speed if it was already an output.

In `controller.NewYellowBoard`, KEY6 (**PA13**) is configured with `PinInputPullup` **after** KEY5 (**PA14**):

```text
PA14.Configure(Pullup) -> CRH[27:24] = 1000  (ok)  and PA15 MODE[0] := 1
PA13.Configure(Pullup) -> CRH[23:20] = 1000  (ok)  and PA14 MODE[0] := 1   <-- PA14 becomes 1001
```

`1001` means CNF=10, MODE=01: an *alternate-function push-pull output at 10 MHz*. With SWJ disabled, the alternate function on PA14 is just the GPIO output. The pin is driven, the internal pull-up is gone, and the button can never pull it low. That is the "MODE14 bit 0 is stuck at 1" symptom.

## Evidence

### 1. Decompiled TinyGo build (`tmp/decompiled/stm32_cube.c`)

The compiler even merged the mask: it clears 5 bits and sets bit 4 of the next nibble.

| Code (decompiled) | Meaning |
|---|---|
| `GPIOA_CRH = (CRH & 0xe0ffffff) \| 0x18000000` | PA14 pull-up, **sets bit 28 = PA15 MODE[0]** |
| `GPIOA_CRH = (CRH & 0xfe0fffff) \| 0x01800000` | PA13 pull-up, **sets bit 24 = PA14 MODE[0]** |
| `GPIOA_CRH = (CRH & 0xfffe0fff) \| 0x00018000` | PA11 pull-up, **sets bit 16 = PA12 MODE[0]** |
| `GPIOA_CRL = (CRL & 0xfffffe0f) \| 0x00000180` | PA1 pull-up, **sets bit 8 = PA2 MODE[0]** |
| `GPIOB_CRL = (CRL & 0xfffe0fff) \| 0x00018000` | PB3 pull-up, **sets bit 16 = PB4 MODE[0]** |
| `GPIOC_CRL = (CRL & 0xfffe0fff) \| 0x00018000` | PC3 pull-up, **sets bit 16 = PC4 MODE[0]** |

### 2. Register snapshots (`tmp/reg_ref/go_output.txt` vs `reference.txt`)

Every anomaly in your old register dump is explained by this bug:

| Pin | Expected | Read back (Go) | Corrupted by |
|---|---|---|---|
| PA2 | `0100` (untouched) | `0101` | PA1 (Switch1) pull-up |
| PA12 | `0100` (untouched) | `0101` | PA11 (KEY7) pull-up |
| **PA14** | `1000` | **`1001`** | **PA13 (KEY6) pull-up, configured after PA14** |
| PA15 | `0100` (untouched) | `1001` | PA14 (KEY5) pull-up |
| PB4 | `0100` (untouched) | `1001` | PB3 (KEY4) pull-up |
| PC4 | `0010` (LE output 2 MHz) | `0011` (50 MHz) | PC3 (Switch2) pull-up, configured after PC4 |
| PC1, PC2, PC3 | `1000` | `1000` | Corrupted by PC0/PC1/PC2, then fixed when the pin's own `Configure` ran later |

The C reference build (`reference.txt`) masks correctly, so it shows `1000` everywhere. That's why "it works in C".

### 3. What the bug report observations mean now

| Observation in the TinyGo issue | Explanation |
|---|---|
| Analog input → `0001` | Analog `0000`, then PA13's pull-up config adds bit 24 → `0001` |
| Manual write `0x8<<24` → `1001` | The register was printed **after** PA13 was configured in the same program (`tmp/cmd/deez_nuts`) |
| `Output2MHz` → `Output50MHz` (`0010` → `0011`) | Same leak. Output modes have no bit 4 themselves, but PA13's pull-up still sets bit 24 |
| "SWJ disable doesn't seem to change anything" | Correct. SWJ is unrelated to the readback. It still matters for **using** the pin (see below) |

Prediction: the minimal repro in `tmp/cmd/test/main.go` configures **only** PA14. It should read PA14 = `1000` (correct) and PA15 = `1001`. That is experiment E1.

## Other things found while investigating (also need fixing)

1. **SWJ still has to be disabled.** PA13 = SWDIO, PA14 = SWCLK, PA15 = JTDI, PB3 = JTDO, PB4 = NJTRST. KEY4/5/6 and SW3 sit on those pins, so the code needs `AFIO_MAPR.SWJ_CFG = 0b100`.
   - Side effect: once the firmware runs this, **SWD can no longer attach** to a running chip. To recover, use BOOT0 (system bootloader) or an ST-Link with *connect under reset*.
   - Tip for debug builds: delay the SWJ disable by about 2 s after reset, or skip it while a key is held.
2. **`AFIO_MAPR` read-modify-write hazard.** On STM32F1, `SWJ_CFG` is write-only and reads back undefined. Any later RMW of `MAPR` (e.g. `SetMAPR_CAN_REMAP`, or TinyGo's own UART/I2C/TIM remap code) can write back a different `SWJ_CFG` and re-enable JTAG.
   - On this GD32 it reads back `100`, so it is currently harmless. Still, always write `SWJ_CFG` **last**, or keep it in the written value. ST's HAL writes `SWJ_CFG=111` ("no change") when it changes other remap bits.
   - The existing `CANRemap()` is unnecessary anyway.
3. **Buttons are active low** (they short to GND; see the schematic notes), but the interrupts use `PinRising`. That fires on **release**. Use `PinFalling`, and add debouncing (mechanical bounce will trigger several EXTI interrupts).
4. **EXTI line sharing.** EXTI lines are per pin number, not per port. PB3 (KEY4) and PC3 (SW2) share EXTI3, and PA1 (SW1) and PC1 (KEY2) share EXTI1. Only one pin per line can raise interrupts, so keep the switches polled.
5. **Missing inputs.** SW3 = PA15 and SW4 = PA12 aren't configured at all. Your base has a 4-way DIP switch, and only SW1/SW2 are exposed on the front panel.

## Workaround (no TinyGo change needed)

Configure with the value that has no bit 4, then set the pull direction through ODR. `Pin.High()` writes BSRR, which sets the ODR bit, which selects pull-**up** in input-pull mode.

```go
func NewInputPin(pin machine.Pin) InputPin {
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown}) // 0x08: no bit-4 spill
	pin.High()                                                        // ODR=1 -> pull-up
	return InputPin{Pin: pin}
}
```

Alternative: write CRL/CRH directly with a masked helper. Rejected, because it duplicates driver logic.
Fragile alternative: configure pins in descending pin order. Rejected, because it still corrupts the highest pin (PA15, PB4, ...).

## Upstream fix (for the TinyGo issue / a PR)

One-line fix in `machine_stm32f103.go`:

```go
port.CRL.ReplaceBits(uint32(config.Mode)&0xf, 0xf, pos)
...
port.CRH.ReplaceBits(uint32(config.Mode)&0xf, 0xf, pos)
```

I can draft the issue update or the PR text if you want. Any other TinyGo target that passes "spare" mode bits to `ReplaceBits` should be checked as well.

## Hardware experiments (tomorrow)

Each experiment prints `GPIOx_CRL/CRH/IDR/ODR` for ports A/B/C over UART1 at boot and again after each step.

Implemented as one program, `cmd/gpio_diag` (steps 1–4 ≈ E1, E2, E4, E3). `tmp/cmd/test` and `tmp/cmd/deez_nuts` import the removed `registers` package and no longer build.

| # | Program | Expected result |
|---|---|---|
| E1 | `tmp/cmd/test` as is (PA14 pull-up only) | PA14 = `1000`, PA15 = `1001` |
| E2 | Configure PA14 then PA13 with `PinInputPullup` | PA14 = `1001` (bug reproduced) |
| E3 | Same as E2 but with the workaround helper, for every input | Every nibble clean. Pressing KEY5 → `IDR14` reads 0 and the EXTI14 callback fires |
| E4 | Workaround **without** `SWJ_CFG=100` | PA14/PA13 buttons don't work (debug port owns them). Shows SWJ disable is still needed |
