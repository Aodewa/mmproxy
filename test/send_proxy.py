import socket 

target_host = "127.0.0.1"
# target_port = 2222
target_port = 8443


client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

client.connect((target_host,target_port))

client.send(b"PROXY TCP4 0.0.0.0 0.0.0.0 8443 443\r\nGET / HTTP/1.1\r\nHost: baidu.com\r\n\r\n")
response = client.recv(4096)
print(response)
