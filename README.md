## 说明

socks5 项目 fork 自 https://github.com/armon/go-socks5，原项目已经多年未更新了，加上想实现支持 server、支持 http proxy 的需求，于是 fork 一份来支持这些需求。

## 有什么功能？
- 支持 socks5 server
- 支持 http proxy

## 如何使用？

```bash
go install github.com/iamlongalong/proxyserver@latest

# 运行
proxyserver
```

参数控制：(环境变量)
- PROXY_USER: 代理用户名 (默认 空)
- PROXY_PASSWORD: 代理密码 (默认 空)
- SOCKS5_PROXY_PORT: socks5 监听端口 (默认 10801)
- HTTP_PROXY_PORT: http proxy 监听端口 (默认 10802)

例如，可以使用如下:
```bash
PROXY_USER=user PROXY_PASSWORD=pass proxyserver
```

## 客户端
你可以直接用:
```bash
# 使用 http proxy
export HTTP_PROXY=http://xx.xx.xx.xx:10802; export HTTPS_PROXY=http://xx.xx.xx.xx:10802;

# 使用 socks5 proxy
export HTTP_PROXY=socks5://xx.xx.xx.xx:10801; export HTTPS_PROXY=socks5://xx.xx.xx.xx:10801;
```

当然，如果你也直接设置成全局代理。
