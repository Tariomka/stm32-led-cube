param(
    [string] $Image = "firmware\FullDump_512K.bin",
    [string] $In = ".agent\external\re\Firmware.c",
    [string] $Out = ".agent\external\re\Firmware.annotated.c"
)

# Appends "/* DAT_x=0x40010C00 GPIOB+0x0 */" to lines whose literal-pool constants point at peripherals.
$names = @{
    '40000000' = 'TIM2'; '40000400' = 'TIM3'; '40000800' = 'TIM4'; '40000C00' = 'TIM5'; '40001000' = 'TIM6'; '40001400' = 'TIM7'
    '40003000' = 'IWDG'; '40003800' = 'SPI2'; '40003C00' = 'SPI3'; '40004400' = 'USART2'; '40004800' = 'USART3'
    '40005400' = 'I2C1'; '40006C00' = 'BKP'; '40007000' = 'PWR'; '40010000' = 'AFIO'; '40010400' = 'EXTI'
    '40010800' = 'GPIOA'; '40010C00' = 'GPIOB'; '40011000' = 'GPIOC'; '40011400' = 'GPIOD'; '40012400' = 'ADC1'
    '40012C00' = 'TIM1'; '40013000' = 'SPI1'; '40013400' = 'TIM8'; '40013800' = 'USART1'; '40018000' = 'SDIO'
    '40020000' = 'DMA1'; '40020400' = 'DMA2'; '40021000' = 'RCC'; '40022000' = 'FLASH'
    'E000E000' = 'SCS'; 'E000E100' = 'NVIC'; 'E000ED00' = 'SCB'
}
# PowerShell parses 0xE000xxxx literals as negative Int32, so use unsigned values.
$scsStart = [uint32]3758153728; $nvic = [uint32]3758153984; $scb = [uint32]3758157056; $scsEnd = [uint32]3758157824
$img = [IO.File]::ReadAllBytes((Resolve-Path $Image))
$rx = [regex]'\b(?:PTR_DAT|_DAT|DAT)_(080[0-3][0-9a-f]{4})\b'
$sb = New-Object Text.StringBuilder
foreach ($line in [IO.File]::ReadLines((Resolve-Path $In))) {
    $ann = @()
    foreach ($m in $rx.Matches($line)) {
        $off = [Convert]::ToInt32($m.Groups[1].Value, 16) - 0x08000000
        if ($off -lt 0 -or $off + 4 -gt $img.Length) { continue }
        $v = [BitConverter]::ToUInt32($img, $off)
        if (-not (($v -ge 0x40000000 -and $v -lt 0x40030000) -or ($v -ge $scsStart -and $v -lt $scsEnd))) { continue }
        $base = [uint32]($v -band [uint32]4294966272)
        if ($v -ge $scsStart) { $base = if ($v -ge $scb) { $scb } elseif ($v -ge $nvic) { $nvic } else { $scsStart } }
        $key = '{0:X8}' -f $base
        $n = if ($names.ContainsKey($key)) { $names[$key] } else { "?$key" }
        $ann += ('{0}=0x{1:X8} {2}+0x{3:X}' -f $m.Value, $v, $n, ($v - $base))
    }
    if ($ann.Count) { [void]$sb.AppendLine($line + '   /* ' + (($ann | Select-Object -Unique) -join '; ') + ' */') }
    else { [void]$sb.AppendLine($line) }
}
[IO.File]::WriteAllText((Join-Path (Get-Location) $Out), $sb.ToString())
"Wrote $Out"
