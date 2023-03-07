# GZAuth - Simple Your Auth for Web Service

gzauth是一个用于服务访问控制的项目，支持多种认证方式，包括Basic Auth、BearToken、GitHub登录和飞书登录等。它可以使用Docker进行部署，并支持私有化部署。此外，gzauth完全开源。

## 特性

- [x] 支持 Basic Auth 认证
- [x] 支持 BearToken 认证
- [x] 支持 GitHub 登录
- [x] 支持飞书登录
- [x] 使用Docker容器化部署
- [x] 支持私有化部署

## 安装

### 方案一、Docker Compose 部署（推荐）

1. 创建 `docker-compose.yml`，这里使用 Basic Auth 类型:

```yaml
# 使用 basic auth
services:
  gzauth:
    restart: unless-stopped
    image: whatwewant/gzauth:latest
    ports:
      - 8080:8080
    environment:
      # basic auth
      AUTH_TYPE: basicauth
      USERNAME: <YOUR_USERNAME>
      PASSWORD: <YOUR_PASSWORD>
      UPSTREAM: https://httpbin.org
```

替换上述环境变量的值为您自己的值。

2. 启动容器：

```bash
$ docker-compose up -d
```

### 方案二、二进制部署

```bash
# 安装服务器管理框架 Zmicro
$ curl -o- https://raw.githubusercontent.com/zcorky/zmicro/master/install | bash

# 安装 gzauth
$ zmicro package install gzauth

# 运行
$ zmicro gzauth basicauth --username <YOUR_USERNAME> --password <YOUR_PASSWORD> --upstream <YOUR_WEBSERVICE>
```


## 更多案例
* 1. 使用 `BearerToken`

```yaml
# docker-compose.yml
services:
  gzauth:
    restart: unless-stopped
    image: whatwewant/gzauth:latest
    ports:
      - 8080:8080
    environment:
      AUTH_TYPE: bearertoken
      TOKEN: <YOUR_TOKEN>
      UPSTREAM: https://httpbin.org
```

* 2. 使用 `GitHub 登录`

```yaml
# docker-compose.yml
services:
  gzauth:
    restart: unless-stopped
    image: whatwewant/gzauth:latest
    ports:
      - 8080:8080
    environment:
      AUTH_TYPE: github
      CLIENT_ID: <GITHUB_OAUTH2_CLIENT_ID>
      CLIENT_SECRET: <GITHUB_OAUTH2_CLIENT_SECRET>
      REDIRECT_URI: <GITHUB_OAUTH2_REDIRECT_URI>
      UPSTREAM: https://httpbin.org
```

* 2. 使用 `飞书登录`

```yaml
# docker-compose.yml
services:
  gzauth:
    restart: unless-stopped
    image: whatwewant/gzauth:latest
    ports:
      - 8080:8080
    environment:
      AUTH_TYPE: feishu
      CLIENT_ID: <FEISHU_OAUTH2_CLIENT_ID>
      CLIENT_SECRET: <FEISHU_OAUTH2_CLIENT_SECRET>
      REDIRECT_URI: <FEISHU_OAUTH2_REDIRECT_URI>
      UPSTREAM: https://httpbin.org
```

## 使用

请参阅 [USAGE.md](./USAGE.md) 文件了解如何使用gzauth。

## 贡献

欢迎您参与贡献gzauth！请参阅 [CONTRIBUTING.md](./CONTRIBUTING.md) 文件了解更多信息。

## 许可证

gzauth采用MIT许可证。请参阅 [LICENSE.md](./LICENSE.md) 文件了解详细信息。
