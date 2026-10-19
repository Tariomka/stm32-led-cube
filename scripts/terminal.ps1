# Serial terminal for the cube console. Typed lines are sent as "#?<line>?!" (the vendor console
# framing), everything received is printed. With -Command it sends one command, prints the reply
# for -Seconds and exits. Ctrl+C quits.
param(
    [string] $Port = "COM4",
    [int] $Baud = 115200,
    [string] $Command = "",
    [double] $Seconds = 2,
    [switch] $Raw # send lines as typed, without "#?...?!"
)

function Format-Line([string] $line) {
    if ($Raw) { return $line + "`r`n" }
    return "#?" + $line + "?!"
}

$serial = New-Object System.IO.Ports.SerialPort $Port, $Baud, ([System.IO.Ports.Parity]::None), 8, ([System.IO.Ports.StopBits]::One)
$serial.Encoding = [System.Text.Encoding]::GetEncoding(28591) # bytes as-is
$serial.DtrEnable = $false
$serial.RtsEnable = $false
$serial.Open()
try {
    if ($Command) {
        $serial.Write((Format-Line $Command))
        $end = (Get-Date).AddSeconds($Seconds)
        while ((Get-Date) -lt $end) {
            if ($serial.BytesToRead -gt 0) { [Console]::Write($serial.ReadExisting()) }
            Start-Sleep -Milliseconds 20
        }
        return
    }

    Write-Host "Connected to $Port at $Baud. Type a command (e.g. Help) and press Enter. Ctrl+C quits." -ForegroundColor Cyan
    $line = ""
    while ($true) {
        if ($serial.BytesToRead -gt 0) { [Console]::Write($serial.ReadExisting()) }
        while ([Console]::KeyAvailable) {
            $key = [Console]::ReadKey($true)
            if ($key.Key -eq [ConsoleKey]::Enter) {
                [Console]::WriteLine()
                if ($line) { $serial.Write((Format-Line $line)) }
                $line = ""
            } elseif ($key.Key -eq [ConsoleKey]::Backspace) {
                if ($line.Length -gt 0) { $line = $line.Substring(0, $line.Length - 1); [Console]::Write("`b `b") }
            } elseif ($key.KeyChar -ne [char]0) {
                $line += $key.KeyChar
                [Console]::Write($key.KeyChar)
            }
        }
        Start-Sleep -Milliseconds 10
    }
} finally {
    $serial.Close()
}
