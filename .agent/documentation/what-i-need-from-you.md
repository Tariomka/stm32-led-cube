# What I need from you (tomorrow)

## Hardware to bring / have on the desk

- [ ] The cube (base + core board) and a **Type-C cable** for power.
- [ ] **USB-UART adapter**, 3.3 V logic (CH340 / CP2102 / FT232). Bonus if it has **DTR/RTS** pins (for the no-jumper workflow).
- [ ] A **jumper cap** for BOOT0 (or a female-female wire).
- [ ] A few female-female Dupont wires.
- [ ] *If you have one:* **ST-Link V2** (clone is fine). Gives fast dumps, real debugging, and no BOOT0 juggling.
- [ ] *If you have one:* a cheap **logic analyzer** (24 MHz / 8 ch "Saleae clone"). Lets us capture the original scan timing.
- [ ] *If you have it:* the iCubeSmart **PC "3D software"**, or a link to it.

## Software (I can install these myself if you say OK)

| What | Why | How |
|---|---|---|
| USB-UART driver | COM port for the adapter | Usually automatic. CH340: WCH driver |
| **stm32flash** or **STM32CubeProgrammer** | UART bootloader read/write | stm32flash: single exe (open source). CubeProgrammer: installer from st.com (needs an ST account) |
| Ghidra + JDK 21 | Reverse-engineer the vendor app | `scoop bucket add extras; scoop bucket add java; scoop install temurin21-jdk ghidra` |
| (optional) GNU make | Keep using the `Makefile` | `scoop install make` (the Makefile needs Windows fixes anyway) |

## Order of operations (about 1 hour)

1. **Unplug the Pico W** from the core board's GND/TXD/RXD/+5V header, and plug the USB-UART adapter in instead. Tell me which **COM port** appears.
2. **Full flash dump** (BOOT0 jumper + Reset, then read 512 KB twice + option bytes). **Nothing gets written before this.**
3. I analyze the dump (a few minutes). Then either:
   - flash the complete original and record its behavior (console log, videos per DIP group, key behavior), or
   - skip that if the second half is gone.
4. KEY5 experiments E1–E4 (each is a build + flash + read the console).
5. Quick hardware checks:
   - Does the base's Type-C show up as a COM port?
   - Is there an SD slot on the back of the core board?
   - "Walk one LED" test to confirm the bit/LED order.

## Decisions I need from you (any time)

1. **Keep the vendor IAP bootloader** at `0x08000000`, with our app at `0x08006000`? Or replace everything with our app from `0x08000000` (simpler, as now)?
2. **Public repo vs vendor data:** commit the extracted animation data and fonts, or keep them local/generated?
3. **Scope:** which original features matter to you? All four DIP groups? IR remote? PC "3D software" protocol compatibility? Cascade?
4. Is it OK for me to touch `.vscode/settings.json` (`go.terminal.activateEnvironment: false`) and replace or fix the `Makefile` for Windows?
