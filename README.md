## 使用 golang 分配 ip

ipam.json 文件导入多段 IP
```yaml
{
    "ranges": [
      {
        "start": "10.172.16.2",
        "end": "10.172.16.3"
      },
      {
        "start": "10.172.17.2",
        "end": "10.172.17.3"
      }
    ]
}
```