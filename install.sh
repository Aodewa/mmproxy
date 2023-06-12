#!/bin/bash

# download the mmproxy
URL="https://github.com/Aodewa/mmproxy/releases/download/v1.0.0/mmproxy_1.0.0.tar.gz"

echo "Download mmproxy from ${URL}"

cd /tmp
rm -rf mmproxy.tar.gz
wget ${URL} -O mmproxy.tar.gz
if [ $? -ne 0 ]; then
    echo 'download failed'
    exit 1
fi

tar -xzvf mmproxy.tar.gz
mv ./mmproxy /usr/local/bin/mmproxy

chmod +x /usr/local/bin/mmproxy

echo '{
    "router": [
        {
            "from":"0.0.0.0:1111",
            "to":"127.0.0.1:2222"
        }
    ],
    "p":"tcp",
    "mark":0,
    "v":0,
    "listeners":1
}' > /etc/mmproxy.json

# add service
echo '[Unit]
Description=mmproxy
After=network.target

[Service]
Type=simple
LimitNOFILE=65535
ExecStartPost=/sbin/ip rule add from 127.0.0.1/8 iif lo table 123
ExecStartPost=/sbin/ip route add local 0.0.0.0/0 dev lo table 123
ExecStartPost=/sbin/ip -6 rule add from ::1/128 iif lo table 123
ExecStartPost=/sbin/ip -6 route add local ::/0 dev lo table 123
ExecStart=/usr/local/bin/mmproxy -c /etc/mmproxy.json
ExecStopPost=/sbin/ip rule del from 127.0.0.1/8 iif lo table 123
ExecStopPost=/sbin/ip route del local 0.0.0.0/0 dev lo table 123
ExecStopPost=/sbin/ip -6 rule del from ::1/128 iif lo table 123
ExecStopPost=/sbin/ip -6 route del local ::/0 dev lo table 123
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target

' > /etc/systemd/system/mmproxy.service


echo 1 > /proc/sys/net/ipv4/conf/eth0/route_localnet

systemctl daemon-reload
