# gNMI Test 

Powered by `ygnmi`.

```bash
$ go run *.go
Supported Encodings:
  - JSON
  - JSON_IETF
  - ASCII
gNMI Version: 0.7.0

Path: elem:{name:"interfaces"} elem:{name:"interface" key:{key:"name" value:"Ethernet3"}} elem:{name:"subinterfaces"} elem:{name:"subinterface" key:{key:"index" value:"0"}} elem:{name:"ipv4"}
Address: 203.0.113.190\30

Path: elem:{name:"interfaces"} elem:{name:"interface" key:{key:"name" value:"*"}} elem:{name:"subinterfaces"} elem:{name:"subinterface" key:{key:"index" value:"0"}} elem:{name:"ipv4"} elem:{name:"addresses"}
Value 0: Address: 172.16.52.14\30
Value 1: Address: 100.64.20.249\30
Value 2: Address: 203.0.113.190\30
Value 3: Address: 198.51.100.1\30
Value 4: Address: 192.0.2.1\32
```