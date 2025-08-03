# Requires PowerShell 7 or later
# This script will download the AlmaLinux 9 rootfs, create a temporary WSL instance,
# attempt to install laravel, run tests, and then clean up the instance.

# Define variables
$InstanceName = "AlmaLinux-9-Laravel-Test"
$RootfsDir = "./$InstanceName-rootfs"
$RootfsUrl = "https://repo.almalinux.org/almalinux/9/BaseOS/x86_64/os/images/rootfs.tar.xz"
$RootfsFile = "almalinux9-rootfs.tar.xz"

# --- Main Script Execution ---
try {
    # 1. Download the AlmaLinux 9 rootfs tarball
    Write-Host "Downloading AlmaLinux 9 rootfs..." -ForegroundColor Green
    Invoke-WebRequest -Uri $RootfsUrl -OutFile $RootfsFile

    # 2. Build a new temporary WSL instance of AlmaLinux 9
    Write-Host "Building a new temporary WSL instance of AlmaLinux 9..." -ForegroundColor Green
    # The -ErrorAction Stop will stop execution if the command fails
    wsl --import $InstanceName $RootfsDir $RootfsFile -ErrorAction Stop
    Write-Host "Successfully imported instance '$InstanceName'." -ForegroundColor Green

    # Set AlmaLinux 9 as the default WSL instance for this script's context
    Write-Host "Setting '$InstanceName' as the temporary default instance..." -ForegroundColor Yellow
    wsl --set-default $InstanceName

    # 3. Attempt to install "laravel" using the built-in package manager (dnf)
    Write-Host "Attempting to install 'laravel' via dnf..." -ForegroundColor Green
    # Run the install command and check the exit code
    wsl --distribution $InstanceName dnf -y install laravel
    if ($LASTEXITCODE -ne 0) {
        throw "Error: Failed to install 'laravel'."
    }
    Write-Host "SUCCESS: 'laravel' package installed." -ForegroundColor Green

    # 4. Run a series of test commands
    Write-Host "Running test commands inside the '$InstanceName' instance..." -ForegroundColor Green

    # Test 1: Check if the 'laravel' command exists
    Write-Host "Test 1: Checking for 'laravel' command..." -ForegroundColor Yellow
    wsl --distribution $InstanceName command -v laravel
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: 'laravel' command found." -ForegroundColor Green
    } else {
        Write-Host "FAILURE: 'laravel' command not found." -ForegroundColor Red
    }

    # Test 2: Check the version of laravel
    Write-Host "Test 2: Getting laravel version..." -ForegroundColor Yellow
    wsl --distribution $InstanceName laravel --version
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Laravel version output successfully." -ForegroundColor Green
    } else {
        Write-Host "FAILURE: Could not get laravel version." -ForegroundColor Red
    }

    # Test 3: Attempt to create a new laravel project
    Write-Host "Test 3: Creating a new laravel project..." -ForegroundColor Yellow
    wsl --distribution $InstanceName "cd /tmp && laravel new my-test-app"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: New laravel project 'my-test-app' created." -ForegroundColor Green
    } else {
        Write-Host "FAILURE: Failed to create new laravel project." -ForegroundColor Red
    }

    Write-Host "All tests completed." -ForegroundColor Green
}
catch {
    Write-Host "An error occurred: $($_.Exception.Message)" -ForegroundColor Red
}
finally {
    # 5. Cleanup: Unregister and remove the temporary instance
    Write-Host "Starting cleanup process..." -ForegroundColor Yellow

    # Terminate the instance gracefully
    Write-Host "Terminating WSL instance '$InstanceName'..." -ForegroundColor Yellow
    wsl --terminate $InstanceName -ErrorAction SilentlyContinue

    # Unregister the instance
    Write-Host "Unregistering WSL instance '$InstanceName'..." -ForegroundColor Yellow
    wsl --unregister $InstanceName -ErrorAction SilentlyContinue

    # Remove the downloaded rootfs and the temporary directory
    if (Test-Path $RootfsFile) {
        Write-Host "Removing temporary file '$RootfsFile'..." -ForegroundColor Yellow
        Remove-Item -Path $RootfsFile -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path $RootfsDir) {
        Write-Host "Removing temporary directory '$RootfsDir'..." -ForegroundColor Yellow
        Remove-Item -Path $RootfsDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    Write-Host "Cleanup finished. Script completed." -ForegroundColor Green
}
