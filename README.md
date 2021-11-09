# cDomain
 利用天眼查查询企业备案

[下载地址](https://github.com/canc3s/cDomain/releases)

## 介绍

可以通过两种方式查询自己想要的企业子公司

1. `-n` 参数：利用（http://beian.tianyancha.com）接口给出的关键字先进行查询。（方便，但会搜索结果受关键字影响，假如会出现一些公司名类同的公司的备案）

   ![image-20210301143501471](https://cdn.jsdelivr.net/gh/canc3s/picBed/img/2021/ee1775ac6ca1ba9bd4d126596b2e4707083.png)

2. `-i` 参数：利用给出的公司id对该公司进行查询。（准确，结果唯一，但需要自己先去查找一级公司，该接口目前没有限制不需要 `cookie`）

3. `-f` 参数：对文件里的所有关键字和id进行查询。因为我比较推荐用id查询，而且为了方便多次递归查询，读文件时会先去该行尝试匹配是否存在公司id，假如不存在就把该行作为关键字进行查询。因此查询可以直接把`cSubsidiary`的结果文件当作`cDomain`的输入文件。

4. 因为天眼查风控比较严格，所以使用时会出现几种情况。一、因为某个ip一段时间内查询次数过多，所以查询时会自动跳到登陆界面，这种情况需要使用一个手机号进行登陆，然后增加cookie去继续查询。二、海外ip或者云服务器ip访问天眼查会显示海外用户，所以最好使用正常的出口ip进行查询。三、假如短时间呢使用很多很多次查询的话，天眼查会有人机判断的验证码，需要手动打开天眼查网站进行一下人机验证。（情况较少）

## 用法

```
admin@admin cSubsidiary % go run cDomain.go -h


 ██████╗██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗
██╔════╝██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║
██║     ██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║
██║     ██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║
╚██████╗██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
 ╚═════╝╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝
							v0.0.4
		https://github.com/canc3s/cDomain

Usage of cDomain:
  -c string
    	天眼查的Cookie
  -delay int
    	请求之间的延迟时间(秒)
  -f string
    	包含公司ID号码的文件
  -i string
    	公司ID号码
  -n string
    	公司名称
  -no-color
    	No Color
  -o string
    	结果输出的文件(可选)
  -silent
    	Silent mode
  -timeout int
    	连接超时时间(秒) (default 15)
  -verbose
    	详细模式
  -version
    	显示软件版本号
```

查询子公司

```
admin@admin cSubsidiary % go run cSubsidiary.go -n 字节跳动


 ██████╗██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗
██╔════╝██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║
██║     ██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║
██║     ██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║
╚██████╗██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
 ╚═════╝╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝
							v0.0.4
		https://github.com/canc3s/cDomain

[INFO] 正在查询关键字 字节跳动
[Warning] find IP : 114.246.10.65
[Warning] find IP : 119.167.189.11
dgo.ink
myzijie.com
toutiaocloud.net
toutiaocloud.com
toutiaocloud.cn
bytedance.cn
bytedance.org
bytedance.com
bytedance.net
zjtdchina.cn
bytecdn.cn
bytedns.net
bytedance.tj.cn
shuiwu360.com
byteimg.com
bytefcdn.com
video518.com
jinritoutiao.js.cn
eat9.cn
cdndns1.com
cdndns2.com
bytedns2.com
bytedns1.com
syzjtd.com
mykailu.com
hzzqw.cn
yxlgzs.cn
nmklw.com
shangtout.com
```

## 其他

软件难免有一些问题，假如大家发现，欢迎大家提意见或者建议。

还有一个工具 `cSubsidiary` 我一般两个一起使用，参考[文章](https://canc3s.github.io/2021/03/01/cSubsidiary和cDomain使用指南/)

## Changelog

* 增加请求延迟功能，防止触发反爬虫（-delay 默认为0，不开启）
