# Schematics cheat sheet

Source PDFs: `docs/diagrams/`. Rendered PNGs (3x scale): `img/`.

Re-render (no extra tools needed; uses the built-in Windows PDF renderer):

```powershell
.\.agent\tools\render-pdf.ps1 -PdfPath "docs\diagrams\<file>.pdf" -OutDir "$env:TEMP\schem_img" -Scale 3 -TilesX 3 -TilesY 2
```

Tiles (`_rYcX.png`) are easier to read than full pages.

Status legend: **[seen]** = read directly from the rendered image, **[?]** = needs checking on the real board.

---

## Yellow core board (`Yellow small board schematic`): "LETOO_RGB-CUBE_CORE V0.1", 2014-12-16 [seen]

- The footprint is for an STM32F103RCT6 (LQFP64), but the chip fitted is a GD32F103RET6 (512 KB flash, 64 KB RAM).
- 8 MHz crystal (Y2). The NRST line has a 10K pull-up, 1 uF to GND and the J8 reset button.
- **BOOT0 = J2**, **BOOT1 (PB2) = J7**. Both are 3-pin jumpers: pin1 = 3V3, pin2 = BOOTx, pin3 = GND.
- **J11 "COM4"**: TX1 (PA9), RX1 (PA10), USB5V, GND. This is USART1, which the system bootloader uses. Pin order is roughly 1 = GND, 2 = TX1, 3 = RX1, 4 = 5V [?].
- Onboard LEDs: **D1 = LED1 = PB9**, **D2 = LED2 = PB8**. Each goes from 3V3 through a 1K resistor to the pin, so the LED is **active low**. D3 is the power LED.
- **24C02 EEPROM** (marked C16): SCL = PB6, SDA = PB7 (I2C1, 10K pull-ups). A0..A2 are on GND, so the address is 0x50. The original firmware's "Regedit" key/value store most likely lives here.
- **SD card (CON3)** on SDIO: D0 = PC8, D1 = PC9, D2 = PC10, D3 = PC11, CLK = PC12, CMD = PD2, card-detect = PB5. The IAP bootloader updates the firmware from `update.hex` on the SD card.
- PB4 (NJTRST) is not connected.
- PA13 net = "DD1". PA12 = "USB_DP", PA11 = "USB_DM" (names only; they go to the base connector).

### U26 board-to-base connector (yellow board side) [seen]

| Pin | Net | MCU | Pin | Net | MCU |
|---|---|---|---|---|---|
| 1 | AIN1 | PA1 | 20 | AIN0 | PA0 |
| 2 | PC3 | PC3 | 19 | PC2 | PC2 |
| 3 | PC1 | PC1 | 18 | PC0 | PC0 |
| 4 | RX2 | PA3 | 17 | TX2 | PA2 |
| 5 | DB5V | 5V in | 16 | GND | |
| 6 | RST | NRST | 15 | PB3 | PB3 |
| 7 | PA15 | PA15 | 14 | PA14 | PA14 |
| 8 | USB_DP | PA12 | 13 | DD1 | PA13 |
| 9 | 154EN | PA8 | 12 | USB_DM | PA11 |
| 10 | PC6 | PC6 | 11 | PC7 | PC7 |
| 21 | SPI1_CS | PA4 | 26 | PC5 | PC5 |
| 22 | SPI1_CLK | PA5 | 27 | GA | PB0 |
| 23 | SPI1_MISO | PA6 | 28 | GB | PB1 |
| 24 | SPI1_MOSI | PA7 | 29 | GC | PB10 |
| 25 | PC4 | PC4 | 30 | GD | PB11 |

---

## Base board `3D16RGB` (the one you said you have) [seen: page 1, page 2]

### Functions on U26

| MCU | Function | Notes |
|---|---|---|
| PC0 | KEY1 | Button to GND. **No external pull-up**, so the internal pull-up is required. Active low. |
| PC1 | KEY2 | same |
| PC2 | KEY3 | same |
| PB3 | KEY4 | Also JTDO, so JTAG must be disabled |
| **PA14** | **KEY5** | Also SWCLK, so SWJ must be fully disabled (`SWJ_CFG=100`) |
| PA13 | KEY6 | Also SWDIO, so SWJ must be fully disabled |
| PA11 | KEY7 | |
| NRST | RESET (K8) | Hardware reset button on the base |
| PA1 | SW1 | DIP switch SWR1. 4.7K external pull-up to 3V3, ON connects to GND |
| PC3 | SW2 | same |
| PA15 | SW3 | same. Also JTDI, so JTAG must be disabled. **Not in the current code yet** |
| PA12 | SW4 | same. **Not in the current code yet** |
| PC6 | IR | IRM-8601M2 receiver, 22K pull-up. NEC-style remote (probably) |
| PA0 | MUSIC_ADC | Microphone M1 into an LM358 amplifier, then to the ADC (audio-reactive modes) |
| PA2 / PA3 | TX2 / RX2 | J22 (XJK4: 5V, TX2, RX2, GND) plus J1/J3 3.5 mm "cascade" jacks (级联插座) |
| PC4 | LE | Via the 74HC245 (U27, 3.3 V to 5 V level shift) |
| PC5 | OE | Via the 74HC245. OE is active low on the MBI5024 |
| PA7 | SPI1_MOSI | Via the 74HC245 into the SDI of the first MBI5024 |
| PA5 | SPI1_CLK | Via the 74HC245 |
| PA6 | SPI1_MISO | SDO from the end of the whole MBI5024 chain (can be read back for a self-test) |
| PB0, PB1, PB10, PB11 | GA, GB, GC, GD | 74HC154 address inputs A..D |
| PA8 | 154EN | Drives both G1 and G2 of the 74HC154. **Low = enabled** |
| PC7 | 154EN2 | Not used on this sheet [?] |

### Layer (anode) drive [seen]

- The 74HC154 output `Y(n)` goes to net `G(n+1)`, then into an APM4953 dual P-MOSFET gate, and that MOSFET drives `CC(n+1)` (layer anode, 5 V). Outputs are active low, so the selected layer's P-MOSFET turns on.
- CC1..CC16 exist (16 layers). An 8-layer cube uses `D=0` (CC1..CC8).

### Column (cathode) drive [seen: page 2]

- MBI5024 constant-current sinks with 680R R-EXT, all wired as **one long daisy chain**:
  `SPI1_MOSI → U33 → U34 → … → U12 → SPI_REDOUT` (page 2, **red**, T1..T256),
  then page 3 (green, T257..T512), then page 4 (blue, T513..T768), and finally `SPI1_MISO`.
- 3 × 16 chips × 16 outputs = 768 outputs, which is 256 RGB columns, i.e. a **16 × 16 footprint**.
- ⚠ **Open question:** the current Go code shifts 24 bytes per layer (192 bits = 12 chips = 8 × 8 × 3). That matches a populated 8×8 base, not the full 3D16RGB chain shown here. On the real board: how many MBI5024 chips are fitted, and which T-pads are the 8×8 LEDs soldered to?

---

## Base board `8X8X8RGB` (newer Type-C revision) [seen: page 1 top-left]

- U26 mapping is the same as 3D16RGB, **except that pin 12 is `USB_DM` (PA11) and there is no KEY7**. Only KEY1..KEY6 exist.
- SWR1 DIP is labeled "音频模式/动画模式/测试…" (audio mode / animation mode / test …).
- IR pull-up is 10K. The 74HC245 is wired the same way (OE, LE, MOSI, CLK).

## Base board `3D8RGB`

Not read yet.
