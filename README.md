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
Address: 137.0.41.190\30

Interface: Ethernet1  -> Address: 10.126.52.14\30
Interface: Ethernet2  -> Address: 137.3.20.249\30
Interface: Ethernet3  -> Address: 137.0.41.190\30
Interface: Ethernet4  -> Address: 137.0.41.1\30
Interface: Loopback0  -> Address: 137.2.20.1\32

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

## Router NTP config

```bash
localhost#sh run | i ntp
ntp server 100.64.1.1
ntp server 172.16.2.1
```
