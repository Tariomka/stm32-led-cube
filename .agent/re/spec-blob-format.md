# Spec: show-data blob ("CartoonDatas")

Status: ✅ **verified.** It is defined by the decompiled PC software (`_3D_Guan.Compile.Compile` / `Compress`) and confirmed by parsing the original blob with [../tools/blobdump](../tools/blobdump/main.go).

All integers are **little-endian**. "str16" is a GB2312 string, NUL-padded to 16 bytes (max 15 chars).

## Location in the original firmware

| Item | Value |
|---|---|
| Start | `0x0800B000` (`Compile.Buff_Start` = 134262784) |
| Max size | 160 000 B (`Max_Buff`), i.e. up to `0x08032100`, where the firmware's second code region starts |
| Our cube's blob | 32 100 B, `0x0800B000–0x08012D63` |
| Hex template marker | The first blob record line starts `:10B00000 65873412…` (`ParseHexHeadAndTail` cuts there) |

## Header (68 B)

| Off | Type | Field |
|---|---|---|
| 0 | u32 | magic `0x12348765` |
| 4 | u32,u32 | config: offset, length |
| 12 | u32,u32 | palettes: offset, length |
| 20 | u32,u32 | graphs: offset, length |
| 28 | u32,u32 | indexes: offset, length |
| 36…60 | 4 × (u32,u32) | cartoon groups 1–4: offset, length |

Offsets are relative to the blob start. Sections are packed in this order: config, palette, graph, index, group1–4.

## Config

| Off | Type | Field |
|---|---|---|
| 0 | u8[7] | key function for KEY1…KEY7 (see *Key functions*) |
| 7 | u8 | N = number of displays (cubes) |
| 8 + 8·i | i16 x, i16 y, i16 z, u8 dir, u8 cubeType | display i |

- `dir`: 0 = bottom, 1 = top.
- `cubeType`: 0 null, 1 = 8×8×8, 2 = 8×16×8, 3 = 16×16×16, 4 = 8×32×8 (dimensions are Long×Height×Width).
- Our cube: keys `[2,3,4,6,8,7,9]`, one display at (0,0,0), dir 0, type 1.

## Palettes (32 B each)

`str16 name, u8[16] colors`. Colors are **6-bit RGB `0bRRGGBB`**, with 2 bits (4 levels) per channel.

Example: `red = 00 30 30 …`, `red-series = 00 30 31 32 33 3B 3F`.

If an element names no palette, the engine uses the **default palette**:
`00 03 0C 0F 30 33 3C 3F 00 01 04 05 10 11 14 15` (black, blue, green, cyan, red, magenta, yellow, white, then the same at half intensity).

## Graphs (variable length)

| Off | Type | Field |
|---|---|---|
| 0 | str16 | name |
| 16 | u32 | record size = 28 + align4(dataLen) |
| 20 | u32 | dataLen |
| 24 | u8 L, u8 H, u8 W | dimensions (Long, Height, Width) |
| 27 | u8 | compression |
| 28 | u8[align4(dataLen)] | data |

- Voxel order: `index = x + y·L + z·L·H` (x < L, y < H, z < W).
- Voxel values are **palette indexes** (0–15).

### Compression (`Compress.cs`)

| Value | Name | Bits/idx | Format |
|---|---|---|---|
| 0 | NoCompress | 8 | raw `L·H·W` bytes |
| 1 / 2 / 3 | Normal_2 / _4 / _8 | 4 / 2 / 1 | `u8 n, u8 lut[n]`, then packed indexes, LSB first, `8/bits` per byte |
| 4 / 5 / 6 | Continue_2 / _4 / _8 | 4 / 2 / 1 | `u8 n, u8 lut[n]`, then RLE tokens |

- `lut` maps a packed index to the voxel value. It is the sorted set of values present.
- RLE token, with `bits` = 4/2/1 and `mask = 2^bits−1`:
  - byte `b` with bit 7 = 0: `count = b >> bits`, `value = lut[b & mask]`.
  - byte `b` with bit 7 = 1, followed by `b2`: `count = ((b & 0x7F) << (8−bits)) | (b2 >> bits)`, `value = lut[b2 & mask]`.
- Normal decoder quirk: it consumes `dataLen − n − 2` full bytes, then one last byte, stopping when the output is full.

The original blob has 79 graphs: e.g. `8x8x8`, `HELLO`, `cubeamp1..4`, `Circle1..5`, `rain1/2`, `Fly_1..4`, `Wave1/2`, `AudioG1/2` (8×9×8), `Demo-Text` (64×8×8).

## Indexes (variable length)

`str16 name, u32 len, u8[align4(len)] data`.

The engine reads the data as an array of u32/i16/u16/i8/u8, depending on the caller.

TitleStyle uses it as a **path of `(i16 x, i16 z)` pairs** (4 B per step). The text scrolls along these columns. Examples: `AllAround` (116 B = 29 steps), `TwoSide` (60 B = 15 steps).

## Cartoon groups

Each group is a sequence of cartoon records:

| Off | Type | Field |
|---|---|---|
| 0 | u32 | record size (incl. this 8-byte header) |
| 4 | u8 type, u8[3] 0 | 0 = Normal, 1 = Script, 2 = Audio |
| 8 | … | body |

Groups map to the DIP switch: group 1 = Dynamic (41 cartoons), 2 = Audio (10), 3 = Test LED (4), 4 = Text (1).

### Normal / Audio body

| Off | Type | Normal | Audio |
|---|---|---|---|
| 0 | i32 | frameCount | audioTime (seconds; 0 = forever) |
| 4 | i16 | loop count | audioType |
| 6 | u16 | delay (ms per frame, before speed scaling) | delay (ms) |
| 8 | elements… | | |

Element:

| Off | Type | Field |
|---|---|---|
| 0 | u32 | element size = 8 + n·propSize |
| 4 | u32 | kind: 0 = SetGraph (56 B props), 1 = DrawBall (20 B), 2 = TitleStyle (56 B) |
| 8 | props[n] | |

Every property starts with `i16 startFrame, i16 endFrame` (inclusive). The property is active when `start ≤ frame ≤ end`.

| Kind | Property layout after start/end |
|---|---|
| SetGraph (56) | `str16 graph`, `str16 palette`, `i16 x0,x1`, `i16 y0,y1`, `i16 z0,z1`, `i16 angle0,angle1`, `i16 axis` (0 = X, 1 = Y, 2 = Z), `i16 mode` (0 = Set, 1 = Mix/OR) |
| DrawBall (20) | `i16 x, y, z` (×100 units), `i16 r0, r1` (×100), `i16 hollow`, `u8 colorIn`, `u8 colorOut`, `u8[2] 0` |
| TitleStyle (56) | `str16 graph` (font strip), `str16 palette`, `str16 index` (path), `i16 high` (y position), `i16 style` (0 = DIR_1, 1 = DIR_2) |

- Interpolation is linear: `v = v0 + (v1 − v0)·(f − start)/(end − start)`, with integer division. At `f == start` use v0, at `f == end` use v1.
- Names may contain `$` expressions (next section).

### Script body

The body is an array of u32 words (`CompileScript.GetScriptList`):

```text
word[0]   = index of the first init instruction (= 2 + resource words)
word[1]   = index of the OnLoop event entry (0 = none)
word[2..] = resources (int constants: 1 word; strings: GB2312 + NUL, padded to whole words)
then      = init code, then functions, then events
```

Instruction = `(opcode << 16) | (u16)arg`, where arg is sign-extended as i16. Opcodes are listed in [spec-engine.md](spec-engine.md#script-vm).

None of the vendor demo blobs contain scripts.

## `$` name expressions (resource names)

Resolved by the engine (`Simulator.dll` `FUN_100044a0`, exported as `GetSourcePath`):
- `prefix$F` or `prefix$C`, optionally followed by `{+,-,*,/,%} number` and further `op number` terms.
- `F`/`C` select a frame-based counter (exact meaning to be confirmed from the disassembly. `GetSourcePath(name, frameIndex, startFrame)` suggests absolute frame vs. frame relative to start).
- The result is `prefix + decimal(value)`, e.g. `Fly_$C+1` → `Fly_1`, `Fly_2`, …
- Used only by the 8X8X32 demo, not ours.

## Tooling

```powershell
cd .agent\tools\blobdump; go run . ..\..\..\firmware\FullDump_512K.bin -v      # or any vendor *.hex
```

In agent terminals, strip `tinygo\goroot-` entries from `PATH` first.
