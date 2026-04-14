$ErrorActionPreference = "Stop"

$repo = if ($env:D2A_REPO) { $env:D2A_REPO } else { "fanfeilong/dot_20_arch_draft" }
$version = if ($env:D2A_VERSION) { $env:D2A_VERSION } else { "latest" }
$installDir = if ($env:D2A_INSTALL_DIR) { $env:D2A_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "d2a\bin" }
$baseUrlOverride = $env:D2A_BASE_URL

function Get-ReleaseBase {
    if ($baseUrlOverride) {
        return $baseUrlOverride
    }

    if ($version -eq "latest") {
        return "https://github.com/$repo/releases/latest/download"
    }

    return "https://github.com/$repo/releases/download/$version"
}

function Get-Arch {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture

    switch ($arch) {
        "X64" { return "amd64" }
        "Arm64" { return "arm64" }
        default { throw "unsupported architecture: $arch" }
    }
}

$arch = Get-Arch
$asset = "d2a_windows_${arch}.zip"
$baseUrl = Get-ReleaseBase
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("d2a-install-" + [System.Guid]::NewGuid().ToString("N"))
$archivePath = Join-Path $tmpDir $asset

New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
    Write-Host "Downloading $asset from $repo ($version)..."
    Invoke-WebRequest -Uri "$baseUrl/$asset" -OutFile $archivePath

    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force
    Copy-Item -Path (Join-Path $tmpDir "d2a.exe") -Destination (Join-Path $installDir "d2a.exe") -Force

    Write-Host "Installed d2a to $(Join-Path $installDir 'd2a.exe')"
    Write-Host "Run: d2a help"
} finally {
    if (Test-Path $tmpDir) {
        Remove-Item -Path $tmpDir -Recurse -Force
    }
}
