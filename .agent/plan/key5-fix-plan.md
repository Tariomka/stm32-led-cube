# Plan: fix KEY5 (and the rest of the input pins)

Status: **done, confirmed on hardware.** Only step 8 (upstream issue/PR) is left.

Background: [../investigation/key5-pa14.md](../investigation/key5-pa14.md)

## Steps

1. **Diagnostic program** (`tmp/cmd/gpio_diag/main.go`, throwaway):
   - Prints a decoded nibble table for GPIOA/B/C (CRL, CRH, IDR, ODR), plus MAPR and EXTI.
   - Runs experiments E1–E4, selected by a constant or by DIP SW1/SW2 at boot.
2. **Run E1–E4 on hardware.** Confirm the bit-4 leak and that `SWJ_CFG=100` is still required.
3. **Fix in `internal/component/pin.go`:** `NewInputPin` configures `PinInputPulldown`, then calls `pin.High()` for the pull-up. One place, so every button and switch benefits.
4. **Clean up `internal/component/registers/swj_config.go`:**
   - Keep only `DisableJTAG` (`SWJ_CFG=100`), written last.
   - Drop `CANRemap`, and drop the AFIO clock disable (EXTI needs AFIO anyway).
   - Optional: in debug builds, delay the SWJ disable so SWD can still attach.
5. **Buttons:**
   - Change to `PinFalling` (active low).
   - Add debouncing: ignore edges within about 30 ms of the last accepted one, using `time.Now()` / ticks inside the callback, or poll in a 1 ms task.
   - Move the `println` calls out of the ISR (set flags instead).
6. **Add SW3 (PA15) and SW4 (PA12)** as polled inputs. Read SW1/SW2 at boot for the group selection, like the original.
7. **Verify:**
   - All seven keys and four switches print the right state.
   - Register dump is clean: no `1001` / `0101` on unconfigured pins.
   - LE pin PC4 is back to `0010`.
8. **Upstream:** update your TinyGo issue with the root cause and the one-line `&0xf` patch. Optionally open a PR (I can draft it).

## Done when

KEY5 toggles Run/Pause reliably, and none of the 11 inputs bounce or leak into neighboring pins.
