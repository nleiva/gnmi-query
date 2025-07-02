# gNMI Test Application

A quick gNMI (gRPC Network Management Interface) test application that demonstrates different network device management operations using the [ygnmi](https://github.com/openconfig/ygnmi) Go package to interact with Arista devices.

## Features

This application showcases three key gNMI operations:

1. **Single Path Query** - Retrieve specific configuration values
2. **Wildcard Path Queries** - Retrieve multiple values using path wildcards  
3. **Configuration Reconciliation** - Continuously monitor and enforce desired state

## Prerequisites

- Go 1.24.4 or later
- Network connectivity to target Arista device
- Valid credentials for the target device

## Configuration

Set up proxy bypass if needed:
```bash
export no_proxy=$no_proxy,<taget-ip-address>
```

## Usage

Run the application:
```bash
cd arista/
go run *.go
```

### Sample Output

The application will connect to the Arista device and perform the following operations:

#### 1. Connection and Capabilities
```
Supported Encodings:
  - JSON
  - JSON_IETF
  - ASCII
gNMI Version: 0.7.0
```

#### 2. Single Path Query
Retrieves IPv4 configuration for a specific interface:
```
Path: /interfaces/interface[name=Ethernet3]/subinterfaces/subinterface[index=0]/ipv4
Address: 100.0.41.190/30
```

#### 3. Wildcard Path Queries
Retrieves IPv4 addresses for all interfaces:
```
Interface: Ethernet1  -> Address: 10.126.52.14/30
Interface: Ethernet2  -> Address: 100.3.20.249/30
Interface: Ethernet3  -> Address: 100.0.41.190/30
Interface: Ethernet4  -> Address: 100.0.41.1/30
Interface: Loopback0  -> Address: 100.2.20.1/32
```

#### 4. Configuration Reconciliation
Monitors and enforces NTP server configuration:

```
Path: elem:{name:"system"} elem:{name:"ntp"}
>>>>> unexpected cfg diff detected:
   &arista.System_Ntp{
        ... // 2 identical fields
        Enabled: nil,
        NtpKey:  nil,
        Server: map[string]*arista.System_Ntp_Server{
                "100.64.1.1": &{Address: &"100.64.1.1"},
-               "172.16.2.1": &{Address: &"172.16.2.1"},
        },
  }

config enforced at: 2025-06-25 17:14:07 for origin:"openconfig" elem:{name:"system"} elem:{name:"ntp"}
```

When a configuration drift is detected, the reconciler automatically corrects it by applying the desired state.

## Architecture

The application consists of three main components:

### `main.go`
- **Single Path Queries**: Uses `ygnmi.Get()` to retrieve specific configuration values
- **Wildcard Queries**: Uses `ygnmi.LookupAll()` to retrieve multiple values with path wildcards
- **Reconciliation**: Uses `ygnmi.NewReconciler()` to continuously monitor and enforce desired configuration state

### `client.go`
- Establishes secure gRPC connection to the Arista device
- Handles TLS configuration with certificate verification options
- Implements capabilities negotiation to determine supported features
- Creates ygnmi client instance for high-level operations

### `creds.go`
- Implements per-RPC credentials for authentication
- Provides username/password authentication mechanism
- Supports TLS requirement enforcement

## Key Dependencies

- **[ygnmi](https://github.com/openconfig/ygnmi)**: High-level gNMI client library
- **[ygot](https://github.com/openconfig/ygot)**: YANG-to-Go code generation and utilities
- **[gnmi](https://github.com/openconfig/gnmi)**: Core gNMI protocol definitions
- **[YANG data structures](https://github.com/nleiva/yang-data-structures)**: Generated Go structs for Arista device models

## Configuration Reconciliation Details

The reconciler demonstrates intent-based network management:

1. **Desired State Definition**: Defines the intended NTP server configuration
2. **Continuous Monitoring**: Watches for configuration changes in real-time
3. **Drift Detection**: Compares actual vs. desired state using deep comparison
4. **Automatic Remediation**: Applies corrections when drift is detected
5. **Audit Trail**: Logs all configuration changes with timestamps

This approach ensures network devices maintain their intended configuration even when manual changes occur.

Example router NTP configuration before reconciliation:

```bash
localhost#sh run | i ntp
ntp server 100.64.1.1
ntp server 172.16.2.1
```
