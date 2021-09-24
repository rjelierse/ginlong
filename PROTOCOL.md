# Protocol

**Firmware Version**: MW_08_0501_1.58

The protocol is a binary encoding. Any number data is using little endian data coding.

## Message structure
```
00000000  a5 01 00 10 47 00 0c c8  d2 dd 2a 00 05 15        |....G.....*...|
```

### Fields
* `0x00` - `binary`: Message start, always `0xA5`.
* `0x01` - `uint8`: Payload length.
* `0x02` - `uint16`: 4096
* `0x04` - `uint8`: Message type
* `0x05` - `uint8`: Counter 1 (increments on message from monitoring system) - value is ACK'ed on next message from logging stick.
* `0x06` - `uint8`: Counter 2 (increments on message from logging stick) - value is ACK'ed on next message from monitoring system.
* `0x07` - `uint32`: The data logger serial number.
* `0x0b` - `binary`: Message payload (**note** the actual end address changes based on payload length)
* `0x0c` - `uint8`: Message checksum. The checksum is calculated by adding the value of each byte from the payload length to the end of the payload.
* `0x0d` - `binary`: Message terminator, always `0x15`.

## Response structure
```
00000000  a5 0a 00 10 12 5b 2x c8  d2 dd 2a 00 01 d6 9e 48
00000010  61 78 00 00 00 06 15
```

### Fields
* `0x00` - `binary`: Message start, always `0xA5`.
* `0x01` - `uint8`: Payload length (always `0x0A`, response has fixed payload size it seems).
* `0x02` - `uint16`: 4096
* `0x04` - `uint8`: Message type (usually: `data logger message type - 0x30`).
* `0x05` - `uint8`: Counter 1 (increments on message from monitoring system) - value is ACK'ed on next message from logging stick.
* `0x06` - `uint8`: Counter 2 (increments on message from logging stick) - value is ACK'ed on next message from monitoring system.
* `0x07` - `uint32`: The data logger serial number.
* `0x0b` - `binary`: USR1
* `0x0c` - `binary`: USR2
* `0x0d` - `uint32`: Unix timestamp (time since epoch 01-01-1970T00:00:00)
* `0x0e` - `uint32`: 120 (?)
* `0x15` - `uint8`: Message checksum. The checksum is calculated by adding the value of each byte from the payload length to the end of the payload.
* `0x16` - `binary`: Message terminator, always `0x15`.

NOTE 1: Bytes `0x0b` - `0x14` are the message payload.
NOTE 2: The `USR1` and `USR2` fields seem to always have a fixed value based on the message type.

## Message types
So far, I've observed the following message types.

### `0x41` - Data logger handshake
```
00000000  02 7d e3 74 00 59 01 00  00 00 00 00 00 05 3c 78  |.}.t.Y........<x|
00000010  01 64 03 4d 57 5f 30 38  5f 30 35 30 31 5f 31 2e  |.d.MW_08_0501_1.|
00000020  35 38 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |58..............|
00000030  00 00 00 00 00 00 00 00  00 00 00 98 d8 63 e4 73  |.............c.s|
00000040  b2 31 39 32 2e 31 36 38  2e 31 2e 32 35 32 00 00  |.192.168.1.252..|
00000050  00 f8 4f 01 01 05                                 |..O...|
```

#### Fields
* `0x00` - `uint8`: Command type
* `0x13` - `string(40)`: Data logger firmware version (NUL padded).
* `0x3b` - `byte(6)`: Hardware address (each byte represents a part of the HW address as-is)
* `0x41` - `string(15)`: Network address (NUL padded).

#### Response fields
* `USR1`: `0x02`
* `USR2`: `0x01`

### `0x42` - Inverter data
```
00000000  81 07 05 17 f0 e4 00 39  06 00 00 00 00 00 00 01  |.......9........|
00000010  00 e9 07 00 00 31 36 30  45 33 31 32 30 34 31 35  |.....160E3120415|
00000020  30 30 36 38 20 51 01 be  08 0b 00 07 00 01 00 00  |0068 Q..........|
00000030  00 00 00 06 00 00 00 00  00 5f 09 88 13 64 00 00  |........._...d..|
00000040  00 50 00 00 00 8c 41 00  00 00 00 00 00 01 00 00  |.P....A.........|
00000050  00 00 00 00 00 00 00 00  00 00 00 01 00 00 00 01  |................|
00000060  00 5a 0e 00 00 f8 2a e8  03 9c 00 00 00 80 00 00  |.Z....*.........|
00000070  00 fc 00 00 00 24 00 7b  05 00 00 12 01 00 00 34  |.....$.{.......4|
00000080  00 04 00 64 00 00 00 15  00 08 00 13 00 11 00 04  |...d............|
00000090  00 0d 00 e3 00 26 00 40  00 00 00 01 00 00 00 00  |.....&.@........|
000000a0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
000000b0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
000000c0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
000000e0  00 00 00 00 00 00 00 03  00                       |.........|
```

#### Fields
* `0x01` - `uint16`: Sensor.
* `0x15` - `string(16)`: Inverter serial number.
* `0x25` - `uint16`: Inverter temperature [10^-1 C].
* `0x27` - `uint16`: PV array 1 voltage [10^-1 V].
* `0x29` - `uint16`: PV array 2 voltage [10^-1 V].
* `0x2b` - `uint16`: PV array 1 current [10^-1 A].
* `0x2d` - `uint16`: PV array 2 current [10^-1 A].
* `0x2f` - `uint16`: Unknown.
* `0x31` - `uint16`: Unknown.
* `0x33` - `uint16`: Grid output current, phase 1 [10^-1 A].
* `0x35` - `uint16`: Grid output current, phase 2 [10^-1 A].
* `0x37` - `uint16`: Grid output current, phase 3 [10^-1 A].
* `0x39` - `uint16`: Grid voltage [10^-1 V].
* `0x3b` - `uint16`: Grid frequency [10^-2 Hz].
* `0x3d` - `uint32`: Unknown.
* `0x41` - `uint32`: Daily generation [10^-2 kWh].
* `0x45` - `uint32`: Total generation [10^-1 kWh].
* `0x61` - `uint16`: V_Bus [10^-1 V].
* `0x63` - `uint16`: H_V_Bus [10^-1 V].
* `0x65` - `uint16`: Power limit [10^-2 %].
* `0x67` - `uint16`: Grid power factor [10^-3].
* `0x69` - `uint16`: Total PV input [W].
* `0x6d` - `uint32`: Generated this month [kWh].
* `0x71` - `uint32`: Generated last month [kWh].
* `0x75` - `uint16`: Generated yesterday [10^-1 kWh].
* `0x77` - `uint32`: Generated this year [kWh].
* `0x83` - `uint16`: Power grid: total apparent power [VA].
* `0x87` - `uint16`: Year of measurement result.
* `0x89` - `uint16`: Month of measurement result.
* `0x8b` - `uint16`: Day of measurement result.
* `0x8d` - `uint16`: Hour of measurement result.
* `0x8f` - `uint16`: Minute of measurement result.
* `0x91` - `uint16`: Second of measurement result.
* `0x93` - `uint16`: Inverter model.
* `0x95` - `uint16`: DSP version.
* `0x97` - `uint16`: HCI version.

#### Response fields
* `USR1`: `0x81`
* `USR2`: `0x01`

### `0x43` - Data logger wifi settings
```
00000000  81 ff e0 75 00 05 00 00  00 00 00 00 00 00 00 48  |...u...........H|
00000010  61 7a 65 73 00 00 00 00  00 00 00 00 00 00 00 00  |azes............|
00000020  00 00 00 00 00 00 00 00  00 00 00 00 00 64 01     |.............d.|
```

#### Fields
* `0x0f-0x2c` - `string`: Wifi SSID

### `0x47` - Ping
```
00000000  00                                                |.|
```

#### Fields
None

#### Response fields
* `USR1`: `0x00`
* `USR2`: `0x01`

### `0x48` - Unknown
```
00000000  01 ee e1 75 00 f4 00 00  00 a9 49 01 60 01 05 2c  |...u......I.`..,|
00000010  84 e5 71 5e 4d 57 5f 30  38 5f 30 35 30 31 5f 31  |..q^MW_08_0501_1|
00000020  2e 35 38 00 00 00 00 00  00 00 00 00 00 00 00 00  |.58.............|
00000030  00 00 00 00 00 00 00 00  00 00 00 00              |............|
```

#### Fields
* `0x15-0x22` - `string`: Data logger firmware version.

#### Response fields
* `USR1`: `0x01`
* `USR2`: `0x01`
