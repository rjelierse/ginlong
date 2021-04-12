# Protocol

The protocol is a binary encoding. Any number data is using little endian data coding.

## Message structure
```
00000000  a5 01 00 10 47 00 0c c8  d2 dd 2a 00 05 15        |....G.....*...|
```

### Fields
* `0x00-0x00` - `binary`: Message start, always `0xA5`.
* `0x01-0x01` - `uint8`: Payload length.
* `0x02-0x03` - `uint16`: Unknown, always `0x0010` (4096)
* `0x04-0x04` - `uint8`: Message type (?)
* `0x05-0x05` - `uint8`: Protocol version (?)
* `0x06-0x06` - `uint8`: Message counter. This field increments after each acknowledged message. If the response does not ack the message number, it will resend the message.
* `0x07-0x0a` - `uint32`: The data logger serial number.
* `0x0b-0x0b` - `binary`: Message payload (**note** the actual end address changes based on payload length)
* `0x0c-0x0c` - `uint8`: Message checksum. The checksum is calculated by adding the value of each byte from the payload length to the end of the payload.
* `0x0d-0x0d` - `binary`: Message terminator, always `0x15`.

## Message types
So far, I've observed the following message types.

### `0x41` - Data logger handshake
The handshake message does not have a protocol version.
```
00000000  02 7d e3 74 00 59 01 00  00 00 00 00 00 05 3c 78  |.}.t.Y........<x|
00000010  01 64 03 4d 57 5f 30 38  5f 30 35 30 31 5f 31 2e  |.d.MW_08_0501_1.|
00000020  35 38 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |58..............|
00000030  00 00 00 00 00 00 00 00  00 00 00 98 d8 63 e4 73  |.............c.s|
00000040  b2 31 39 32 2e 31 36 38  2e 31 2e 32 35 32 00 00  |.192.168.1.252..|
00000050  00 f8 4f 01 01 05                                 |..O...|
```

#### Fields
* `0x13-0x21` - `string`: Data logger firmware version.
* `0x3b-0x40` - `[]byte`: Hardware address.
* `0x41-0x4f` - `string`: Network address.

### `0x42` - Inverter data
The data format differs based on the protocol version.

For `0x02`, this is a sample payload:
```
// 00000000  01 07 05 fa e1 74 00 5c  01 00 00 e5 bb 00 60 01  |.....t.\......`.|
// 00000010  00 83 0f 00 00 31 36 30  45 33 31 32 30 34 31 35  |.....160E3120415|
// 00000020  30 30 36 38 20 fd 01 22  09 0d 00 62 00 01 00 00  |0068 .."...b....|
// 00000030  00 00 00 56 00 00 00 00  00 ab 09 88 13 5c 08 00  |...V.........\..|
// 00000040  00 80 02 00 00 f2 17 00  00 00 00 00 00 01 00 00  |................|
// 00000050  00 00 00 00 00 00 00 00  00 00 00 01 00 00 00 01  |................|
// 00000060  00 36 0f 00 00 f8 2a e8  03 f3 08 00 00 4d 00 00  |.6....*......M..|
// 00000070  00 90 00 00 00 60 00 53  01 00 00 12 01 00 00 34  |.....`.S.......4|
// 00000080  00 04 00 5c 08 00 00 15  00 04 00 0d 00 0e 00 22  |...\..........."|
// 00000090  00 0e 00 e3 00 26 00 40  00 00 00 01 00 00 00 00  |.....&.@........|
// 000000a0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000b0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000c0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000e0  00 00 00 00 00 00 00 03  00                       |.........|
```

For `0x03`, this is a sample payload:
```
// 00000000  81 07 05 cb e0 74 00 28  00 00 00 00 00 00 00 01  |.....t.(........|
// 00000010  00 82 0f 00 00 31 36 30  45 33 31 32 30 34 31 35  |.....160E3120415|
// 00000020  30 30 36 38 20 fa 01 fa  08 0c 00 5f 00 01 00 00  |0068 ......_....|
// 00000030  00 00 00 53 00 00 00 00  00 a5 09 87 13 02 08 00  |...S............|
// 00000040  00 6c 02 00 00 f2 17 00  00 00 00 00 00 01 00 00  |.l..............|
// 00000050  00 00 00 00 00 00 00 00  00 00 00 01 00 00 00 01  |................|
// 00000060  00 fe 0e 00 00 f8 2a e8  03 87 08 00 00 4d 00 00  |......*......M..|
// 00000070  00 90 00 00 00 60 00 53  01 00 00 12 01 00 00 34  |.....`.S.......4|
// 00000080  00 04 00 f8 07 00 00 15  00 04 00 0d 00 0e 00 1d  |................|
// 00000090  00 04 00 e3 00 26 00 40  00 00 00 01 00 00 00 00  |.....&.@........|
// 000000a0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000b0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000c0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
// 000000e0  00 00 00 00 00 00 00 03  00                       |.........|
```

#### Fields
* `0x15-0x23` - `string`: Inverter serial number.
* ...
* `0x93-0x94` - `uint16`: Inverter model.
* `0x95-0x96` - `uint16`: Inverter firmware version (slave).
* `0x97-0x98` - `uint16`: Inverter firmware version (main).

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

### `0x48` - Unknown
```
00000000  01 ee e1 75 00 f4 00 00  00 a9 49 01 60 01 05 2c  |...u......I.`..,|
00000010  84 e5 71 5e 4d 57 5f 30  38 5f 30 35 30 31 5f 31  |..q^MW_08_0501_1|
00000020  2e 35 38 00 00 00 00 00  00 00 00 00 00 00 00 00  |.58.............|
00000030  00 00 00 00 00 00 00 00  00 00 00 00              |............|
```

#### Fields
* `0x15-0x22` - `string`: Data logger firmware version.
