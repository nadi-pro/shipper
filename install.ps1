# Check if running with administrative privileges
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "Please run this script as an administrator."
    exit
}

$GH_REPO = "nadi-pro/shipper"
$TIMEOUT = 90

# Get the current logged-in user
$USERNAME = $env:USERNAME

$VERSION = (Invoke-WebRequest -Uri "https://api.github.com/repos/$GH_REPO/releases/latest" -UseBasicParsing).Content | ConvertFrom-Json | Select-Object -ExpandProperty tag_name
if (-not $VERSION) {
    Write-Host "`nThere was an error trying to check what is the latest version of shipper.`nPlease try again later.`n"
    exit 1
}

$OS_type = $env:PROCESSOR_ARCHITECTURE
switch ($OS_type) {
    "AMD64", "x86_64" {
        $OS_type = "amd64"
    }
    "x86", "i386" {
        $OS_type = "386"
    }
    "ARM64" {
        $OS_type = "arm64"
    }
    default {
        Write-Host "OS type not supported"
        exit 2
    }
}

$GH_REPO_BIN = "shipper-${VERSION}-windows-${OS_type}.tar.gz"

# Create tmp directory
$TMP_DIR = New-TemporaryFile -Directory | Select-Object -ExpandProperty FullName
Write-Host "Change to temporary directory $TMP_DIR"
Set-Location $TMP_DIR

Write-Host "Downloading shipper $VERSION"
$LINK = "https://github.com/$GH_REPO/releases/download/$VERSION/$GH_REPO_BIN"

Invoke-WebRequest -Uri $LINK -OutFile "$TMP_DIR\$GH_REPO_BIN"
if ($?) {
    Write-Host "Error downloading"
    exit 2
}

$BINARY_PATH = "C:\Program Files\Nadi-Pro\Shipper"
$null = New-Item -Path $BINARY_PATH -ItemType Directory -Force

Copy-Item -Path "$TMP_DIR\shipper.exe" -Destination $BINARY_PATH -Force
if ($?) {
    exit 2
}

$BINARY_DIRECTORY = "C:\ProgramData\Nadi-Pro\Shipper"
$null = New-Item -Path $BINARY_DIRECTORY -ItemType Directory -Force

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/nadi-pro/shipper/master/nadi.reference.yaml" -OutFile "$BINARY_DIRECTORY\nadi.yaml"
Write-Host "Downloaded nadi.reference.yaml and saved as $BINARY_DIRECTORY\nadi.yaml."

Remove-Item -Path $TMP_DIR -Recurse -Force
Write-Host "Installed successfully to $BINARY_PATH\shipper.exe"

# Create the service
$SERVICE_NAME = "Shipper"
$SERVICE_PATH = "$BINARY_PATH\shipper.exe"
$SERVICE_CONFIG_PATH = "$BINARY_DIRECTORY\nadi.yaml"

$serviceParams = @{
    Name             = $SERVICE_NAME
    BinaryPathName   = "$SERVICE_PATH --config=$SERVICE_CONFIG_PATH --record"
    DisplayName      = $SERVICE_NAME
    Description      = "Shipper Service"
    StartupType      = "Automatic"
    Credential       = "LocalSystem"
    DependsOn        = @("tcpip")
    ErrorControl     = "Normal"
    ServiceArguments = @()
}

$service = New-Service @serviceParams

if ($service) {
    Write-Host "Service created successfully."
} else {
    Write-Host "Failed to create the service."
}
