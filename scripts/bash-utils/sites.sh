#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/debug.sh
. $(dirname "${BASH_SOURCE[0]}")/console.sh

function generate_tls() {
    local fname=$1

    # 生成私钥 
    openssl genrsa -out ${fname}_PKEY.pem 2048

    # 生成证书签名请求
    openssl req -new -key ${fname}_PKEY.pem -out certrequest.csr 

    # 自签名证书 
    openssl x509 -req -days 365 -in certrequest.csr -signkey ${fname}_PKEY.pem -out ${fname}_CERT.crt

    info "证书信息:" 
    openssl x509 -in ${fname}_CERT.crt -text -noout  

    info "密钥信息:"
    openssl rsa -in ${fname}_PKEY.pem -check

    info "证书和密钥已生成!"
    info "${fname}_CERT.crt 为证书文件" 
    info "${fname}_PKEY.pem 为私钥文件"
}