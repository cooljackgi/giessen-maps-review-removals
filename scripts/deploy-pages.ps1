param(
    [string]$CName = "",
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"

function Resolve-Tool {
    param(
        [string]$CommandName,
        [string[]]$FallbackPaths
    )

    $cmd = Get-Command $CommandName -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }

    foreach ($path in $FallbackPaths) {
        if (Test-Path $path) {
            return $path
        }
    }

    throw "Required tool not found: $CommandName"
}

function Invoke-Step {
    param(
        [string]$FilePath,
        [string[]]$Arguments,
        [string]$WorkingDirectory
    )

    Write-Host ">" $FilePath ($Arguments -join " ")
    & $FilePath @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code $LASTEXITCODE"
    }
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$gitExe = Resolve-Tool "git" @(
    "C:\Program Files\Git\cmd\git.exe",
    "C:\Program Files\Git\bin\git.exe"
)
$goExe = Resolve-Tool "go" @(
    "C:\Program Files\Go\bin\go.exe"
)

Push-Location $repoRoot
try {
    if (-not $SkipBuild) {
        Invoke-Step $goExe @("run", "./cmd/charts") $repoRoot
        Invoke-Step $goExe @("run", "./cmd/dashboard") $repoRoot
    }

    $publicDir = Join-Path $repoRoot "public"
    $chartsDir = Join-Path $repoRoot "output\charts"
    $metadataFile = Join-Path $repoRoot "output\metadata.json"
    $placesCsv = Join-Path $repoRoot "output\places.csv"
    $dashboardFile = Join-Path $chartsDir "giessen_dashboard.html"

    foreach ($requiredPath in @($chartsDir, $metadataFile, $placesCsv, $dashboardFile)) {
        if (-not (Test-Path $requiredPath)) {
            throw "Missing required build artifact: $requiredPath"
        }
    }

    if (Test-Path $publicDir) {
        Remove-Item -Recurse -Force $publicDir
    }
    New-Item -ItemType Directory $publicDir | Out-Null
    New-Item -ItemType Directory (Join-Path $publicDir "charts") | Out-Null
    New-Item -ItemType Directory (Join-Path $publicDir "data") | Out-Null
    New-Item -ItemType File (Join-Path $publicDir ".nojekyll") | Out-Null

    Copy-Item $dashboardFile (Join-Path $publicDir "index.html")
    Copy-Item (Join-Path $chartsDir "*") (Join-Path $publicDir "charts") -Recurse -Force
    Copy-Item $metadataFile (Join-Path $publicDir "data\metadata.json")
    Copy-Item $placesCsv (Join-Path $publicDir "data\places.csv")

    if ($CName.Trim() -ne "") {
        Set-Content -Path (Join-Path $publicDir "CNAME") -Value $CName.Trim() -NoNewline
    }

    $remoteUrl = (& $gitExe remote get-url origin).Trim()
    if (-not $remoteUrl) {
        throw "Could not resolve git remote 'origin'"
    }

    $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("giessen-pages-" + [System.Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory $tempDir | Out-Null

    try {
        & $gitExe ls-remote --exit-code --heads $remoteUrl gh-pages *> $null
        $hasGhPages = ($LASTEXITCODE -eq 0)

        if ($hasGhPages) {
            Invoke-Step $gitExe @("clone", "--quiet", "--branch", "gh-pages", "--single-branch", $remoteUrl, $tempDir) $repoRoot
        } else {
            Invoke-Step $gitExe @("clone", "--quiet", $remoteUrl, $tempDir) $repoRoot
            Push-Location $tempDir
            try {
                Invoke-Step $gitExe @("checkout", "--orphan", "gh-pages") $tempDir
            } finally {
                Pop-Location
            }
        }

        Get-ChildItem -Force $tempDir | Where-Object { $_.Name -ne ".git" } | Remove-Item -Recurse -Force
        Copy-Item (Join-Path $publicDir "*") $tempDir -Recurse -Force

        Push-Location $tempDir
        try {
            Invoke-Step $gitExe @("add", "-A") $tempDir
            & $gitExe diff --cached --quiet
            if ($LASTEXITCODE -eq 0) {
                Write-Host "gh-pages is already up to date."
            } else {
                Invoke-Step $gitExe @("commit", "-m", "Deploy GitHub Pages site") $tempDir
                Invoke-Step $gitExe @("push", "-u", "origin", "gh-pages") $tempDir
            }
        } finally {
            Pop-Location
        }
    } finally {
        if (Test-Path $tempDir) {
            Remove-Item -Recurse -Force $tempDir
        }
    }

    Write-Host "GitHub Pages deployment complete."
    Write-Host "Now set GitHub Pages source to branch 'gh-pages' and folder '/(root)'."
} finally {
    Pop-Location
}
