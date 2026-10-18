package player

import (
	"strconv"

	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
)

// Script VM opcodes (op << 16 | int16 argument), see .agent/re/spec-engine.md.
const (
	opLdrSystemA     = 0
	opLdrLocalA      = 2
	opLdrLocalAddA   = 4
	opLdrStackBeginA = 5
	opLdrStackEndA   = 6
	opLdrSystemB     = 7
	opLdrLocalB      = 9
	opLdrLocalAddB   = 11
	opLdrStackBeginB = 12
	opLdrStackEndB   = 13
	opStrSystemA     = 14
	opStrStackBeginA = 15
	opStrStackEndA   = 16
	opSwapAB         = 17
	opPush           = 18
	opPop            = 19
	opPushA          = 20
	opPopA           = 21
	opGoto           = 22
	opCallSystem     = 24
	opCallLocal      = 26
	opRetFunction    = 27
	opRetEvent       = 28
	opCreateRef      = 29
	opDestroyRef     = 30
	opMinus          = 32
	opNot            = 33
	opLogicalNot     = 34
	opInc            = 35
	opDec            = 36
	opAdd            = 37
	opSub            = 38
	opMul            = 39
	opDiv            = 40
	opMod            = 41
	opShl            = 42
	opShr            = 43
	opAnd            = 44
	opOr             = 45
	opXor            = 46
	opEq             = 47
	opNe             = 48
	opLogicalAnd     = 49
	opLogicalOr      = 50
	opGt             = 51
	opGe             = 52
	opLt             = 53
	opLe             = 54
	opCopy           = 55
	opCopyStr        = 56
	opCat            = 57
	opJz             = 58
	opJnz            = 59
)

// Script system functions (CALL_SYSTEM index)
const (
	sysTraceN = iota
	sysString
	sysWait
	sysClean
	sysPoint
	sysLine
	sysGraph
	sysGetIndex
	sysGetRandomColor
	sysGetTick
	sysGetAdc
	sysGetRealAdc
	sysEnabledKey
	sysGetKey
)

// Virtual memory layout. Pointers held in A/B are byte addresses in this space.
const (
	stackWords   = 200
	heapWords    = 1024
	tempBuffers  = 5
	tempSize     = 128
	maxSteps     = 20000 // instructions per Step before yielding
	stackBase    = 0
	heapBase     = stackBase + 4*stackWords
	tempBase     = heapBase + 4*heapWords
	bodyBase     = tempBase + tempBuffers*tempSize
	sysVarsCount = 0
)

type vm struct {
	p       *Player
	mem     []byte
	bodyLen int
	a, b    int32
	pc      int // byte address
	sp      int // byte address
	heap    int
	temp    int
	phase   int // 0 = init, 1 = OnLoop, 2 = done
	onLoop  int
	failed  bool
}

func newVM(p *Player, body string) *vm {
	v := &vm{p: p, mem: make([]byte, bodyBase+len(body)), bodyLen: len(body)}
	copy(v.mem[bodyBase:], body)
	v.sp, v.heap = stackBase, heapBase
	if len(body) >= 8 {
		v.pc = bodyBase + 4*int(blob.U32(body, 0))
		v.onLoop = int(blob.U32(body, 4))
	} else {
		v.phase = 2
	}
	return v
}

// stepScript runs the script until it waits, ends or uses up its instruction budget.
func (p *Player) stepScript(c blob.Cartoon, now uint32) (rendered bool, waitMs int, done bool) {
	v := p.vm
	if v == nil {
		return false, 0, true
	}
	for i := 0; i < maxSteps; i++ {
		if v.phase == 2 || v.failed {
			return false, 0, true
		}
		if wait, yielded := v.exec(now); yielded {
			return true, p.frameWait(wait), false
		}
	}
	return false, 0, false
}

func (v *vm) word(addr int) int32 {
	if addr < 0 || addr+4 > len(v.mem) {
		v.fail("bad address")
		return 0
	}
	return int32(uint32(v.mem[addr]) | uint32(v.mem[addr+1])<<8 | uint32(v.mem[addr+2])<<16 | uint32(v.mem[addr+3])<<24)
}

func (v *vm) setWord(addr int, val int32) {
	if addr < 0 || addr+4 > len(v.mem) {
		v.fail("bad address")
		return
	}
	v.mem[addr], v.mem[addr+1], v.mem[addr+2], v.mem[addr+3] = byte(val), byte(val>>8), byte(val>>16), byte(val>>24)
}

func (v *vm) fail(msg string) {
	if !v.failed {
		v.failed = true
		v.p.log("Script error: " + msg)
	}
}

func (v *vm) str(addr int32) string {
	start := int(addr)
	if start <= 0 || start >= len(v.mem) {
		return ""
	}
	end := start
	for end < len(v.mem) && v.mem[end] != 0 {
		end++
	}
	return string(v.mem[start:end])
}

func (v *vm) setStr(addr int, s string) {
	if addr <= 0 || addr+len(s)+1 > len(v.mem) {
		v.fail("bad string address")
		return
	}
	copy(v.mem[addr:], s)
	v.mem[addr+len(s)] = 0
}

// tempString stores s in one of the rotating temporary buffers and returns its address.
func (v *vm) tempString(s string) int32 {
	if len(s) >= tempSize {
		s = s[:tempSize-1]
	}
	addr := tempBase + v.temp*tempSize
	v.temp = (v.temp + 1) % tempBuffers
	v.setStr(addr, s)
	return int32(addr)
}

func (v *vm) arg(n int) int32 { return v.word(v.sp - 4*n) }

func (v *vm) stackOut() {
	v.fail("Stack Out!")
}

func (v *vm) dst(useB bool, val int32) {
	if useB {
		v.b = val
	} else {
		v.a = val
	}
}

func b2i(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// exec runs one instruction; yielded reports a Wait call with its duration.
func (v *vm) exec(now uint32) (wait int, yielded bool) {
	ins := uint32(v.word(v.pc))
	if v.failed {
		return 0, false
	}
	v.pc += 4
	op, arg := int(ins>>16), int(int16(ins))
	toB := arg != 0
	switch op {
	case opLdrSystemA, opLdrSystemB:
		v.dst(op == opLdrSystemB, 0)
	case opLdrLocalA, opLdrLocalB:
		v.dst(op == opLdrLocalB, v.word(bodyBase+4*arg))
	case opLdrLocalAddA, opLdrLocalAddB:
		v.dst(op == opLdrLocalAddB, int32(bodyBase+4*arg))
	case opLdrStackBeginA, opLdrStackBeginB:
		v.dst(op == opLdrStackBeginB, v.word(stackBase+4*arg))
	case opLdrStackEndA, opLdrStackEndB:
		v.dst(op == opLdrStackEndB, v.word(v.sp+4*arg))
	case opStrSystemA:
	case opStrStackBeginA:
		v.setWord(stackBase+4*arg, v.a)
	case opStrStackEndA:
		v.setWord(v.sp+4*arg, v.a)
	case opSwapAB:
		v.a, v.b = v.b, v.a
	case opPush:
		v.sp += 4 * arg
		if v.sp > stackBase+4*stackWords {
			v.stackOut()
		}
	case opPop:
		v.sp -= 4 * arg
	case opPushA:
		if v.sp >= stackBase+4*stackWords-4 {
			v.stackOut()
			break
		}
		v.setWord(v.sp, v.a)
		v.sp += 4
	case opPopA:
		v.sp -= 4
		v.a = v.word(v.sp)
	case opGoto:
		v.pc += 4 * arg
	case opCallSystem:
		return v.system(arg, now)
	case opCallLocal:
		if v.sp >= stackBase+4*stackWords {
			v.stackOut()
			break
		}
		v.setWord(v.sp, int32(v.pc))
		v.sp += 4
		v.pc = bodyBase + 4*arg
	case opRetFunction:
		v.sp -= 4
		v.pc = int(v.word(v.sp))
	case opRetEvent:
		if v.phase == 0 && v.onLoop != 0 {
			v.phase, v.pc = 1, bodyBase+4*v.onLoop
		} else {
			v.phase = 2
		}
	case opCreateRef:
		if v.heap+4*arg > tempBase {
			v.fail("out of string memory")
			break
		}
		for i := 0; i < arg; i++ {
			v.setWord(v.heap+4*i, 0)
		}
		v.setWord(v.sp, int32(v.heap))
		v.sp += 4
		v.heap += 4 * arg
	case opDestroyRef:
		v.heap -= 4 * arg
	case opMinus:
		v.dst(toB, -v.a)
	case opNot:
		v.dst(toB, ^v.a)
	case opLogicalNot:
		v.dst(toB, b2i(v.a == 0))
	case opInc:
		v.dst(toB, v.a+1)
	case opDec:
		v.dst(toB, v.a-1)
	case opAdd:
		v.dst(toB, v.a+v.b)
	case opSub:
		v.dst(toB, v.a-v.b)
	case opMul:
		v.dst(toB, v.a*v.b)
	case opDiv, opMod:
		if v.b == 0 {
			v.fail("division by zero")
			break
		}
		if op == opDiv {
			v.dst(toB, v.a/v.b)
		} else {
			v.dst(toB, v.a%v.b)
		}
	case opShl:
		v.dst(toB, v.a<<(uint32(v.b)&31))
	case opShr:
		v.dst(toB, v.a>>(uint32(v.b)&31))
	case opAnd:
		v.dst(toB, v.a&v.b)
	case opOr:
		v.dst(toB, v.a|v.b)
	case opXor:
		v.dst(toB, v.a^v.b)
	case opEq:
		v.dst(toB, b2i(v.a == v.b))
	case opNe:
		v.dst(toB, b2i(v.a != v.b))
	case opLogicalAnd:
		v.dst(toB, b2i(v.a != 0 && v.b != 0))
	case opLogicalOr:
		v.dst(toB, b2i(v.a != 0 || v.b != 0))
	case opGt:
		v.dst(toB, b2i(v.a > v.b))
	case opGe:
		v.dst(toB, b2i(v.a >= v.b))
	case opLt:
		v.dst(toB, b2i(v.a < v.b))
	case opLe:
		v.dst(toB, b2i(v.a <= v.b))
	case opCopy:
		for i := 0; i < arg; i++ {
			v.setWord(int(v.a)+4*i, v.word(int(v.b)+4*i))
		}
	case opCopyStr:
		v.setStr(int(v.a), v.str(v.b))
	case opCat:
		s := v.tempString(v.str(v.a) + v.str(v.b))
		v.dst(toB, s)
	case opJz:
		if v.a == 0 {
			v.pc += 4 * arg
		}
	case opJnz:
		if v.a != 0 {
			v.pc += 4 * arg
		}
	}
	return 0, false
}

func (v *vm) system(fn int, now uint32) (int, bool) {
	p := v.p
	switch fn {
	case sysTraceN:
		p.log(v.str(v.arg(1)))
	case sysString:
		v.a = v.tempString(strconv.Itoa(int(v.arg(1))))
	case sysWait:
		return int(v.arg(1)), true
	case sysClean:
		p.canvas.Clear(byte(v.arg(1)))
	case sysPoint:
		p.canvas.Point(int(v.arg(1)), int(v.arg(2)), int(v.arg(3)), byte(v.arg(4)), int(v.arg(5)))
	case sysLine:
		p.canvas.Line(int(v.arg(1)), int(v.arg(2)), int(v.arg(3)), int(v.arg(4)), int(v.arg(5)), int(v.arg(6)),
			byte(v.arg(7)), int(v.arg(8)))
	case sysGraph:
		g, ok := p.assets.graph(v.str(v.arg(1)))
		if !ok {
			break
		}
		pal, ok := p.assets.palette(v.str(v.arg(2)))
		if !ok {
			break
		}
		p.canvas.Graph(g, pal, int(v.arg(3)), int(v.arg(4)), int(v.arg(5)), int(v.arg(6)), int(v.arg(7)), int(v.arg(8)))
	case sysGetIndex:
		v.a = 0
		if idx, ok := p.assets.index(v.str(v.arg(1))); ok {
			v.a = int32(blob.IndexValue(idx, int(v.arg(2)), int(v.arg(3))))
		}
	case sysGetRandomColor:
		v.a = int32(p.randomColor())
	case sysGetTick:
		v.a = int32(now)
	case sysGetAdc, sysGetRealAdc:
		v.a = 0
		if p.ADC != nil {
			v.a = int32(p.ADC(int(v.arg(1)), fn == sysGetRealAdc))
		}
	case sysEnabledKey:
		p.capturedKey = 0
		p.keyCapture = v.arg(1) != 0
	case sysGetKey:
		v.a = int32(p.capturedKey)
		p.capturedKey = 0
	}
	return 0, false
}

// randomColor steps through the vendor's fixed 31-color cycle.
func (p *Player) randomColor() byte {
	p.colorStep++
	if p.colorStep > 30 {
		p.colorStep = 0
	}
	u := p.colorStep + 1
	if u < 29 {
		if u > 23 {
			u = p.colorStep + 3
		}
	} else {
		u = p.colorStep + 4
	}
	r := u % 4
	switch u / 4 {
	case 0:
		return byte(r)
	case 1:
		return byte(r<<2 | 3)
	case 2:
		return byte(r<<4 | 0xF)
	case 3:
		return byte((3-r)*4) | 0x33
	case 4:
		return byte(3-r) | 0x30
	case 5:
		if r != 3 {
			return byte(^(r << 4) & 0x30)
		}
		return 4
	case 6:
		return byte(r << 2)
	case 7:
		return byte(r<<4 | 0xC)
	case 8:
		return byte((r*-0x10 - 0xD | (3-r)*4) & 0x3C)
	}
	return 0
}
