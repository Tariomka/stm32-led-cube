# Show engine

The firmware plays shows with a port of the factory show engine. It produces the same frames as the factory engine for the same show data.

## Model

- **Canvas**: one 8×8×8 voxel buffer per display (cube), 1 byte per voxel = `0bRRGGBB`. Index `(x − X0) + (y − Y0)·L + (z − Z0)·L·H`, where (X0, Y0, Z0) is the display's position from the show's display table. Writes outside a display are clipped.
- **Show**: cartoons + key function table (KEY1…KEY7).
- **Cartoon** kinds: Normal (keyframed elements), Audio (frame chosen by the microphone level), Script (a small program).
- A show comes either from show data ([show-data.md](show-data.md)) or from Go code (the built-in shows, [below](#built-in-shows)). Both render identically.

## Player

Defaults at power-up (before saved settings are applied): play model = all, speed index 4 (1.00), not paused, display on, first cartoon.

```text
loop:
  start the current cartoon
  Normal: repeat loopCount times:
            for frame in 0..frameCount-1:
              clear canvas, draw every element's properties active on this frame
              show it, wait delay / speed
  Audio:  until the duration (seconds, 0 = forever) has passed:
            clear canvas, draw the elements at frame = audio level (0..63)
            show it, wait delay / speed
  Script: run the program; each Wait(ms) shows the canvas and waits ms / speed
  at the end: play model "all" -> next cartoon (wrapping), "single" -> same cartoon again
  a key that selects a cartoon ends the current one at once
```

- **Speed** table (index 0–7): 0.1, 0.2, 0.5, 0.7, **1.0**, 1.5, 2.0, 3.0. Frame wait = `delay × 10 / (speed × 10)` ms, integer.
- **Pause** stops the player; the last frame stays on the cube.
- **Display off** blanks the LEDs; the player keeps running.
- A cartoon with nothing to draw (0 frames, 0 loops) is skipped.
- All timing uses a millisecond show clock; in a cascade the slaves use the master's clock ([serial-protocol.md](serial-protocol.md#cascade-uart2)).

## Key functions

Each key's function code comes from the show's key table. The same codes are used by the IR remote and the console `Action` command.

| Code | Function | Log line |
|---|---|---|
| 0 | none | |
| 1 | Replay the current cartoon | `Replay Cartoon` |
| 2 | Previous cartoon (wraps to the last) | `Previous Cartoon` |
| 3 | Next cartoon (wraps to the first) | `Next Cartoon` |
| 4 | Speed up (max index 7) | `Add speed to %.2f` |
| 5 | Speed reset (index 4) | `Speed reset to %.2f` |
| 6 | Speed down (min index 0) | `Sub speed to %.2f` |
| 7 | Toggle play model (all ↔ single) | `Play Model=%d` |
| 8 | Toggle pause | `DeviceSuppend:%d` |
| 9 | Toggle display off | `Display off=%d` |
| 10–19 | Play cartoon 0–9 (ignored if the show has fewer) | `Change Cartoon %d` |

Built-in shows use keys `[2, 3, 4, 6, 8, 7, 9]`: KEY1 previous, KEY2 next, KEY3 speed up, KEY4 speed down, KEY5 pause, KEY6 play model, KEY7 display off.

When a script has called `EnabledKey(1)`, key functions are not executed; the code is kept for the script's `GetKey()`.

## Renderers

### Graph

Draws a graph (voxel picture of palette indexes) at a position, optionally rotated.

1. Interpolate x, y, z and angle for the frame.
2. Resolve name expressions, look up the graph and the palette (default palette if none).
3. `angle` is taken modulo 360.
   - 0°: copy the graph's whole box with its minimum corner at (x, y, z), clipped to the display.
   - Otherwise: rotate about the axis (0 = X, 1 = Y, 2 = Z) through the graph's centre. Each line of voxels along the axis moves as a whole to its rotated position in the perpendicular plane, computed with a quarter-wave sine table scaled by 1000 (91 entries, 0–90°) and fixed-point math (`/100000`); positions outside the display are skipped.
4. Each voxel's color is `palette[index]`. Mode **Set** (0) overwrites the canvas voxel, so index 0 (normally black) erases what is behind the graph's box; mode **Mix** (1) ORs the color in, so black is transparent.

### Ball

Centre (x, y, z) and radius r are in 1/100 voxel; r is interpolated. For each voxel with squared distance d² from the centre (in the same units):

| Condition | Color |
|---|---|
| `d² ≤ r²` and `d² > (r − 90)²` | `colorOut` (shell) |
| else `d² > (r − 180)²` | `colorIn` |
| else, if not hollow | `colorIn` |

Colors are raw `0bRRGGBB` values (no palette).

### Title

Scrolls a font strip (a wide graph, e.g. 64×8×8) along a path:

- The path is an index of `(x, z)` column positions.
- Each step (`frame − start`) shifts the strip one column along the path, so text runs around the cube like a marquee.
- `high` is the y position; style 0 runs forward, 1 in reverse.
- Height is clipped to the display.

### Element frame quirk

Inside one element, after a property that interpolates (a frame strictly between its start and end), the following properties of the same element see the frame **shifted by that property's start**. This matches the factory engine and some factory cartoons depend on it.

## Audio level

Audio cartoons pick their frame from the microphone loudness (0–63).

- ADC1 channel 0 (PA0) converts continuously (55.5-cycle sample time, 12 bit).
- A 1 kHz timer interrupt (TIM3) stores the latest sample in a ring of 20 (initialised to 1050) and tracks the ring's max and min.
- `up = max < 1251 ? 0 : max < 2300 ? (max − 1250)·62/1050 + 1 : 63`
- `down = min ≥ 850 ? 0 : min < −199 ? 63 : (850 − min)·62/1050 + 1`
- `level = max(up, down)`.
- The cartoon's level mode (the Loop field): mode 2 makes the level fall slowly: when the new level is lower than the held one, the held level drops by 1 per frame; otherwise it jumps up. Other modes use the level directly.

The factory audio cartoons all use mode 0, 10 s duration and 10 ms delay.

## Script VM

Scripts are compiled programs ([show-data.md](show-data.md#script-body)). The factory shows contain none, but uploaded shows may.

Machine:

- Registers A and B (32-bit), program counter.
- Stack of 200 words (`Stack Out!` on overflow), heap of 1024 words for strings, 5 rotating 128-byte buffers for string results.
- Execution starts at `word[0]` (init) and runs until `RET_EVENT`, then runs the OnLoop event (`word[1]`) until `RET_EVENT` again, then the cartoon ends.
- The player runs at most 20 000 instructions per step, so a script without `Wait` can't freeze the cube.
- Errors (bad address, division by zero, out of memory) print `Script error: <reason>` and end the cartoon.

Instruction = `opcode << 16 | i16 arg`. For unary and binary operators `arg ≠ 0` writes the result to B instead of A.

| Opcode | Name | Effect |
|---|---|---|
| 0 / 7 | LDR_SYSTEM_A / B | A/B = system variable (none defined, 0) |
| 2 / 9 | LDR_LOCAL_A / B | A/B = resource word `arg` |
| 4 / 11 | LDR_LOCAL_ADD_A / B | A/B = address of resource `arg` (strings) |
| 5 / 12 | LDR_STACK_BEGIN_A / B | A/B = stack word `arg` from the bottom (globals) |
| 6 / 13 | LDR_STACK_END_A / B | A/B = word at SP + `arg` (negative = locals, parameters) |
| 14 | STR_SYSTEM_A | no effect |
| 15 / 16 | STR_STACK_BEGIN_A / STR_STACK_END_A | store A like the loads |
| 17 | SWAP_A_B | swap A and B |
| 18 / 19 | PUSH n / POP n | SP ± n words |
| 20 / 21 | PUSH_A / POP_A | push A / pop into A |
| 22 | GOTO rel | PC += rel words (from the next instruction) |
| 24 | CALL_SYSTEM idx | system function (table below) |
| 26 | CALL_LOCAL addr | push the return address, jump to word `addr` |
| 27 | RET_FUNCTION | pop the return address |
| 28 | RET_EVENT | end of init / event |
| 29 / 30 | CREATE_REF n / DESTROY_REF n | allocate n zeroed heap words (address pushed) / free them |
| 32–36 | minus, bit not, logical not, ++, −− | on A |
| 37–46 | + − × ÷ % << >> & \| ^ | A op B |
| 47–54 | == != && \|\| > >= < <= | 0 / 1 |
| 55 | COPY n | copy n words from \*B to \*A |
| 56 | COPY_STR | copy the string at B to A |
| 57 | CAT | concatenate the strings at A and B into a temporary buffer |
| 58 / 59 | JZ / JNZ rel | jump if A == 0 / A ≠ 0 |

System functions (arguments on the stack, first argument at SP − 1 word; result in A):

| # | Name | Arguments → result |
|---|---|---|
| 0 | TraceN | (string) → prints a log line |
| 1 | String | (int) → string |
| 2 | Wait | (ms) → shows the canvas, waits ms / speed |
| 3 | Clean | (color) → fills the canvas |
| 4 | Point | (x, y, z, color, mix) |
| 5 | Line | (x1, y1, z1, x2, y2, z2, color, mix) → 3-D line |
| 6 | Graph | (graph, palette, x, y, z, axis, angle, mode) |
| 7 | GetIndex | (index, i, type) → element *i* (type 0 = u32, 1 = i16, 2 = u16, 3 = i8, 4 = u8) |
| 8 | GetRandomColor | () → next color of a fixed 31-color cycle |
| 9 | GetTick | () → show clock, ms |
| 10 | GetAdc | (channel) → audio level 0–63 |
| 11 | GetRealAdc | (channel) → raw 12-bit microphone sample |
| 12 | EnabledKey | (on) → capture keys for the script |
| 13 | GetKey | () → last captured key function code, then cleared |

## Built-in shows

The Dynamic, Audio and Test modes play built-in shows: the factory cartoons written as Go code (`internal/show`), so they need no show data.

- A cartoon is a function returning a keyframe description: frame count, loop count, delay and elements. Each element is a list of **Graph** props (start/end frame, shape, palette, from/to pose x, y, z, angle, axis, mode) and **Ball** props (start/end, centre, radius from/to, hollow, colors).
- Shapes are ASCII art: one line per layer from the top (y = H−1) down; each line holds W rows (z = 0, 1, …) separated by spaces, each with L voxels (x = 0, 1, …). A voxel is `.` (palette index 0) or a hex digit 1–F.
- Palettes are 16 colors (`0bRRGGBB`), by name.
- The Go cartoons are generated once from the factory show data with `go run ./tools/showgen`; a test checks that every frame equals the factory rendering.

| Mode | Cartoons |
|---|---|
| Dynamic | 41 animations |
| Audio | 10 microphone cartoons |
| Test | 4 LED test patterns |

The factory Text group (one scrolling-text cartoon) is not built in; Text mode plays uploaded show data instead.
