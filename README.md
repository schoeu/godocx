# godocx

> 简单快速的markdown文件浏览工具

## 安装

```
go install
```

## 启动


```
go run main.go "配置文件路径"
```

## 说明

* 配置中的路径均支持相对路径与绝对路径。
* 配置中只有`path`字段为必须,其他字段都是可选。
* 该文档平台有独立日志,日志路径可配置

## 配置参数

这里列出全部参数,除了`path`参数外,其他参数按需添加即可。

```
{
  // 监听端口,默认为8910,可选
  "port": "8910",

  // markdown文档路径,支持相对,绝对路径,必选
  "path": "/home/work/docx",

  // 需要忽略的目录名,不能被markdown正确解析的目录都应该加到这里来,可选
  "ignoreDir": ["img",".git",".svn"],

  // 是否debug状态, 非debug状态会启用缓存,可选
  "debug": true,

  // header条标题,可选
  "headText": "PSFE-DOC",

  // 展示主题,可选, 开箱自带两套皮肤default,antd,默认为default.
  "theme": "default",

  // 预处理脚本定制,填写脚本地址即可,可选
  "preprocessscript":"",

  // web title,可选
  "title": "PSFE",

  // 默认文档路径,支持相对,绝对路径,可选
  "index": "/readme.md",

  // 技术支持,可选
  // 邮箱填写: mailto:xx@xxx.com
  // Hi填写: baidu://message/?id=用户名,可以直接调起Hi
  "supportInfo": "baidu://:xx@xxx.com",

  // 默认false, 开启报警后报错会发送邮件,可选
  "waringFlag": false,

  // 报警邮箱配置,可选
  "warningEmail": {
    "host": "smtp.163.com",
    "port": 25,
    "user": "xx@163.com",
    "pass": "password",
    "from": "xx@163.com", // 发件人
    "to": "xx@xxx.com", // 收件人
    "subject": "DOCX error"  // 邮件标题
  },

  // 文件夹命名配置文件路径,可选
  "dirsConfName": "map.json",

  // 链接配置,展示位置为右上角,可以配置其他链接,可选
  "extUrls": {
    "label": "友情链接",
    "links":[{
      "name": "栅格文档PMD",
      "url": "http://sfe.baidu.com/pmd/doc/"
    },{
      "name": "MIP文档",
      "url": "http://mip.baidu.com"
    },{
      "name": "SUPERFRAME文档",
      "url": "http://superframe.baidu.com"
    }]
  },
  // 文件夹命名配置
  "dirNames":[
         {"dir1": {
           "name": "dir1",
           "sort": 1
         }},
         {"dir2": {
           "name": "dir2",
           "sort": 2
         }},
         {"dir3": {
           "name": "dir3"
         }}
     ]
}

```
## 主题

开箱自带两套皮肤`default`,`antd`,默认为`default`主题。
目录为`themes/default`与`themes/antd`,如想换其他主题请自行替换或开发。
