package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const blobBase = 0x0800B000

var opNames = map[uint16]string{
	0: "LDR_SYSTEM_A", 1: "LDR_GLOBAL_A", 2: "LDR_LOCAL_A", 3: "LDR_GLOBAL_ADD_A", 4: "LDR_LOCAL_ADD_A",
	5: "LDR_STACK_BEGIN_A", 6: "LDR_STACK_END_A", 7: "LDR_SYSTEM_B", 8: "LDR_GLOBAL_B", 9: "LDR_LOCAL_B",
	10: "LDR_GLOBAL_ADD_B", 11: "LDR_LOCAL_ADD_B", 12: "LDR_STACK_BEGIN_B", 13: "LDR_STACK_END_B",
	14: "STR_SYSTEM_A", 15: "STR_STACK_BEGIN_A", 16: "STR_STACK_END_A", 17: "SWAP_A_B", 18: "PUSH", 19: "POP",
	20: "PUSH_A", 21: "POP_A", 22: "GOTO", 23: "CALL_BEFORE", 24: "CALL_SYSTEM", 25: "CALL_GLOBAL", 26: "CALL_LOCAL",
	27: "RET_FUNCTION", 28: "RET_EVENT", 29: "CREATE_REF", 30: "DESTORY_REF", 32: "minus", 33: "not_calc",
	34: "not_logic", 35: "adds", 36: "subs", 37: "add", 38: "sub", 39: "mul", 40: "div", 41: "mod", 42: "shift_l",
	43: "shift_r", 44: "and_calc", 45: "or_calc", 46: "xor", 47: "equal", 48: "not_equal", 49: "and_logic",
	50: "or_logic", 51: "big_then", 52: "big_e_then", 53: "small_then", 54: "small_e_then", 55: "COPY",
	56: "COPY_STR", 57: "CAT", 58: "JZ", 59: "JNZ",
}

var sysFuncs = []string{"TraceN", "String", "Wait", "Clean", "Point", "Line", "Graph", "GetIndex",
	"GetRandomColor", "GetTick", "GetAdc", "GetRealAdc", "EnabledKey", "GetKey"}

func loadImage(path string) []byte {
	if strings.HasSuffix(strings.ToLower(path), ".bin") {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		return b
	}
	img := make([]byte, 0x80000)
	for i := range img {
		img[i] = 0xFF
	}
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var base uint32
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if len(l) < 11 || l[0] != ':' {
			continue
		}
		n, _ := strconv.ParseUint(l[1:3], 16, 8)
		a, _ := strconv.ParseUint(l[3:7], 16, 16)
		t, _ := strconv.ParseUint(l[7:9], 16, 8)
		switch t {
		case 4:
			v, _ := strconv.ParseUint(l[9:13], 16, 16)
			base = uint32(v) << 16
		case 0:
			for k := 0; k < int(n) && 9+2*k+2 <= len(l); k++ {
				v, _ := strconv.ParseUint(l[9+2*k:11+2*k], 16, 8)
				off := int(base+uint32(a)) - 0x08000000 + k
				if off >= 0 && off < len(img) {
					img[off] = byte(v)
				}
			}
		}
	}
	return img
}

func str16(b []byte) string {
	e := 0
	for e < 16 && b[e] != 0 {
		e++
	}
	return string(b[:e])
}

func u32(b []byte, o int) uint32 { return binary.LittleEndian.Uint32(b[o:]) }
func i16(b []byte, o int) int16  { return int16(binary.LittleEndian.Uint16(b[o:])) }

func main() {
	img := loadImage(os.Args[1])
	verbose := len(os.Args) > 2 && os.Args[2] == "-v"
	blob := img[blobBase-0x08000000:]
	fmt.Printf("magic 0x%08X\n", u32(blob, 0))
	names := []string{"config", "palette", "graph", "index", "group1", "group2", "group3", "group4"}
	var off, ln [8]int
	for i := range names {
		off[i] = int(u32(blob, 4+8*i))
		ln[i] = int(u32(blob, 8+8*i))
		fmt.Printf("%-8s off=0x%05X len=%d\n", names[i], off[i], ln[i])
	}
	end := off[7] + ln[7]
	fmt.Printf("blob size %d (end 0x%08X)\n", end, blobBase+end)

	c := blob[off[0]:]
	fmt.Printf("\nkeys: %v\n", c[:7])
	for k := 0; k < int(c[7]); k++ {
		o := 8 + 8*k
		fmt.Printf("display %d: x=%d y=%d z=%d dir=%d type=%d\n", k, i16(c, o), i16(c, o+2), i16(c, o+4), c[o+6], c[o+7])
	}

	fmt.Printf("\npalettes (%d):\n", ln[1]/32)
	for o := off[1]; o < off[1]+ln[1]; o += 32 {
		fmt.Printf("  %-16s % X\n", str16(blob[o:]), blob[o+16:o+32])
	}

	fmt.Println("\ngraphs:")
	ng := 0
	compNames := []string{"None", "N2", "N4", "N8", "C2", "C4", "C8"}
	for o := off[2]; o < off[2]+ln[2]; {
		total := int(u32(blob, o+16))
		dl := int(u32(blob, o+20))
		L, H, W, ct := blob[o+24], blob[o+25], blob[o+26], blob[o+27]
		cn := "?"
		if int(ct) < len(compNames) {
			cn = compNames[ct]
		}
		if verbose {
			fmt.Printf("  %-16s %2dx%2dx%2d %-4s data=%d\n", str16(blob[o:]), L, H, W, cn, dl)
		}
		ng++
		if total <= 0 {
			break
		}
		o += total
	}
	fmt.Printf("  count=%d\n", ng)

	fmt.Println("\nindexes:")
	for o := off[3]; o < off[3]+ln[3]; {
		l := int(u32(blob, o+16))
		fmt.Printf("  %-16s len=%d % X\n", str16(blob[o:]), l, blob[o+20:o+20+min(l, 24)])
		o += 20 + ((l + 3) &^ 3)
	}

	for g := 0; g < 4; g++ {
		fmt.Printf("\ngroup %d:\n", g+1)
		n := 0
		for o := off[4+g]; o < off[4+g]+ln[4+g]; {
			size := int(u32(blob, o))
			typ := blob[o+4]
			body := blob[o+8 : o+size]
			n++
			switch typ {
			case 0, 2:
				frameOrTime := int32(u32(body, 0))
				loopOrType, delay := i16(body, 4), i16(body, 6)
				var els []string
				for e := 8; e < len(body); {
					es := int(u32(body, e))
					et := u32(body, e+4)
					cnt := 0
					switch et {
					case 0, 2:
						cnt = (es - 8) / 56
					case 1:
						cnt = (es - 8) / 20
					}
					desc := fmt.Sprintf("%s×%d", []string{"Graph", "Ball", "Title"}[et], cnt)
					if verbose && (et == 0 || et == 2) {
						var gn []string
						for p := 0; p < cnt; p++ {
							q := e + 8 + 56*p
							s := fmt.Sprintf("[%d-%d]%s/%s", i16(body, q), i16(body, q+2), str16(body[q+4:]), str16(body[q+20:]))
							if et == 0 {
								s += fmt.Sprintf(" x%d>%d y%d>%d z%d>%d a%d>%d axis%d mode%d",
									i16(body, q+36), i16(body, q+38), i16(body, q+40), i16(body, q+42), i16(body, q+44),
									i16(body, q+46), i16(body, q+48), i16(body, q+50), i16(body, q+52), i16(body, q+54))
							}
							gn = append(gn, s)
						}
						desc += "{" + strings.Join(gn, " ") + "}"
					}
					els = append(els, desc)
					if es <= 0 {
						break
					}
					e += es
				}
				kind := "Normal frames"
				if typ == 2 {
					kind = "Audio time"
				}
				fmt.Printf("  #%-2d %-6s %s=%d loop/type=%d delay=%d  %s\n", n, map[byte]string{0: "Normal", 2: "Audio"}[typ], kind, frameOrTime, loopOrType, delay, strings.Join(els, ", "))
			case 1:
				words := len(body) / 4
				fmt.Printf("  #%-2d Script words=%d init@%d OnLoop@%d\n", n, words, u32(body, 0), u32(body, 4))
				if verbose {
					disasm(body)
				}
			default:
				fmt.Printf("  #%-2d type=%d size=%d\n", n, typ, size)
			}
			o += size
		}
	}
}

func disasm(body []byte) {
	words := len(body) / 4
	initAt := int(u32(body, 0))
	for w := 2; w < initAt && w < words; w++ {
		v := u32(body, 4*w)
		fmt.Printf("      %04X: res 0x%08X %q\n", w, v, strings.Trim(string(body[4*w:4*w+4]), "\x00"))
	}
	for w := initAt; w < words; w++ {
		v := u32(body, 4*w)
		op := uint16(v >> 16)
		arg := int16(v)
		name := opNames[op]
		extra := ""
		if op == 24 && int(arg) < len(sysFuncs) {
			extra = " ; " + sysFuncs[arg]
		}
		fmt.Printf("      %04X: %-18s %d%s\n", w, name, arg, extra)
	}
}
