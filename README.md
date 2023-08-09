This is the minimalist `Tiktok` project of `TikCoder` development team in the sixth back-end youth training camp of ByteDance.

# 1. 项目部署

## docker 安装

首次安装 Docker 之前，需要添加 Docker 安装源。添加之后，我们就可以从已经配置好的源，安装和更新 Docker。添加 Docker 安装源的命令如下：

```bash
yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

正常情况下，直接安装最新版本的 Docker 即可，因为最新版本的 Docker 有着更好的稳定性和安全性。你可以使用以下命令安装最新版本的 Docker。

```bash
yum install docker-ce docker-ce-cli containerd.io
```

安装完成后，使用以下命令启动 Docker。

```Shell
systemctl start docker
```

可以通过下面这个命令查看， docker 是否启动成功。

```Shell
systemctl status docker
```

如果显示了 running，说明 docker 已经启动成功了

## 安装并初始化MySQL

HOST: 本地

宿主机端口: 8086

容器内mysql端口：3306

容器名：mysql

DB_NAME：tiktok

mysql用户：root

mysql密码：123456

用docker启动一个mysql容器

```bash
docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d -e MYSQL_DATABASE=tiktok -p 8086:3306 mysql:8.0
```

查看mysql容器是否启动

```bash
docker ps
```



*目前没有写脚本，没有用 gorm 自动表，需手动执行 init_db.sql*

```mysql
SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
USE
tiktok;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论id，自增主键',
    `user_id`     bigint(20) NOT NULL COMMENT '评论发布用户id',
    `video_id`    bigint(20) NOT NULL COMMENT '评论视频id',
    `content`     varchar(255) NOT NULL COMMENT '评论内容',
    `create_time` datetime     NOT NULL COMMENT '评论发布时间',
    PRIMARY KEY (`id`),
    KEY           `videoIdIdx` (`video_id`) USING BTREE COMMENT '评论列表使用视频id作为索引-方便查看视频下的评论列表'
) ENGINE=InnoDB AUTO_INCREMENT=1206 DEFAULT CHARSET=utf8 COMMENT='评论表';

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `follow_id`   bigint(20) NOT NULL COMMENT '用户id',
    `follower_id` bigint(20) NOT NULL COMMENT '关注的用户',
    PRIMARY KEY (`id`),
    UNIQUE KEY `followIdtoFollowerIdIdx` (`follow_id`,`follower_id`) USING BTREE,
    KEY           `FollowIdIdx` (`follow_id`) USING BTREE,
    KEY           `FollowerIdIdx` (`follower_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1096 DEFAULT CHARSET=utf8 COMMENT='关注表';

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`
(
    `id`       bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`  bigint(20) NOT NULL COMMENT '点赞用户id',
    `video_id` bigint(20) NOT NULL COMMENT '被点赞的视频id',
    PRIMARY KEY (`id`),
    UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
    KEY        `userIdIdx` (`user_id`) USING BTREE,
    KEY        `videoIdx` (`video_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1229 DEFAULT CHARSET=utf8 COMMENT='点赞表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`               bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
    `user_name`        varchar(255) NOT NULL COMMENT '用户名',
    `password`         varchar(255) NOT NULL COMMENT '用户密码',
    `name`             varchar(255) NOT NULL DEFAULT '默认用户名' COMMENT '该用户的默认用户名',
    `work_count`       bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户的作品数量',
    `follow_count`     bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户关注其他用户个数',
    `follower_count`   bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户粉丝个数',
    `total_favorited`  varchar(255) NOT NULL DEFAULT '0' COMMENT '该用户被喜欢的视频数量',
    `favorite_count`   bigint(20) NOT NULL DEFAULT 0 COMMENT '该用户喜欢的视频数量',
    `signature`        varchar(1024) COMMENT '签名',
    `avatar`           varchar(1024) COMMENT '用户头像',
    `background_image` varchar(1024) COMMENT '主页背景',
    PRIMARY KEY (`id`),
    KEY                `name_password_idx` (`user_name`,`password`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20044 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`             bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
    `author_id`      bigint(20) NOT NULL COMMENT '视频作者id',
    `play_url`       varchar(255) NOT NULL COMMENT '播放url',
    `cover_url`      varchar(255) NOT NULL COMMENT '封面url',
    `favorite_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '视频的点赞数量',
    `comment_count`  bigint(20) NOT NULL DEFAULT 0 COMMENT '视频的评论数量',
    `publish_time`   bigint(20) NOT NULL COMMENT '发布时间戳',
    `title`          varchar(255) DEFAULT NULL COMMENT '视频名称',
    PRIMARY KEY (`id`),
    KEY              `time` (`publish_time`) USING BTREE,
    KEY              `author` (`author_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8 COMMENT='视频表';


```

## Minio 安装

先创建对应的目录（项目配置文件中 配置的就是如下目录）

```bash
mkdir -p /root/tiktok/video
mkdir -p /root/tiktok/pic
```

安装 docker

```bash
docker run --name minio \
-p 9000:9000 \
-p 9999:9999 \
-d --restart=always \
-e "MINIO_ACCESS_KEY=admin123" \
-e "MINIO_SECRET_KEY=admin123" \
-v /home/minio/data:/data \
-v /home/minio/config:/root/.minio \
minio/minio server /data \
--console-address '0.0.0.0:9999' -address ":9000"
```

> 在UI界面，浏览器登录 39.106.47.128:9999(自己的IP地址)
>
> ps：配置文件的端口名不要改，仍然是 9000
>
> - 用户名、密码：admin123
> - Access Keys—Create access key
> - 将 Access Key 和 Secret Key 填入配置文件

## ffmpeg 安装

```bash
git clone https://gitee.com/mirrors/ffmpeg.git  
cd ffmpeg/
./configure --disable-x86asm 
make
make install
```

## 项目运行

```go
go run ./main test
```



# 2. 架构设计

## 1. 项目介绍

## 2. 实现功能

### 通用功能

- [x] 密码加密
- [x] JWT鉴权
- [x] 视频ID雪花算法
- [ ] 分布式锁
- [ ] 缓存
- [ ] 日志
- [ ] RPC框架

### 基础功能

- [x] 视频feed流
- [x] 视频投稿
- [x] 个人信息
- [x] 用户登录
- [x] 用户注册
- [ ] 视频投稿

## 3. 技术

| 技术         | go1.20（重新init，可降低版本） |
| ------------ | ------------------------------ |
| http框架     | Gin                            |
| rom框架      | gorm                           |
| 数据库       | MySQL                          |
| 缓存         |                                |
| 加密         | bcrypt                         |
| 对象存储服务 | minio                          |
| 日志管理     |                                |
| 视频         | ffmpeg                         |
| 鉴权         | JWT                            |
| 分布式锁     |                                |
| RPC框架/协议 |                                |

## 4. 框架图

> 客户端请求——>handler（参数校验，JWT，token）——>service——>db

## 5. 数据库表设计

### User表

| 字段             | 数据类型 | 说明                 |
| ---------------- | -------- | -------------------- |
| id               | bigint   | 主键、自增           |
| user_name        | varchar  | 用户名               |
| password         | varchar  | 密码                 |
| name             | varchar  | 用户名称             |
| follow_count     | bigint   | 关注其他用户         |
| follower_count   | bigint   | 粉丝数               |
| total_favorited  | bigint   | 用户被喜欢的视频数量 |
| favorite_count   | bigint   | 用户喜欢的视频数量   |
| signature        | varchar  | 签名                 |
| avatar           | varchar  | 头像                 |
| background_image | varchar  | 主页背景             |

### 视频表

| 字段           | 数据类型 | 说明         |
| -------------- | -------- | ------------ |
| id             | bigint   | 主键、自增   |
| author_id      | bigint   | 视频作者ID   |
| play_url       | varchar  | 播放连接     |
| cover_url      | varchar  | 视频封面连接 |
| favorite_count | bigint   | 点赞数量     |
| comment_count  | bigint   | 评论数量     |
| publish_time   | bigint   | 发布时间     |
| title          | varchar  | 标题         |

### 喜欢表

| 字段     | 数据类型 | 说明       |
| -------- | -------- | ---------- |
| id       | bigint   | 主键、自增 |
| user_id  | bigint   | 用户ID     |
| video_id | bigint   | 视频ID     |

### 评论表

| 字段        | 数据类型 | 说明       |
| ----------- | -------- | ---------- |
| id          | bigint   | 主键、自增 |
| user_id     | bigint   | 用户ID     |
| video_id    | bigint   | 视频ID     |
| content     | varchar  | 评论内容   |
| create_time | varchar  | 评论时间   |

### 关系表

| 字段        | 数据类型 | 说明         |
| ----------- | -------- | ------------ |
| id          | bigint   | 主键、自增   |
| follow_id   | bigint   | 用户ID       |
| follower_id | bigint   | 被关注用户ID |

# 3. 模块设计

## 用户模块

### 用户注册

