# Show data format

The show data is one binary block (the "cartoon data"). The vendor "3D Software" compiles it into `Update.hex` at flash address `0x0800B000`; this firmware stores it in its show slot ([storage.md](storage.md#flash-map)) and plays it in Text mode.

- All integers are **little-endian**.
- `str16` = text in GB2312 (ASCII for plain names), NUL padded to 16 bytes (max 15 characters).
- `align4(n)` = n rounded up to a multiple of 4.
- Maximum size: 160 000 bytes (`0x0800B000–0x080320FF`).

## Header (68 bytes)

| Offset | Type | Field |
|---|---|---|
| 0 | u32 | Magic `0x12348765` (bytes `65 87 34 12`) |
| 4 | u32 offset, u32 length | Config |
| 12 | u32, u32 | Palettes |
| 20 | u32, u32 | Graphs |
| 28 | u32, u32 | Indexes |
| 36, 44, 52, 60 | u32, u32 each | Cartoon groups 1–4 |

Offsets are from the start of the block. Sections follow in this order: config, palettes, graphs, indexes, groups 1–4.

## Config

| Offset | Type | Field |
|---|---|---|
| 0 | u8[7] | Key function of KEY1…KEY7 ([show-engine.md](show-engine.md#key-functions)) |
| 7 | u8 | N = number of displays (cubes) |
| 8 + 8·i | i16 x, i16 y, i16 z, u8 dir, u8 cubeType | Display *i* |

- x, y, z: position of the cube in the show's coordinate space.
- `dir`: 0 = bottom, 1 = top (mounting direction; this firmware supports 0 only).
- `cubeType`: 0 none, 1 = 8×8×8, 2 = 8×16×8, 3 = 16×16×16, 4 = 8×32×8 (Long × Height × Width).

The factory show uses keys `[2, 3, 4, 6, 8, 7, 9]` and one display at (0, 0, 0), dir 0, type 1.

## Palettes (32 bytes each)

| Offset | Type | Field |
|---|---|---|
| 0 | str16 | Name |
| 16 | u8[16] | Colors, `0bRRGGBB` |

A graph voxel value 0–15 selects a palette entry. When an element names no palette, the **default palette** applies:

```text
00 03 0C 0F 30 33 3C 3F 00 01 04 05 10 11 14 15
```

(black, blue, green, cyan, red, magenta, yellow, white, then black and the same colors at the lowest level).

## Graphs (variable size)

| Offset | Type | Field |
|---|---|---|
| 0 | str16 | Name |
| 16 | u32 | Record size = 28 + align4(dataLen) |
| 20 | u32 | dataLen |
| 24 | u8 L, u8 H, u8 W | Size: Long (x), Height (y), Width (z) |
| 27 | u8 | Compression |
| 28 | u8[align4(dataLen)] | Data |

Voxel order: `index = x + y·L + z·L·H`. Values are palette indexes 0–15.

### Compression

| Value | Name | Bits per voxel | Data |
|---|---|---|---|
| 0 | None | 8 | `L·H·W` raw bytes |
| 1 / 2 / 3 | Packed 4 / 2 / 1 | 4 / 2 / 1 | `u8 n`, `u8 lut[n]`, then packed indexes, LSB first, `8/bits` per byte |
| 4 / 5 / 6 | Run length 4 / 2 / 1 | 4 / 2 / 1 | `u8 n`, `u8 lut[n]`, then run tokens |

- `lut` maps a packed index to the voxel value (the sorted set of values used).
- Run token, with `bits` = 4, 2 or 1 and `mask = 2^bits − 1`:
  - byte `b` with bit 7 = 0: `count = b >> bits`, `value = lut[b & mask]`.
  - byte `b` with bit 7 = 1, followed by `b2`: `count = ((b & 0x7F) << (8 − bits)) | (b2 >> bits)`, `value = lut[b2 & mask]`.
- Packed decoding stops when `L·H·W` voxels are filled.

## Indexes (variable size)

| Offset | Type | Field |
|---|---|---|
| 0 | str16 | Name |
| 16 | u32 | len |
| 20 | u8[align4(len)] | Data |

Generic tables. Title elements read them as a path of `(i16 x, i16 z)` pairs, 4 bytes per step (e.g. `AllAround` = 29 steps around the cube's edge). Scripts can read them as arrays of u32, i16, u16, i8 or u8.

## Cartoon groups

Group *n* is a sequence of cartoon records:

| Offset | Type | Field |
|---|---|---|
| 0 | u32 | Record size including this 8-byte header |
| 4 | u8 type, u8[3] 0 | 0 = Normal, 1 = Script, 2 = Audio |
| 8 | … | Body |

In the factory firmware the groups belong to the DIP modes (group 1 Dynamic: 41 cartoons, group 2 Audio: 10, group 3 Test: 4, group 4 Text: 1). This firmware plays **all groups in order** of an uploaded show in Text mode.

### Normal and Audio body

| Offset | Type | Normal | Audio |
|---|---|---|---|
| 0 | i32 | Frame count | Duration in seconds (0 = until changed) |
| 4 | i16 | Loop count | Level mode |
| 6 | u16 | Delay, ms per frame at speed 1.00 | Delay |
| 8 | elements… | | |

Element:

| Offset | Type | Field |
|---|---|---|
| 0 | u32 | Element size = 8 + n·propSize |
| 4 | u32 | Kind: 0 = Graph (56-byte props), 1 = Ball (20), 2 = Title (56) |
| 8 | props[n] | |

Every property starts with `i16 startFrame, i16 endFrame` (inclusive); it draws on frames `start ≤ f ≤ end`.

| Kind | Fields after start/end |
|---|---|
| Graph | `str16 graph`, `str16 palette`, `i16 x0, x1`, `i16 y0, y1`, `i16 z0, z1`, `i16 angle0, angle1`, `i16 axis` (0 = X, 1 = Y, 2 = Z), `i16 mode` (0 = Set, 1 = Mix) |
| Ball | `i16 x, y, z` (1/100 voxel), `i16 r0, r1` (1/100 voxel), `i16 hollow`, `u8 colorIn`, `u8 colorOut`, `u8[2]` 0 |
| Title | `str16 graph` (font strip), `str16 palette`, `str16 index` (path), `i16 high` (y position), `i16 style` (0 = forward, 1 = reverse) |

Values with a 0/1 pair are interpolated linearly with integer division: `v = v0 + (v1 − v0)·(f − start)/(end − start)`; `v0` at `f = start`, `v1` at `f = end`.

### Name expressions

Graph and palette names may contain a frame expression: `prefix$F` or `prefix$C`, followed by any number of `<op><number>` terms (`+ - * / %`, applied left to right).

- `F` = frame number, `C` = frame − property start.
- The result is `prefix` + the decimal value, e.g. `Fly_$C%4+1` → `Fly_1`, `Fly_2`, `Fly_3`, `Fly_4`, `Fly_1`, …
- An invalid expression resolves to an empty name (nothing is drawn).

### Script body

An array of u32 words:

```text
word[0]   index of the first init instruction (= 2 + number of resource words)
word[1]   index of the OnLoop event entry (0 = none)
word[2..] resources: int constants (1 word) and strings (GB2312 + NUL, padded to whole words)
then      init code, functions, events
```

Instruction word = `(opcode << 16) | (u16)arg`; `arg` is a signed 16-bit value. The instruction set is in [show-engine.md](show-engine.md#script-vm).

## Factory show contents

The factory show data is 32 100 bytes with 79 graphs (e.g. `8x8x8`, `HELLO`, `cubeamp1..4`, `Circle1..5`, `rain1/2`, `Fly_1..4`, `Wave1/2`, `AudioG1/2` at 8×9×8, the 64×8×8 font strip `Demo-Text`), its palettes (e.g. `red` = `00 30 30 …`, `red-series` = `00 30 31 32 33 3B 3F`) and the index paths `AllAround` and `TwoSide`. It has no scripts.
