//go:build windows

// Command stm32-led-cube runs the project's tasks on Windows: building, flashing and dumping the
// firmware, the serial console, show uploads and the tests. Run it from the repository root:
//
//	go run . <task> [arguments]
//
// CUBEPROG overrides the STM32CubeProgrammer CLI path; ICUBE_DOWNLOAD is the path of the vendor
// "Download Software.exe" for the upload task.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	defaultPort = "COM4"
	defaultBaud = "115200"
	firmwarePkg = "./cmd/8x8_rgb_cube"
	hexFile     = `bin\main.hex`
)

var (
	cubeProgrammer   = env("CUBEPROG", `C:\Program Files\STMicroelectronics\STM32Cube\STM32CubeProgrammer\bin\STM32_Programmer_CLI.exe`)
	downloadSoftware = os.Getenv("ICUBE_DOWNLOAD")
)

type task struct {
	name, usage, help string
	run               func(args []string) error
}

var tasks []task

func init() {
	tasks = []task{
		{"build", "[package]", `TinyGo build into bin\main.hex (default ` + firmwarePkg + `)`, build},
		{"flash", "[COMx] [file]", "Flash over the ROM bootloader (fit BOOT0 and press Reset first); builds first when no file is given", flash},
		{"dump", "[COMx]", `Read the whole 512 KB flash into firmware\dump_<time>.bin (fit BOOT0 and press Reset first)`, dump},
		{"terminal", "[-raw] [COMx] [baud] [command]", "Serial console: typed lines are sent as #?line?!; with a command, sends it once and prints the reply", terminal},
		{"monitor", "[COMx] [baud]", "Prints everything the cube sends", monitor},
		{"upload", "[COMx]", "Starts the vendor Download Software at 115200 baud for a show upload", upload},
		{"test", "[go test flags]", "Runs the host tests", test},
		{"oracle", "[go test flags]", "Runs the tests 32-bit, including the comparison with the vendor Simulator.dll", oracle},
	}
}

func main() {
	if len(os.Args) > 1 {
		for _, task := range tasks {
			if task.name == os.Args[1] {
				if err := task.run(os.Args[2:]); err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					os.Exit(1)
				}
				return
			}
		}
	}

	fmt.Println("Usage: go run . <task> [arguments]   (defaults: " + defaultPort + ", " + defaultBaud + " baud)")
	fmt.Println()
	for _, task := range tasks {
		fmt.Printf("  %-9s %-32s %s\n", task.name, task.usage, task.help)
	}
	os.Exit(2)
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

// arg returns args[index] or fallback when it is missing.
func arg(args []string, index int, fallback string) string {
	if index < len(args) {
		return args[index]
	}
	return fallback
}

// execute runs a program in the foreground with the console attached.
func execute(extraEnv []string, name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	command.Env = append(os.Environ(), extraEnv...)
	return command.Run()
}

func build(args []string) error {
	pkg := arg(args, 0, firmwarePkg)
	if err := os.MkdirAll("bin", 0o755); err != nil {
		return err
	}

	fmt.Println("Building " + pkg + " ...")
	err := execute(nil, "tinygo", "build", "-o", hexFile, "-target", `targets\gd32f103re.json`,
		"-opt", "s", "-panic", "print", "-gc", "conservative", "-size", "short", pkg)
	if err == nil {
		fmt.Println("Built " + hexFile)
	}
	return err
}

func flash(args []string) error {
	port, file := arg(args, 0, defaultPort), arg(args, 1, "")
	if file == "" {
		file = hexFile
		if err := build(nil); err != nil {
			return err
		}
	}
	return execute(nil, cubeProgrammer, "-c", "port="+port, "br=115200", "-w", file, "-v", "-g", "0x08000000")
}

func dump(args []string) error {
	out := `firmware\dump_` + time.Now().Format("20060102_150405") + ".bin"
	return execute(nil, cubeProgrammer, "-c", "port="+arg(args, 0, defaultPort), "br=115200", "-u", "0x08000000", "0x80000", out)
}

func upload(args []string) error {
	if downloadSoftware == "" {
		return errors.New(`set ICUBE_DOWNLOAD to the vendor "Download Software.exe" (in the iCubeSmart 3D Software folder)`)
	}
	if _, err := os.Stat(downloadSoftware); err != nil {
		return fmt.Errorf("vendor Download Software not found (ICUBE_DOWNLOAD): %w", err)
	}
	// Started detached: in the tool, pick the Update.hex made by "3D Software - RGB" and download.
	return exec.Command(downloadSoftware, arg(args, 0, defaultPort), "115200").Start()
}

func test(args []string) error {
	return execute(nil, "go", append(append([]string{"test"}, args...), "./tests/...")...)
}

func oracle(args []string) error {
	return execute([]string{"GOARCH=386"}, "go", append(append([]string{"test", "-count=1"}, args...), "./tests/...")...)
}

func terminal(args []string) error {
	flags := flag.NewFlagSet("terminal", flag.ContinueOnError)
	raw := flags.Bool("raw", false, "send lines as typed (with \\r\\n), without #?...?!")
	if err := flags.Parse(args); err != nil {
		return err
	}

	rest := flags.Args()
	port, err := openSerial(arg(rest, 0, defaultPort), arg(rest, 1, defaultBaud))
	if err != nil {
		return err
	}
	defer port.close()

	frame := func(line string) []byte {
		if *raw {
			return []byte(line + "\r\n")
		}
		return []byte("#?" + line + "?!")
	}
	if len(rest) > 2 {
		if err := port.write(frame(strings.Join(rest[2:], " "))); err != nil {
			return err
		}
		return port.print(2 * time.Second)
	}

	fmt.Println("Connected to " + port.name + ". Type a command (e.g. Help) and press Enter. Ctrl+C quits.")
	go func() {
		lines := bufio.NewScanner(os.Stdin)
		for lines.Scan() {
			if line := strings.TrimSpace(lines.Text()); line != "" {
				port.write(frame(line))
			}
		}
	}()
	return port.print(0)
}

func monitor(args []string) error {
	port, err := openSerial(arg(args, 0, defaultPort), arg(args, 1, defaultBaud))
	if err != nil {
		return err
	}
	defer port.close()

	fmt.Println("Monitoring " + port.name + ". Ctrl+C quits.")
	return port.print(0)
}

// Serial port through the Win32 API (no extra modules needed).

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetCommState    = kernel32.NewProc("GetCommState")
	procSetCommState    = kernel32.NewProc("SetCommState")
	procSetCommTimeouts = kernel32.NewProc("SetCommTimeouts")
)

// dcb is the Win32 DCB structure.
type dcb struct {
	length    uint32
	baudRate  uint32
	flags     uint32 // bit 0 fBinary; DTR, RTS and flow control bits 0 = off
	reserved  uint16
	xonLim    uint16
	xoffLim   uint16
	byteSize  byte
	parity    byte
	stopBits  byte
	xonChar   byte
	xoffChar  byte
	errorChar byte
	eofChar   byte
	evtChar   byte
	reserved1 uint16
}

type commTimeouts struct {
	readInterval, readTotalMultiplier, readTotalConstant uint32
	writeTotalMultiplier, writeTotalConstant             uint32
}

type serialPort struct {
	name   string
	handle syscall.Handle
}

// openSerial opens a COM port at baud 8N1 with DTR and RTS off; reads return at once.
func openSerial(name, baud string) (*serialPort, error) {
	rate, err := strconv.Atoi(baud)
	if err != nil {
		return nil, fmt.Errorf("baud rate %q: %w", baud, err)
	}

	path, err := syscall.UTF16PtrFromString(`\\.\` + name)
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	port := &serialPort{name: name + " at " + baud, handle: handle}

	state := dcb{length: uint32(unsafe.Sizeof(dcb{}))}
	if ok, _, err := procGetCommState.Call(uintptr(handle), uintptr(unsafe.Pointer(&state))); ok == 0 {
		port.close()
		return nil, fmt.Errorf("GetCommState %s: %w", name, err)
	}
	state.baudRate, state.flags = uint32(rate), 1
	state.byteSize, state.parity, state.stopBits = 8, 0, 0
	if ok, _, err := procSetCommState.Call(uintptr(handle), uintptr(unsafe.Pointer(&state))); ok == 0 {
		port.close()
		return nil, fmt.Errorf("SetCommState %s: %w", name, err)
	}
	timeouts := commTimeouts{readInterval: 0xFFFFFFFF, writeTotalConstant: 1000}
	if ok, _, err := procSetCommTimeouts.Call(uintptr(handle), uintptr(unsafe.Pointer(&timeouts))); ok == 0 {
		port.close()
		return nil, fmt.Errorf("SetCommTimeouts %s: %w", name, err)
	}
	return port, nil
}

func (this *serialPort) close() { syscall.CloseHandle(this.handle) }

func (this *serialPort) write(data []byte) error {
	var written uint32
	if err := syscall.WriteFile(this.handle, data, &written, nil); err != nil {
		return err
	}
	if int(written) != len(data) {
		return errors.New("serial write timed out")
	}
	return nil
}

// print copies everything received to the console for duration (0 = until the program ends).
func (this *serialPort) print(duration time.Duration) error {
	buffer := make([]byte, 256)
	end := time.Now().Add(duration)
	for duration == 0 || time.Now().Before(end) {
		var read uint32
		if err := syscall.ReadFile(this.handle, buffer, &read, nil); err != nil {
			return err
		}
		if read > 0 {
			os.Stdout.Write(buffer[:read])
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}
