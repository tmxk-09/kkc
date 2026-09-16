# UpFile 2.0 (KKC) — 高性能私有化文件极速传输云平台

<div align="center">

![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Vue 3](https://img.shields.io/badge/Frontend-Vue%203%20%7C%20Vite-42b883.svg)
![Go](https://img.shields.io/badge/Backend-Go%201.27%20%7C%20Gin-00ADD8.svg)
![Docker](https://img.shields.io/badge/Deploy-Docker%20%7C%20Linux%20amd64-2496ED.svg)
![Speed](https://img.shields.io/badge/Download-Multi--threaded%20Range%20Accelerated-orange.svg)

**现代化、轻量级、面向千兆带宽设计的极速文件上传下载与空间管理系统**

</div>

---

## 📖 项目简介

**UpFile 2.0** 是一套专为高带宽局域网/互联网场景打造的高性能文件传输系统。全面重构自早期的单体架构，后端采用轻量高效的 **Go 语言**，前端采用 **Vue 3 + Vite**，深度压榨网络带宽瓶颈。

系统内置**文叔叔同款的 Web 端多协程切片并发下载引擎**、**File System Access API 零内存本地流式直写**、**断点续传与秒传**，同时具备完善的**多用户空间配额与空间互赠**功能，支持 Docker 一键离线分发部署。

---

## ⚡ 核心特性

### 1. 🚀 极速多线程下载引擎（并发加速）
- **6 协程 HTTP Range 分片并发**：突破浏览器单连接 TCP 窗口慢启动限制，多通道并行拉取，拉满物理带宽。
- **File System Access API 流式落盘**：优先唤起现代浏览器文件流式句柄，数据边下载边直接按偏移量刷盘，**下载几十 GB 文件内存占用近乎为 0**，杜绝网页卡死。
- **暗黑科技毛玻璃动态看板**：右下角悬浮展示实时传输速率（如 `45.8 MB/s`）、总体完成度以及 6 个并发分片管道跑马灯动画。
- **全场景下载支持**：支持普通单流下载；操作列提供**终端命令一键复制**，即开即得 `aria2c -s 16 -x 16` 与 `curl` 命令，方便 Linux/Mac 终端极速取件。

### 2. 📤 千兆吞吐分片上传与调度
- **并发任务限流调度**：内置前端任务队列，限制单用户同时上传不超过 5 个任务，防止网络拥塞。
- **实时传输变化跟踪**：动态监控每个任务已用时间、传输进度与瞬时上传速率。
- **大文件断点续传 & 极速秒传**：基于文件 MD5 指纹检测，已上传文件瞬间秒传，网络中断可直接从断点继续上传。
- **取件码提取**：支持生成 6 位提取码与自动销毁期限（1天、7天、30天或永久保存）。

### 3. 👥 空间资产化与用户权限系统
- **用户存储配额管控**：管理员可动态调控任意用户的可用空间大小、冻结或解除账户。
- **空间互助赠送机制**：用户间可凭账号将自己的可用存储额度赠送给他人，全程安全密码校验。
- **全员系统消息广播**：支持全员一键群发广播（新注册用户自动补发母版消息）与私信点对点通知。

### 4. 🪶 极简轻量级架构
- 剔除臃肿的传统运行时依赖，采用静态交叉编译二进制与 SQLite 3 (WAL 高并发读写模式)。
- 生产 Docker 容器底座轻盈（镜像仅 **33.8 MB**），整机内存占用仅数十兆。

---

## 🛠️ 技术架构

```text
┌────────────────────────────────────────────────────────┐
│                   Vue 3 + Vite 前端                    │
│   (Element Plus / File System Access API / 响应式看板)  │
└───────────────────────────┬────────────────────────────┘
                            │ HTTP/1.1 长连接 & Range 并发
┌───────────────────────────▼────────────────────────────┐
│                    Nginx 反向代理                      │
│      (长连接连接池 / 缓冲直通 / CORS 暴露切片 Range)      │
└───────────────────────────┬────────────────────────────┘
                            │ 内部高性能通信
┌───────────────────────────▼────────────────────────────┐
│                   Go 后端 (Gin 框架)                    │
│        (JWT 鉴权 / 分片合成 / 空间计算 / 定时清理任务)   │
└───────────────────────────┬────────────────────────────┘
                            │ WAL 高并发读写模式
┌───────────────────────────▼────────────────────────────┐
│                  SQLite 3 嵌入式持久化                  │
└────────────────────────────────────────────────────────┘
```

---

## 📁 目录结构

```text
upfile/
├── .gitignore                     # Git 忽略规则
├── 项目简介.md                    # 本文档
├── upfile-2.0/                    # 2.0 核心服务端与容器化部署
│   ├── server-go/                 # Go 语言高性能后端源码
│   │   ├── handler/               # 接口控制器 (上传/下载/空间/用户等)
│   │   ├── repository/            # SQLite 数据持久层
│   │   ├── service/               # 业务逻辑服务
│   │   └── main.go                # 程序启动入口
│   ├── docker/                    # 部署工具链与脚本
│   │   ├── deploy.sh              # 容器构建与一键部署脚本
│   │   ├── install.sh             # 生产环境一键离线安装脚本
│   │   ├── docker-compose.yml     # 容器编排配置
│   │   └── tune-kernel.sh         # Linux 内核网络与 BBR 加速调优脚本
│   ├── upd/                       # 构建输出的离线包与镜像压缩包 (本地生成)
│   └── 方案.md                    # 业务架构设计文档
└── web/                           # Vue 3 现代前端工程
    ├── src/
    │   ├── views/MainView.vue     # 主工作台（并发下载/上传/管理后台）
    │   └── views/LoginView.vue    # 登录/注册面板
    └── package.json
```

---

## 🚀 快速启动指南

### 1. 前端本地开发
```bash
cd web
npm install
npm run dev
# 本地前端开发服务默认运行在 http://localhost:5173
```

### 2. 后端本地运行
```bash
cd upfile-2.0/server-go
go run main.go
# 服务端默认监听 :9102 端口，支持自动初始化 SQLite 数据库
```

### 3. Docker 容器一键部署
进入 `upfile-2.0/docker` 目录：
```bash
cd upfile-2.0/docker

# 极速编译镜像并启动容器服务
./deploy.sh up
```
- **访问地址**：`http://服务器IP:9110`
- **默认管理员账号**：`admin`
- **默认管理员密码**：见 `config.yaml`（建议首次登录后及时修改）

---

## 🌐 生产服务器网络加速调优建议

为了充分发挥本系统的极速下载潜力，推荐在 Linux 宿主机上开启 **Google BBR 拥塞控制算法**：

```bash
# 执行网络与内核调优脚本 (需 root 权限)
sudo bash upfile-2.0/docker/tune-kernel.sh

# 开启 BBR 拥塞控制
sudo sysctl -w net.core.default_qdisc=fq
sudo sysctl -w net.ipv4.tcp_congestion_control=bbr

# 验证 BBR 开启状态
sysctl net.ipv4.tcp_congestion_control
# 预期输出: net.ipv4.tcp_congestion_control = bbr
```

---

## 📄 开源许可证

本项目基于 [MIT License](LICENSE) 开源。
