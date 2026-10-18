param(
    [Parameter(Mandatory)] [string] $PdfPath,
    [Parameter(Mandatory)] [string] $OutDir,
    [double] $Scale = 3.0,
    [int] $TilesX = 3,
    [int] $TilesY = 3
)

Add-Type -AssemblyName System.Runtime.WindowsRuntime
Add-Type -AssemblyName System.Drawing
$null = [Windows.Storage.StorageFile, Windows.Storage, ContentType = WindowsRuntime]
$null = [Windows.Data.Pdf.PdfDocument, Windows.Data.Pdf, ContentType = WindowsRuntime]
$null = [Windows.Storage.Streams.InMemoryRandomAccessStream, Windows.Storage.Streams, ContentType = WindowsRuntime]

$asTaskGeneric = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object {
        $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncOperation`1'
    })[0]
$asTaskAction = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object {
        $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncAction'
    })[0]

function Await($op, [Type] $resultType) {
    $task = $asTaskGeneric.MakeGenericMethod($resultType).Invoke($null, @($op))
    $task.Wait() | Out-Null
    $task.Result
}
function AwaitAction($op) {
    $task = $asTaskAction.Invoke($null, @($op))
    $task.Wait() | Out-Null
}

New-Item -ItemType Directory -Force $OutDir | Out-Null
$file = Await ([Windows.Storage.StorageFile]::GetFileFromPathAsync((Resolve-Path $PdfPath).Path)) ([Windows.Storage.StorageFile])
$doc = Await ([Windows.Data.Pdf.PdfDocument]::LoadFromFileAsync($file)) ([Windows.Data.Pdf.PdfDocument])
$base = [IO.Path]::GetFileNameWithoutExtension($PdfPath) -replace '[^\w]', '_'

for ($i = 0; $i -lt $doc.PageCount; $i++) {
    $page = $doc.GetPage($i)
    $opts = New-Object Windows.Data.Pdf.PdfPageRenderOptions
    $opts.DestinationWidth = [uint32]($page.Size.Width * $Scale)
    $opts.DestinationHeight = [uint32]($page.Size.Height * $Scale)
    $opts.BackgroundColor = [Windows.UI.Color]::FromArgb(255, 255, 255, 255)

    $stream = New-Object Windows.Storage.Streams.InMemoryRandomAccessStream
    AwaitAction ($page.RenderToStreamAsync($stream, $opts))
    $netStream = [System.IO.WindowsRuntimeStreamExtensions]::AsStreamForRead($stream.GetInputStreamAt(0))
    $bmp = [System.Drawing.Bitmap]::FromStream($netStream)

    $full = Join-Path $OutDir ("{0}_p{1}_full.png" -f $base, ($i + 1))
    $bmp.Save($full, [System.Drawing.Imaging.ImageFormat]::Png)
    "page {0}: {1}x{2} -> {3}" -f ($i + 1), $bmp.Width, $bmp.Height, $full

    $tw = [int][Math]::Ceiling($bmp.Width / $TilesX)
    $th = [int][Math]::Ceiling($bmp.Height / $TilesY)
    for ($ty = 0; $ty -lt $TilesY; $ty++) {
        for ($tx = 0; $tx -lt $TilesX; $tx++) {
            $x = $tx * $tw; $y = $ty * $th
            $w = [Math]::Min($tw, $bmp.Width - $x); $h = [Math]::Min($th, $bmp.Height - $y)
            $rect = New-Object System.Drawing.Rectangle $x, $y, $w, $h
            $tile = $bmp.Clone($rect, $bmp.PixelFormat)
            $tile.Save((Join-Path $OutDir ("{0}_p{1}_r{2}c{3}.png" -f $base, ($i + 1), $ty, $tx)), [System.Drawing.Imaging.ImageFormat]::Png)
            $tile.Dispose()
        }
    }
    $bmp.Dispose(); $netStream.Dispose(); $stream.Dispose(); $page.Dispose()
}
