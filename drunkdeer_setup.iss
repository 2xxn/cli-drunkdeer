[Setup]
AppName=A CLI Tool for DrunkDeer Keyboards
AppVersion=1.0
DefaultDirName={commonpf}\CLIDrunkDeer
DefaultGroupName=CLIDrunkDeer
OutputDir=.
OutputBaseFilename=CLIDrunkDeerInstaller
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin

[Files]
Source: "drunkdeer.exe"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
; Add installation path to system PATH variable
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Flags: preservestringtype

[Run]
; Display a success message in the terminal (hidden cmd window)
Filename: "{cmd}"; Parameters: "/C echo Installation complete! You can now use 'drunkdeer' from any terminal."; Flags: runhidden
