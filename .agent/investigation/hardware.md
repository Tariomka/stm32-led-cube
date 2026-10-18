# Hardware identification (from the photos in `.agent/photos/`, gitignored)

For the schematic-level pin map see [../documentation/schematics/README.md](../documentation/schematics/README.md).

## Core board ("yellow board", actually orange, silkscreen "iCubeSmart")

- GD32F103RET6 (LQFP64), 8.000 MHz crystal (Y2), AMS1117 LDO (U4), 24C02 EEPROM (C16), LEDs D1/D2 (status) and D3 (power).
- **4-pin right-angle header: `GND / TXD / RXD / +5V`.** This is USART1, the bootloader/console port. In the photo it is connected to the Raspberry Pi Pico W.
- **BOOT0**: 2-pin header with R7 (10K). Fit a jumper to boot the system (UART) bootloader. **BOOT1**: unpopulated pads with R8, so BOOT1 = 0, as required.
- **Reset** push button on the core board.
- 2×10 + 1×10 headers to the base (U26 in the schematics).
- No SD-card socket on the top side. **To check:** is there one on the back? The vendor IAP updates from SD.
- **No SWD header.** SWDIO = PA13 and SWCLK = PA14 only reach the base connector, where they are the KEY6 and KEY5 lines.

## Base board: "iCubeSmart 8x8x RGB LED CUBE" (□8x8x8 □8x8x16 □8x8x32)

- **12 × MBI5024** (you confirmed). 12 × 16 = 192 outputs = 8 × 8 × RGB, so **24 bytes per layer**. That matches the current code.
- Layer anode pads **1–32**: up to 32 layers through two 74HC154 (154EN = PA8 for layers 1–16, 154EN2 = PC7 for 17–32), plus P-MOSFETs. Only layers 1–8 are fitted on your cube.
- Power: **USB Type-C** ("8x8x8 5V 3A", "8x8x32 5V 7A" pads), large electrolytics, a "TEST LED" header.
- **Keys** (silkscreen): KEY1 *Previous*, KEY2 *Next*, KEY3 *Speed+*, KEY4 *Speed−*, KEY5 *Run/Pause*, KEY6 *Cycle*, KEY7 *ON/OFF*, plus **RESET**. "The button function can be modified in the 3D software."
- SW3 (PA15) / SW4 (PA12) are not reachable on this base and read low (ON). The original firmware reads them only to print `Config1 %d Config2 %d` at boot (`FUN_08009d9c`, bit0 = PA15 low, bit1 = PA12 low). Config1 is SW1/SW2 the same way (`FUN_08009d5c`).
- **DIP switch (SW1–SW4)**. SW1/SW2 select the effect group, applied after RESET:

  | SW1 | SW2 | Group | Function |
  |---|---|---|---|
  | 1 | 1 | Group1 | Dynamic Effects 1 (default) |
  | 0 | 1 | Group2 | Audio Effects |
  | 1 | 0 | Group3 | Test LED |
  | 0 | 0 | Group4 | Text Effects |

- **Two 3.5 mm jacks** labeled RX / TX: the UART2 "cascade" link (PA2/PA3).
- IR receiver (PC6) and microphone/LM358 (PA0) are present per the schematic. Positions in the photos are to be confirmed.

### Which schematic matches?

None exactly:
- The **pin map** of `3D16RGB` matches (KEY7 = PA11, SW3 = PA15, SW4 = PA12, IR = PC6, MUSIC = PA0).
- The **LED driver count** (12) and **32 layer pads** match the `8X8X8RGB` sheet.
- That sheet's U26 pin 12 is `USB_DM`, while your base has KEY7 there.

Treat it as an intermediate revision. Pin map: 3D16RGB. Chain length: 12 chips.

## Open hardware questions (to answer with the board in hand)

1. Does the base's **Type-C** port carry data (does a COM port appear when it's plugged into the PC)? Or is it power only?
2. Is there an SD-card slot on the back of the core board?
3. Which header pin is MCU TX vs RX? The README says straight `TXD→TXD`, `RXD→RXD`. If the bootloader doesn't answer, swap them.
4. Bit order / LED order in the 24-byte layer frame. The current `led-common-lib` layout already lights LEDs, so it's probably right. Re-verify with a "walk one LED" test.
