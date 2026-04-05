# dlms-cosem-go

**Complete DLMS/COSEM protocol stack for Go** — modular, sans-io implementation with HDLC framing, A-XDR codec, ASN.1 BER/DER, COSEM IC classes, security suites, server/client, and multiple transport layers.

[![Tests](https://img.shields.io/badge/tests-362%20passed-brightgreen)]()
[![Go 1.21+](https://img.shields.io/badge/go-1.21%2B-00ADD8.svg)]()
[![License: BSL 1.1](https://img.shields.io/badge/license-BSL%201.1-orange.svg)]()

## Features

- **ASN.1 BER/DER**: Complete tag-length-value encoding/decoding with property-based tests
- **A-XDR Codec**: DLMS-specific data encoding (octet string, bit string, visible string, integer, etc.)
- **HDLC Framing**: Full LLC/MAC layer with segmentation, CRC-16, and frame parsing
- **COSEM IC Classes**: Interface class implementations for smart metering
- **Security Suites**: Suite 0-5 with AES-GCM, SM4-GCM, SM4-GMAC authentication
- **Server/Client**: DLMS server with association handling and client API
- **Transport Layers**: TCP, serial (RS485), and mock transports

## Project Structure

```
├── asn1/       # ASN.1 BER/DER encoding
├── axdr/       # DLMS A-XDR codec
├── client/     # DLMS client implementation
├── core/       # Core types and constants
├── cosem/      # COSEM interface classes
├── hdlc/       # HDLC framing layer
├── security/   # Cryptographic operations (AES, SM4, GMAC)
├── server/     # DLMS server implementation
└── transport/  # Transport layer abstractions
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ViewWay/dlms-cosem-go/client"
)

func main() {
    // Create DLMS client
    c := client.New(client.Config{
        LogicalAddress: 1,
        PhysicalAddress: 16,
    })
    fmt.Println(c)
}
```

## Testing

```bash
go test ./... -v
```

## Multi-Language Family

| Language | Tests | Lines |
|----------|-------|-------|
| [Python](https://github.com/ViewWay/dlms-cosem) | 5,146 | 37K |
| [Rust](https://github.com/ViewWay/dlms-cosem-rust) | 739 | 21K |
| **Go** | **362** | **~8K** |
| [C++](https://github.com/ViewWay/dlms-cosem-cpp) | 280+ | 6.5K |
| [C](https://github.com/ViewWay/dlms-cosem-c) | 36 | 6.2K |

## License

BSL 1.1
