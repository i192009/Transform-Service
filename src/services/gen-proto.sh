#!/bin/bash

pwd
files="UserService.proto AccountService.proto PrivilegeService.proto StorageService2.proto TransformService2.proto"
    
for file in $files; do
    echo $file
    protoc --proto_path=../proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional ../proto/$file
done
