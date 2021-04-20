## grpc如果使用tls通信，需要对证书进行相关操作，这里使用私有ca签署。

使用openssl工具进行证书操作。

#### 1. 修改openssl配置文件

```bash
#切换到home目录
cd
# 创建操作目录
mkdir ca && cd ca
# 拷贝配置文件到当前目录
cp /etc/pki/tls/openssl.cnf .

#创建空文件
touch /etc/pki/CA/index.txt
#初始化serial文件
echo 00 > /etc/pki/CA/serial
```
修改配置文件
```bash
vi openssl.cnf
#为了方便，[ req_distinguished_name ]段落我修改或添加了一些内容
countryName_default             = CN
stateOrProvinceName_default     = GuangDong
localityName_default            = GuangZhou
0.organizationName_default      = dmai
organizationalUnitName_default  = devops


#[ v3_req ]段落保证有如下内容
subjectAltName = @alt_names

# 文件最后添加如下内容，很重要，需要签发带"使用者备用名称(dns)"的证书用到
[ alt_names ]
DNS.1 = *.grpc.test.com
DNS.2 = dfe.example.org
DNS.3 = ex.abcexpale.net
```




2. 根证书相关操作
##### 2.1. 生成根 CA 的 4096 位长的 RSA 密钥

```bash
# openssl genrsa -out ca.key 4096
Generating RSA private key, 4096 bit long modulus
.....++
........................................++
e is 65537 (0x10001)
```
当前目录生成ca.key文件



```bash
#如果要加密码，需要加选项--aes256
openssl genrsa  -out ca.key 4096 --aes256
```


##### 2.2 生成根证书
使用上一步的私钥ca.key，自签证书ca.crt，有效期100年
```bash
# openssl req -new -x509 -days 36500 -key ca.key -out ca.crt -config openssl.cnf
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [CN]:
State or Province Name (full name) [GuangDong]:
Locality Name (eg, city) [GuangZhou]:
Organization Name (eg, company) [dmai]:
Organizational Unit Name (eg, section) [devops]:
Common Name (eg, your name or your server's hostname) []:ca_server
Email Address []:
```

#### 3. 签发证书

##### 3.1 4096 位长的 RSA 密钥


```bash
# openssl genrsa -out server.key 4096
Generating RSA private key, 4096 bit long modulus
.........++
...........................................++
e is 65537 (0x10001)
```
这一步将生产server.key文件，下面会用到


##### 3.2. 生成证书签发请求
```bash
# openssl req -new -key server.key -out server.csr -config openssl.cnf
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [CN]:
State or Province Name (full name) [GuangDong]:
Locality Name (eg, city) [GuangZhou]:
Organization Name (eg, company) [dmai]:
Organizational Unit Name (eg, section) [devops]:
Common Name (eg, your name or your server's hostname) []:grpc-demo
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```
这一步将生产server.csr文件，下面会用到

##### 3.3. 签发证书

利用ca.crt、server.key、server.csr，生成服务端证书
```bash
# openssl ca -in server.csr -out server.crt -cert ca.crt -keyfile ca.key -extensions v3_req -config openssl.cnf
Using configuration from openssl.cnf
Check that the request matches the signature
Signature ok
Certificate Details:
        Serial Number: 4 (0x4)
        Validity
            Not Before: Apr 20 10:16:17 2021 GMT
            Not After : Apr 20 10:16:17 2022 GMT
        Subject:
            countryName               = CN
            stateOrProvinceName       = GuangDong
            organizationName          = dmai
            organizationalUnitName    = devops
            commonName                = grpc-demo
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:FALSE
            X509v3 Key Usage: 
                Digital Signature, Non Repudiation, Key Encipherment
            X509v3 Subject Alternative Name: 
                DNS:*.grpc.test.com, DNS:dfe.example.org, DNS:ex.abcexpale.net
Certificate is to be certified until Apr 20 10:16:17 2022 GMT (365 days)
Sign the certificate? [y/n]:y


1 out of 1 certificate requests certified, commit? [y/n]y
Write out database with 1 new entries
Data Base Updated
```
这一步生产server.crt证书文件.

查看证书内容,注意证书中的"DNS:*.grpc.test.com, DNS:dfe.example.org, DNS:ex.abcexpale."，这是grpc客户端校验的关键
```bash
# openssl x509 -noout -text  -in server.crt 
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 4 (0x4)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=CN, ST=GuangDong, L=GuangZhou, O=dmai, OU=devops, CN=ca_server
        Validity
            Not Before: Apr 20 10:16:17 2021 GMT
            Not After : Apr 20 10:16:17 2022 GMT
        Subject: C=CN, ST=GuangDong, O=dmai, OU=devops, CN=grpc-demo
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (4096 bit)
                Modulus:
                    00:......(省略)....f7
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:FALSE
            X509v3 Key Usage: 
                Digital Signature, Non Repudiation, Key Encipherment
            X509v3 Subject Alternative Name: 
                DNS:*.grpc.test.com, DNS:dfe.example.org, DNS:ex.abcexpale.net
    Signature Algorithm: sha256WithRSAEncryption
         45:......(省略)....db:2c
```


至此一共产生5个文件+1个配置文件，如下：

```bash
# tree
.
├── ca.crt
├── ca.key
├── openssl.cnf
├── server.crt
├── server.csr
└── server.key
```