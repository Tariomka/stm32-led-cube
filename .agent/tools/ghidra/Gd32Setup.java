// Prepares a raw GD32F103 flash image (loaded at 0x08000000) for analysis:
// memory blocks for SRAM/peripherals/SCS, peripheral labels, functions from vector tables,
// and the show-data blob marked as data. Usage: -preScript Gd32Setup.java
// @category Setup
import ghidra.app.cmd.disassemble.DisassembleCommand;
import ghidra.app.script.GhidraScript;
import ghidra.program.model.address.Address;
import ghidra.program.model.address.AddressSet;
import ghidra.program.model.mem.Memory;
import ghidra.program.model.symbol.SourceType;

public class Gd32Setup extends GhidraScript {
    private static final long APP_VECTORS = 0x08006000L;
    private static final long IAP_VECTORS = 0x08000000L;
    private static final long BLOB_START = 0x0800B000L;
    private static final long BLOB_END = 0x08032100L;

    private static final Object[][] PERIPHERALS = {
        {"TIM2", 0x40000000L}, {"TIM3", 0x40000400L}, {"TIM4", 0x40000800L}, {"TIM5", 0x40000C00L},
        {"TIM6", 0x40001000L}, {"TIM7", 0x40001400L}, {"RTC", 0x40002800L}, {"WWDG", 0x40002C00L},
        {"IWDG", 0x40003000L}, {"SPI2", 0x40003800L}, {"SPI3", 0x40003C00L}, {"USART2", 0x40004400L},
        {"USART3", 0x40004800L}, {"UART4", 0x40004C00L}, {"UART5", 0x40005000L}, {"I2C1", 0x40005400L},
        {"I2C2", 0x40005800L}, {"USB", 0x40005C00L}, {"CAN", 0x40006400L}, {"BKP", 0x40006C00L},
        {"PWR", 0x40007000L}, {"DAC", 0x40007400L}, {"AFIO", 0x40010000L}, {"EXTI", 0x40010400L},
        {"GPIOA", 0x40010800L}, {"GPIOB", 0x40010C00L}, {"GPIOC", 0x40011000L}, {"GPIOD", 0x40011400L},
        {"GPIOE", 0x40011800L}, {"ADC1", 0x40012400L}, {"ADC2", 0x40012800L}, {"TIM1", 0x40012C00L},
        {"SPI1", 0x40013000L}, {"TIM8", 0x40013400L}, {"USART1", 0x40013800L}, {"ADC3", 0x40013C00L},
        {"SDIO", 0x40018000L}, {"DMA1", 0x40020000L}, {"DMA2", 0x40020400L}, {"RCC", 0x40021000L},
        {"FLASH_IF", 0x40022000L}, {"CRC", 0x40023000L},
        {"SysTick", 0xE000E010L}, {"NVIC_ISER", 0xE000E100L}, {"NVIC_IPR", 0xE000E400L}, {"SCB", 0xE000ED00L},
        {"DBGMCU", 0xE0042000L},
    };

    @Override
    protected void run() throws Exception {
        Memory mem = currentProgram.getMemory();
        createBlock(mem, "SRAM", 0x20000000L, 0x10000);
        createBlock(mem, "PERIPH_APB1", 0x40000000L, 0x8000);
        createBlock(mem, "PERIPH_APB2", 0x40010000L, 0x4000);
        createBlock(mem, "PERIPH_AHB", 0x40018000L, 0xC000);
        createBlock(mem, "SCS", 0xE000E000L, 0x1000);
        createBlock(mem, "DBG", 0xE0042000L, 0x100);
        createBlock(mem, "SYSMEM", 0x1FFFF000L, 0x810);

        for (Object[] p : PERIPHERALS) {
            createLabel(toAddr((Long) p[1]), (String) p[0], true, SourceType.USER_DEFINED);
        }

        createLabel(toAddr(BLOB_START), "SHOW_DATA_BLOB", true, SourceType.USER_DEFINED);

        seedVectors(IAP_VECTORS, 76, "iap_vec_");
        seedVectors(APP_VECTORS, 76, "app_vec_");
        seedPointers(0x08006000L, BLOB_START);
        seedPointers(BLOB_END, 0x08039000L);
    }

    // Handlers registered at runtime (IRQ table in RAM, VM opcode tables) are only referenced
    // as Thumb pointers in literal pools.
    private void seedPointers(long start, long end) throws Exception {
        for (long p = start; p + 4 <= end; p += 4) {
            long v = getInt(toAddr(p)) & 0xFFFFFFFFL;
            if ((v & 1) == 0 || v == 0xFFFFFFFFL || !inCode(v & ~1L)) {
                continue;
            }
            Address a = toAddr(v & ~1L);
            if (getFunctionContaining(a) != null) {
                continue;
            }
            new DisassembleCommand(new AddressSet(a), null, true).applyTo(currentProgram, monitor);
            createFunction(a, null);
        }
    }

    private static boolean inCode(long a) {
        return (a >= 0x08006000L && a < BLOB_START) || (a >= BLOB_END && a < 0x08039000L);
    }

    private void createBlock(Memory mem, String name, long start, long size) throws Exception {
        Address a = toAddr(start);
        if (mem.getBlock(a) == null) {
            mem.createUninitializedBlock(name, a, size, false);
        }
    }

    private void seedVectors(long table, int entries, String prefix) throws Exception {
        for (int i = 1; i < entries; i++) {
            long v = getInt(toAddr(table + 4L * i)) & 0xFFFFFFFFL;
            if ((v & 1) == 0 || v < 0x08000000L || v >= 0x08080000L) {
                continue;
            }
            long target = v & ~1L;
            if (target >= BLOB_START && target < BLOB_END) {
                continue;
            }
            Address a = toAddr(target);
            new DisassembleCommand(new AddressSet(a), null, true).applyTo(currentProgram, monitor);
            if (getFunctionAt(a) == null) {
                createFunction(a, prefix + i);
            }
        }
    }
}
