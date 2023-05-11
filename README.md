#  sage

#### 一款自动将mysql表结构（不含数据）导出成csv，从而可以导入到word的效率工具。
#### 由于某些企业流程过于死板，明明有了model还需要在概要设计（word文档）中将表结构整理成表格，整个过程十分繁琐，急需一款效率工具。在经过调研之后没有找到合适的工具，有的也是收费的而且巨贵，于是花了三个小时自己写了一段代码，代码不多，很简单，但是实用。
#### 由于只是暂时解决手边的问题，因此直接拷了一个现有的项目来改，项目是基于kratos的，有感兴趣的朋友可以把工程修改修改，改成一个命令行工具。

#### 另外，某些公司真的是神烦，明明有openapi文档了还是需要在概要设计（word）中将文档整理成表格，无奈，我之后会抽时间写一个转换工具，有兴趣的朋友可以把代码贡献在这个项目里

## 效果展示
| 字段 |类型|是否允许为空|默认值|备注|
| ---- | ---- | ---- | ---- | ---- |
|id	|char(36)	|否	|	|业务唯一Id|
|created_at|	datetime|	否	|	|创建时间|
|updated_at	|datetime	|否|		|更新时间|
|city	|varchar(128)	|是|	|	|city|
|country	|varchar(128)|是|		|country|
|session_key|	varchar(63)|	是	|	|sessionKey|
|app_code	|char(5)	|否	|1|	|应用code|
|login_at	|datetime|	是|		|本次登录时间|

## 安装CLI工具

- 安装kratos

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

- [安装protoc](http://google.github.io/proto-lens/installing-protoc.html)

## 运行项目

```bash
# 1.项目初始化
make build

# 2.生成代码
./bin/server

# 3.修改配置

vim configs/config.yaml

```


