!define APPNAME "SV Printer"
!define COMPANYNAME "SV Tech"
!define DESCRIPTION "Cross-platform local agent for thermal printing"
!define VERSIONMAJOR 0
!define VERSIONMINOR 1
!define VERSIONBUILD 0

Name "${APPNAME}"
OutFile "SV_Printer_Setup.exe"
InstallDir "$PROGRAMFILES64\${COMPANYNAME}\${APPNAME}"

RequestExecutionLevel admin

Page directory
Page instfiles

Section "Install"
    SetOutPath "$INSTDIR"
    File "..\sv-printer.exe"
    File "..\assets\logo.ico"

    WriteUninstaller "$INSTDIR\uninstall.exe"

    CreateDirectory "$SMPROGRAMS\${COMPANYNAME}"
    CreateShortcut "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk" "$INSTDIR\sv-printer.exe" "" "$INSTDIR\logo.ico"
    
    ; Run the agent after installation completes
    Exec "$INSTDIR\sv-printer.exe"
SectionEnd

Section "Uninstall"
    ; Stop the process if running
    ExecWait "taskkill /F /IM sv-printer.exe"
    
    Delete "$INSTDIR\sv-printer.exe"
    Delete "$INSTDIR\logo.ico"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"
    
    Delete "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk"
    RMDir "$SMPROGRAMS\${COMPANYNAME}"

    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "sv-printer"
SectionEnd
