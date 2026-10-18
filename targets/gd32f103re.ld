/* GD32F103RET6: 512 KB flash, 64 KB SRAM. The last 4 KB (two 2 KB pages) hold the console registry. */
MEMORY
{
    FLASH_TEXT (rw) : ORIGIN = 0x08000000, LENGTH = 508K
    RAM (xrw)       : ORIGIN = 0x20000000, LENGTH = 64K
}

_stack_size = 4K;

INCLUDE "targets/arm.ld"
