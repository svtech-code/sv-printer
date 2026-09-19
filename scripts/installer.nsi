!define APPNAME "SV Printer"
!define COMPANYNAME "SV Tech"
!define DESCRIPTION "Cross-platform local agent for thermal printing"

!ifndef VERSIONMAJOR
!define VERSIONMAJOR 0
!endif
!ifndef VERSIONMINOR
!define VERSIONMINOR 1
!endif
!ifndef VERSIONBUILD
!define VERSIONBUILD 0
!endif

!define VERSION "${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}"

Name "${APPNAME}"
OutFile "SV_Printer_Setup.exe"
InstallDir "$PROGRAMFILES64\${COMPANYNAME}\${APPNAME}"

VIProductVersion "${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}.0"
VIFileVersion "${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}.0"
VIAddVersionKey "ProductName" "${APPNAME}"
VIAddVersionKey "CompanyName" "${COMPANYNAME}"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"

RequestExecutionLevel admin

Page directory
Page instfiles

Section "Install"
    ; Ensure any running instance is closed before overwriting
    ExecWait "taskkill /F /IM sv-printer.exe"
    Sleep 1000

    SetOutPath "$INSTDIR"
    File "..\sv-printer.exe"
    File "..\assets\logo.ico"

    WriteUninstaller "$INSTDIR\uninstall.exe"

    CreateDirectory "$SMPROGRAMS\${COMPANYNAME}"
    CreateShortcut "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk" "$INSTDIR\sv-printer.exe" "" "$INSTDIR\logo.ico"
    CreateShortcut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\sv-printer.exe" "" "$INSTDIR\logo.ico"
SectionEnd

Section "Uninstall"
    ; Stop the process if running
    ExecWait "taskkill /F /IM sv-printer.exe"
    ; Wait for Windows to release file handles
    Sleep 1000
    
    Delete "$INSTDIR\sv-printer.exe"
    Delete "$INSTDIR\logo.ico"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"
    
    Delete "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk"
    Delete "$DESKTOP\${APPNAME}.lnk"
    RMDir "$SMPROGRAMS\${COMPANYNAME}"

    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "sv-printer"
SectionEnd
