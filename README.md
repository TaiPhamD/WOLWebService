# WOLWebService
## summary
This is a webservice that runs in 2 modes:
  - Master: This mode is meant for an always on device like a raspberry pi or a unix based router like the UDM PRO/SE
      - It could relay Wake-On-Lan commands to your PC to wake it up.
      - Query slave system to get current OS information
      - Reboot slave system and change its UEFI BootNext variable with the specified BootID
  - Slave: This mode is meant for the client computer that would fulfill restart requests or OS query

A typical use case would be to install this server app on an always on device like a raspberry pi then set it up in Master mode. Then configure config.json to include information about all your LAN devices that you want to be able to send WoL packets, restart command, or query OS info.

You do not need to install this software on the WoL targeted system unless you want to support other APIs like restart to a certain OS for multi-boot system or to query OS info. 

## API
- POST /api/wol   - This API sends a Wake-On-Lan packet to your client PC defined by the alias param. Only available in Master mode.
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1"
}
```
- POST /api/os   - This API query OS type from your client PC defined by the alias param
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1"
}
```
- POST /api/restart   - This API restarts the client PC based on the OS parameter
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1",
        "os": "ubuntu"
}
```

## Config file explanation

```
{
    "port": "9991", <--- Listening port of web appp server
    "api_key": "my_secure_password", <--- Password to authorize api calls
    "fullchain": "certs/fullchain.pem", <--- optional TLS certs for HTTPS hosting
    "priv_key": "certs/privkey.pem",<--- optional TLS certs for HTTPS hosting
    "clients:": [ <--- Add all PCs on your LAN that you want WOL control here
        {
            "alias": "client1", <--- alias used to select the right PC . aka mapping to IP/MAC info.
            "ip": "192.168.2.23",
            "mac": "00:00:00:00:00:00"
        },
        {
            "alias": "client2",
            "ip": "192.168.0.27",
            "mac": "aa:aa:aa:aa:aa:aa" 
        }
    ]
}
```

# Build from source
## Linux
### Prerequisite
- Go lang compiler 
### Build step
- ./build_unix.sh
### Install step
- ./install_linux.sh (It will install a systemd service named wolservice and start it but it wont work yet until you setup a config.json)
- setup config.json based on the examples from [master config](https://github.com/TaiPhamD/WOLWebService/blob/master/server/config/config_master.json)
or [slave config](https://github.com/TaiPhamD/WOLWebService/blob/master/server/config/config_slave.json)
- copy config.json to /opt/wolservice/config/config.json
- sudo systemctl restart wolservice
## Windows
### Prerequisite
- MSYS2 for GCC tool chain (https://www.msys2.org/). Install the following packages:
    - pacman -S --needed base-devel mingw-w64-x86_64-toolchain
    - pacman -S cmake msys2-w32api-headers msys2-w32api-runtime
- Golang 1.16+
- CMake v3+
### Build step
- ./build_windows.bat
### Install step
- ./install_windows.bat (It will install a Windows Service called WOLServerService
- setup config.json based on the examples from [master config](https://github.com/TaiPhamD/WOLWebService/blob/master/server/config/config_master.json)
or [slave config](https://github.com/TaiPhamD/WOLWebService/blob/master/server/config/config_slave.json)
- copy config.json to C:\wolservice\config\config.json
- Restart the WOLServerService from services.msc
