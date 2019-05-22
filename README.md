# Overview
go3status is a simple and lightweigt replacement for i3status utility with programmatic click events

## Install
### Binary instyallation (preffered)
Download latest release on [project release page](https://github.com/pymhd/go3status/releases)
Unpack go3status binary, and base config file

#### From source
```sh
git clone https://github.com/pymhd/go3status.git
cd go3status
make
make install clean
```
go3status binary file will be placed in /usr/local/bin dir 

## Usage
Replace status_command option in i3 config file with:
```sh
status_command /path/to/binary/go3status -config /etc/config.yaml
```

## Click Events
There is 1 common predefined click event for all modules (middle mouse button)
It is used to switch between full and short module output. For every module
difference between full and short form is custom. Try it.

## Modules
### cpu
CPU usage percentage
```sh
- cpu:
    interval: 2s500ms
    prefix: "\uf2db "
    postfix: 
    colors:
      good: "#66b266"
      warn: "#e2c96e"
      crit: "#7f0909"
    levels:
      good: 0-50
      warn: 51-75
      crit: 75-100
    clickEvents:
      left: shell cmd
      right: shell cmd
      wheelUp: shell cmd
      wheelDown: shell cmd
    extra:
```
### Memory
```sh
- memory:
    interval: 5s
    prefix: " "
    postfix:
    colors:
      good: "#66b266"
      warn: "#e2c96e"
      crit: "#7f0909"
    levels:
      good: 0-50
      warn: 51-75
      crit: 75-100
    clickEvents:
      left: shell cmd
      right: shell cmd
      wheelUp: shell cmd
      wheelDown: shell cmd
    extra:
      format: '{{printf "%.1f" $used}}/{{printf "%.1f" $total}} ({{printf "%.0f" $percentage}}%)'
```
Memory module accepts extra arg - format. Go templating is supported, available vars:
  - $used - used memory
  - $total - total memory available
  - $percentage - percentage of used memory

### Disk
```sh
- hdd:
    interval: 10s
    prefix: "\uf120 " 
    postfix: 
    colors:
      good: "#66b266"
      warn: "#e2c96e"
      crit: "#7f0909"
    levels:
      good: 0-60
      warn: 61-80
      crit: 81-100
    clickEvents:
      left: thunar
      right:
      wheelUp:
      wheelDown:
    extra:
      mountPoint: /
      # $avail, $path, $percentage  also available as float64
      #format: '{{printf "%.1f" $used}}/{{printf "%.1f" $total}} ({{printf "%.0f" $percentage}}%)'
      format: '{{printf "%.1f" $avail}}'
```
HDD module acceps required extra arg - mountPoint and optional arg - format
Go templating is supported for format arg, available vars:
  - $avail - Gb available on mountpoint
  - $path - mountPoint
  - $percentage - percent of used space
  - $total - total Gb available
  - $used - Gb used on mountpoint
 
### Time
```sh
- time:
    interval: 1s
    prefix: "\uf274 "
    extra:
      format: 02/01 15:04
    clickEvents:
      left: 'xdg-open "https://calendar.google.com"'
      right: 'xdg-open "https://calendar.google.com"'
```
Memory module accepts extra arg - format. Specify format using "time" Go package layout

### Docker
```sh
- docker:
    interval: 5s
    prefix: " "
    postfix:
    extra:
      # docker version |grep API|tail -n1 |awk '{print $3}'
      clientAPIVersion: "1.39"
      color: "#1bbbd2"
```
Docker module acceps required extra arg - clientAPIVersion and optional - color.
To get API version value run:
```sh
docker version |grep API|tail -n1 |awk '{print $3}'
```
### Weather 
```sh
- weather:
    interval: 15m
    extra:
      location: auto
```
Weather module accepts optional arg - location. If it is missing, "auto" value will be used. Your location will be determined via ipinfo.io service based on your ip address. Or you can specify your city explicit (Moscow, London...)
### Title
This module will get focused window name via i3ipc socket
```sh
- title:
    interval: 1s
    prefix: "(  "
    postfix: "  )"
    extra:
      maxChars: 64
```
Title module accepts optional arg - maxChars. With it only this number of first chars of window name will be showed

### Exec
With exec module you can exec shell command/script and expose it's ouptut. Bash interpreter is used
```
- exec:
    interval: 30m
    prefix: "\uf153 " 
    postfix:
    extra:
      cmd: "cat /proc/loadavg | awk '{print $1}'"
      color: "#bf3e3e"
      timeout: 2s
      cache: true
      update: 1h
```
Exec module accepts several extra args:
  - cmd - command to exec (required)
  - color - color of text output
  - timeout - command timeout
  - cache - true/false - if enabled - command will be executed in background periodically and its value will be stored in cache.
  - updade - interval how often to update cache

### Batt
This module show battery information
```
- batt:
    interval: 1s
    prefix: #"\uf5df "
    postfix:
    short: true
    colors:
      good: "#66b266"
      warn: "#e2c96e"
      crit: "#cc2222"
    levels:
      good: 50-100
      warn: 20-50
      crit: 0-20
``` 
### Network
This module finds out an active network interface, determinate his type and show information about it. Colors and level used only for wireless link quality level information.
```
- network:
    interval: 1s
    short: true
    colors:
      good: "#66b266"
      warn: "#e2c96e"
      crit: "#cc2222"
    levels:
      good: 61-100
      warn: 21-60
      crit: 0-20
```
### VPN
This module collects information about all tun and vpn NICs 
colors and levels: 
good = as minimum one VPN NIC is present and active connected
crit = no VPNs were connected
```
- vpn:
    interval: 1s
    short: true
    colors:
      good: "#66b266"
      crit: "#cc2222"
    levels:
      good: 51-100
      crit: 0-51
```
