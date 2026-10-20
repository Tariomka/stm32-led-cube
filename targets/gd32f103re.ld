/* GD32F103RET6: 512 KB flash, 64 KB SRAM. Flash above 350 KB is data (internal/component/flash.go):
   0x08057800 uploaded show (158 KB), 0x0807F000 console registry + backup (2 x 2 KB). */
MEMORY
{
    FLASH_TEXT (rw) : ORIGIN = 0x08000000, LENGTH = 350K
    RAM (xrw)       : ORIGIN = 0x20000000, LENGTH = 64K
}

_stack_size = 4K;

INCLUDE "targets/arm.ld"
