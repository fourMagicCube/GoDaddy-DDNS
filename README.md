# godaddy动态更新域名 #

----------
#### 1. 从 https://developer.godaddy.com/keys 申请一个Production的API KEY。ote表示测试用，Production表示正式生产环境使用。 ####
#### 2. 进入main.go ####
#### 3. 自定义信息 ####
```
 var domain = ""       //域名：不包括www.

 var key = ""          //自己从godaddy创建的key

 var secret = ""       //自己从godaddy创建的secret

 var logPath="logs/"   //日志地址 (后面要加/)

 var name = "AAAA"     //ipv4:A; ipv6:AAAA
```
    



#### 4. 构建go程序 ####
    $ go build main.go

#### 5. 运行程序  ####
    $ ./main &
#### 6. 日志查看  ####
    $ tail -f logs/info-xxxx-xx-xx.log