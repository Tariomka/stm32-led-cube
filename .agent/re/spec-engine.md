# Spec: show engine (player, renderers, keys, script VM)

Main source: `Simulator.dll` (x86). Decompiled with Ghidra into `.agent/external/re/Simulator.c` (gitignored, regenerate with [README.md](README.md)). The firmware has the same engine strings and structure. Hardware-specific parts are in [firmware-map.md](firmware-map.md).

Legend: ✅ read from code | 🔶 inferred, verify against the oracle or the firmware | ❓ open

## Voxel model ✅

- One buffer per display, sized L×H×W with **1 byte per voxel** = 6-bit color `0bRRGGBB` (2 bits per channel). The DLL reserves 4096 B per cube.
- Index: `(x−X0) + (y−Y0)·L + (z−Z0)·L·H`. Out-of-range writes are clipped per display.
- Palette lookup at draw time: `voxel = palette[graphValue]`, with 16 entries. Palette = named palette bytes, or the **default palette** `00 03 0C 0F 30 33 3C 3F 00 01 04 05 10 11 14 15`.
- Write mode: `Set` (0) overwrites, `Mix` (1) ORs into the existing voxel.

## Player loop (`FUN_10004ca0`) ✅

State at start: `playModel = 1` (all), `speedIdx = 4`, `displayOff = 0`, `pause = 0`, current cartoon = first of the selected group.

```text
loop forever:
  c = current cartoon
  switch c.type:
  Normal: for l in 0..loopCount-1:
            for f in 0..frameCount-1:
              clear all displays to 0
              for each element: render props active at frame f
              present frame (callback / display swap)
              wait delay/speed ms (while paused: keep waiting, re-check every 100 ms)
              if changeRequested: break out to "next selection"
  Script: run VM (init, then OnLoop) until RET_EVENT, then continue
  Audio:  deadline = now + audioTime·1000 (0 = none)
          repeat: clear; render elements with frame = audioLevel(mode); frameWait(delay)
          until deadline passed or changeRequested
          (the DLL always uses frame 0; the firmware uses the mic level, see "Audio level" below ✅)
  after the cartoon: if playModel == 1 -> advance to next (wrap to first); else replay same
  (key actions may already have set "current" and changeRequested; then that wins)
```

- **Speed table** (`speedIdx` 0…7): `0.1, 0.2, 0.5, 0.7, 1.0, 1.5, 2.0, 3.0`. Default index 4 = 1.0.
- Frame wait = `(delay_ms + 0.0) / speed`, truncated to an int ms.
- `displayOff` (key fn 9) blanks the output but the player keeps running 🔶 (the firmware decides what "off" does).

## Audio level (firmware `FUN_08007048`, `FUN_08032584`) ✅

Audio cartoon header: the `FrameCount` field is the duration in seconds (0 = forever), the `Loop` field is the level **mode**.

- ADC1 channel 0 (PA0), continuous conversion, 55.5-cycle sample time, 12-bit right aligned, DMA1 ch1 into one halfword.
- A 1 ms task stores the sample in a ring of 20 (initialised to 1050) and keeps `max` and `min` of the ring.
- `up = max < 1251 ? 0 : max < 2300 ? (max-1250)*62/1050 + 1 : 63`
- `down = min >= 850 ? 0 : min < -199 ? 63 : (850-min)*62/1050 + 1`
- `level = max(up, down)`. Mode 2: if `level < held`, `held--` and return `held` (falls one step per call), else `held = level`.
- The frame rendered is `level` (0..63). All 10 audio cartoons in our blob use mode 0, 10 s, delay 10.
- No original show calls `GetAdc`/`GetRealAdc` (checked by running all groups); their firmware mapping is not traced.

Implemented: `internal/engine/audio` (host tests), `component.Microphone` (TIM3 1 kHz ISR).

## Key functions (`FUN_10004f20`) ✅

The config byte per key (KEY1…KEY7) selects one of:

| Value | Action | Our config |
|---|---|---|
| 0 | none | |
| 1 | Replay current cartoon | |
| 2 | Previous cartoon (wraps to last) | KEY1 |
| 3 | Next cartoon (wraps to first) | KEY2 |
| 4 | Speed up (`speedIdx++` up to 7) | KEY3 |
| 5 | Speed reset (`speedIdx = 4`) | |
| 6 | Speed down (`speedIdx--` down to 0) | KEY4 |
| 7 | Toggle play model (all ↔ single loop) | KEY6 |
| 8 | Toggle pause | KEY5 |
| 9 | Toggle display off | KEY7 |
| 10…19 | Jump to cartoon n−10 of the current group | |

When a script has called `EnabledKey(1)`, key events are captured for `GetKey()` and not executed ✅.

Console messages: `Replay Cartoon`, `Previous/Next Cartoon`, `Add/Sub speed to %.2f`, `Speed reset to %.2f`, `Play Model=%d`, `Display off=%d`, `Change Cartoon %d`.

## Element renderers

### SetGraph (`FUN_100049d0` → `FUN_10005a90` → `FUN_100023a0`) ✅ / 🔶

1. For each active property: interpolate x, y, z, angle (linear, integer division).
2. Resolve the `$` names, look up the graph and palette, and decompress the graph into a scratch buffer. The DLL caches the last decoded graph.
3. `angle %= 360`, normalized to 0…359.
   - `angle == 0`: straight blit with clipping (`FUN_10001310`; it also has special paths for 90°/180°/270° 🔶 to confirm).
   - Otherwise: rotation about `axis` (0 = X, 1 = Y, 2 = Z) by `angle` degrees (`FUN_10001fc0`):
     - Uses a **quarter-wave sine table ×1000** (91 entries, 0…90°): `0,17,34,52,…,999,1000`.
     - Each slice perpendicular to the axis is rotated around its centre `(dim/2)` in fixed point (`/100000`).
     - Nearest voxel is written, with clipping to the display.
4. Write mode Set/Mix as above.

### DrawBall (`FUN_10004b10` → `FUN_10002590`) ✅

Centre (x, y, z) and radius r are in **1/100 voxel** units; r is interpolated. For each voxel centre at distance d (voxel coordinate ×100):
- `d² ≤ r²` and `d² > (r−90)²` → `colorOut`
- else `d² > (r−180)²` → `colorIn`
- else, if not hollow → `colorIn`

Colors are raw 6-bit values (no palette).

### TitleStyle (`FUN_10004bb0` → `FUN_100027f0`) ✅ / 🔶

- Inputs: font strip graph (e.g. `Demo-Text` 64×8×8), palette, index path, `high` (y offset), style, `step = frame − start`.
- The index is read as `(i16 x, i16 z)` pairs. Each step places one column of the strip at the path position, shifting along the path as `step` grows (a marquee around the cube).
- style 0 = DIR_1, 1 = DIR_2 (reverse direction in the strip).
- Height clipping is against display H.

## `$` name resolver (`FUN_100044a0`) 🔶

Grammar: `prefix '$' ('F'|'C') (op number)*`, op ∈ `+ - * / %`. The result is `prefix + itoa(value)`. The meaning of F vs C (absolute frame vs. relative to the property start) is to be confirmed in the assembly.

## Script VM (`FUN_100035e0`, dispatch table `0x1001a5a0`)

Machine model ✅:
- Registers: `A` (`DAT_1001d22c`), `B` (`DAT_1001d218`), `PC` (word pointer into the script body).
- Data stack: `SP` (`DAT_1001d230`), 200 words. Overflow prints `Error: Stack Out!` and halts.
- Ref/heap area (`DAT_1001d220`) for strings (`CREATE_REF n` reserves n words). System variables table (`DAT_10021630`, none defined).
- Execution: run from `word[0]` (init) until `RET_EVENT`, then run the OnLoop event (`word[1]`, if non-zero) until `RET_EVENT`.

Opcodes (`op << 16 | i16 arg`). The "dst" bit means `arg != 0` writes to B instead of A:

| Op | Name | Semantics |
|---|---|---|
| 0 / 7 | LDR_SYSTEM_A/B | A/B = sysvar[arg] |
| 2 / 9 | LDR_LOCAL_A/B | A/B = body[arg] (resource word) |
| 4 / 11 | LDR_LOCAL_ADD_A/B | A/B = &body[arg] (pointer to resource, e.g. string) |
| 5 / 12 | LDR_STACK_BEGIN_A/B | A/B = stackBase[arg] (script globals) |
| 6 / 13 | LDR_STACK_END_A/B | A/B = SP[arg] (arg negative = locals/params) |
| 14 / 15 / 16 | STR_SYSTEM_A / STR_STACK_BEGIN_A / STR_STACK_END_A | store A |
| 17 | SWAP_A_B | ❓ (not decompiled, trivial) |
| 18 / 19 | PUSH n / POP n | SP += n / SP −= n words |
| 20 / 21 | PUSH_A / POP_A | ❓ push A / pop into B 🔶 |
| 22 | GOTO rel | PC += rel (relative to the next instruction) |
| 23 / 24 | CALL_BEFORE / CALL_SYSTEM idx | call the system function idx (table ❓) |
| 26 | CALL_LOCAL addr | call into the script, absolute word index ❓ |
| 27 / 28 | RET_FUNCTION / RET_EVENT | return / stop the run loop |
| 29 / 30 | CREATE_REF n / DESTORY_REF n | allocate n heap words (zeroed, pointer pushed) / free |
| 32–34 | minus, not_calc (~), not_logic (!) | unary on A |
| 35 / 36 | adds / subs | ++ / −− |
| 37–46 | add sub mul div mod shl shr and or xor | `A op B` |
| 47–54 | == != && \|\| > >= < <= | boolean 0/1 |
| 55 | COPY n | copy n words from \*B to \*A |
| 56 | COPY_STR | ❓ |
| 57 | CAT | concatenate strings \*A + \*B into one of 5 rotating 128-B buffers |
| 58 / 59 | JZ / JNZ rel | if A == 0 / A != 0: PC += rel |

System functions (`CALL_SYSTEM` index, from `SystemFunctionDef.cs`). Arguments are passed on the stack (`SP[-1]`, `SP[-2]`, …); the return value goes in A 🔶.

| # | Name | Signature |
|---|---|---|
| 0 | TraceN | (str) → console print |
| 1 | String | (int) → str |
| 2 | Wait | (ms) → present frame + wait (speed-scaled 🔶) |
| 3 | Clean | (color) → fill displays |
| 4 | Point | (x, y, z, color, mix) |
| 5 | Line | (x1, y1, z1, x2, y2, z2, color, mix) → 3-D DDA (`FUN_10001000`) |
| 6 | Graph | (graph, palette, x, y, z, axis, angle, mode) |
| 7 | GetIndex | (indexName, i, type) → element (type 0 = u32, 1 = i16, 2 = u16, 3 = i8, 4 = u8) |
| 8 | GetRandomColor | () → next color from a fixed cycle (`FUN_10002a00`) |
| 9 | GetTick | () → ms |
| 10 | GetAdc | (ch) → audio level 🔶 |
| 11 | GetRealAdc | (ch) → raw ADC 🔶 |
| 12 | EnabledKey | (on) → capture keys for script |
| 13 | GetKey | () → last captured key (then clear 🔶) |

The only event is `OnLoop`.

None of the vendor blobs use scripts. The VM is needed for "all features" and for user-made shows, but it is low priority.

## Not part of the engine

`GetKeyChar` / `GetBuff` / `FUN_10002ae0` and the MD5 code are the PC app's licence/activation check (`FormAbout`). They are irrelevant to the port.

## Open items for Phase 2

- [ ] Disassemble the VM handlers Ghidra didn't make functions for (`0x10003800/3870/38b0/3900/3910/3960/3980/3f00`) and the system-function table.
- [ ] Exact `$F`/`$C` semantics.
- [ ] Unrotated blit (`FUN_10001310`) special cases. Title path/step math in detail.
- [ ] Audio cartoon: how `audioType` and the ADC spectrum feed `AudioG1/G2` (graph 8×**9**×8: 9 levels?) → firmware side.
- [ ] Cross-check each item against the ARM firmware (same strings at `0x08008xxx–0x0800Axxx`).
