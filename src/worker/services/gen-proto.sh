#!/bin/bash

pwd
files="UserService.proto PrivilegeService.proto StorageService2.proto TransformService2.proto"
    
for file in $files; do
    echo $file
    python -m grpc_tools.protoc -I ../../proto --python_out=. --grpc_python_out=. --experimental_allow_proto3_optional ../../proto/$file
    # --proto_path=../grpc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional ../proto/$file
done
