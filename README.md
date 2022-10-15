# virtual-switch
&emsp;&emsp;一款虚拟交换机软件。可以把异地主机接入一个虚拟局域网下，实现异地组网，配合NAT软件还可以实现内网穿透。

&emsp;&emsp;本软件通过虚拟网卡与系统通信，工作在[ OSI 参考模型](https://cloud.tencent.com/developer/article/1870180)的第 2 层，可以转发所有上层流量。下面有两种典型应用场景。  
- 异地组网  
&emsp;&emsp;在一台有公网 IP 的服务器上部署服务端，在内网主机上安装客户端，主机通过虚拟局域网访问内网（服务器不接入虚拟局域网）。  
&emsp;&emsp;如果希望服务器也接入局域网，需要在服务器上也安装一个本软件客户端。  
&emsp;&emsp;本软件不提供自动配置 IP 地址功能，需要为每一台主机配置静态 IP 地址，或者在其中一台装有客户端的主机上配置 DHCP 功能。

- 内网穿透  
&emsp;&emsp;本软件可以为异地主机交换网络流量，自然也可以为内网主机和服务器交换网络流量。只需要在服务器上同时安装客户端，让服务器也接入网络。再配置流量转发功能，比如[ NAT (网络地址转换)](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E5%9C%B0%E5%9D%80%E8%BD%AC%E6%8D%A2)功能，或者像[ Nginx 软件](https://zh.wikipedia.org/wiki/Nginx)中的[反向代理](https://zh.wikipedia.org/wiki/%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86)功能。

```text
                  ╭─────────╮
                  │  Cloud  │
                  │ (Server)│
                  ╰───↑─↑───╯
                      │ │
             ┌────────┘ └────────┐
    ╭┈┈┈┈┈┈┈┈│┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈│┈┈┈┈┈┈┈┈╮
    ┊   ╭────↓────╮         ╭────↓────╮   ┊
    ┊   │  PC-1   │   ...   │  PC-N   │   ┊
    ┊   │ (Client)│         │ (Client)│   ┊
    ┊   ╰─────────╯         ╰─────────╯   ┊
    ┊             Virtual LAN             ┊
    ╰┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈╯
                A. 异地组网模式


            ╭───────╮     ╭───────╮
            │  User │ ... │  User │
            ╰───↑───╯     ╰───↑───╯
                └──────┬──────┘
                  ╭────↓────╮
                  │   NAT   │
                  │ (Server)│
                  ╰───↑─↑───╯
                      │ │
             ┌────────┘ └────────┐
             │                   │
        ╭────↓────╮         ╭────↓────╮
        │   PC-1  │   ...   │   PC-N  │
        │ (Client)│         │ (Client)│
        ╰─────────╯         ╰─────────╯
                 B. 内网穿透模式
```

## 支持环境
- Linux
- macOS
- Windows (需要安装[ OpenVPN TAP Windows 驱动](https://build.openvpn.net/downloads/releases/latest/tap-windows-latest-stable.exe))

## 目录说明
```text
vswitch
 ├─client            # 客户端代码
 ├─cmd
 │  ├─vswitchc       # 客户端入口
 │  └─vswitchs       # 服务端入口
 ├─out               # 编译输出
 ├─pkg
 │  ├─common
 │  ├─config         # 软件配置项
 │  ├─consts         # 常量
 │  ├─kcp            # kcp 协议栈(有修改)
 │  ├─mactable       # Mac 地址表模块
 │  ├─pkgbuf         # 本软件数据包缓存IO
 │  ├─porttable      # 交换机端口表模块
 │  ├─util
 │  │  ├─log         # Beego 日志模块
 │  │  └─util
 │  └─virtualswitch  # 流量交换模块 
 └─server            # 服务端代码
```

## 数据报文结构
```text
             本软件转发流量数据包
      0               16              32
      ┌───────────────┬───────────────┐
      │      Flag     │ Payload Length│
      ├───────────────┴───────────────┤
      │            Payload            │
      └───────────────────────────────┘
     ╱                                 ╲
    ╱                                   ╲
   ╱                                     ╲
    6       6     4       46~1500       4 
┌───────┬───────┬────┬───────────────┬─────┐
│  Des  │  Src  │Type│    Payload    │ FCS │ ← 以太网 MAC 帧
└───────┴───────┴────┴───────────────┴─────┘

Flag          : 帧头标记，固定为 0xBCBC
Payload Length: 有效载荷长度，字节为单位
Payload       : 有效载荷

Des           : 目的 MAC 地址
Src           : 源 MAC 地址
Type          : 上层协议的类型
Payload       : 有效载荷
FCS           : 帧检验序列，CRC校验
```
