# transform v2

## IDE
Use VS Code as the standard ID. To open the project, please use VS Code to open transfrom.code-workspace.

## Project Structure

```txt
*
|-.vscode                     # vs code configuration folder
  |-launch.json               # debug launcher configuration
  |-tasks.json                # task configuration
|-bin
  |-conf
    |-transform.yaml.sample   # a sample for transform2 configuration
|-scripts
  |-bash-utils                # Submodule, there has a lot of very useful bash functions, such as Menu install
  |-cases.json                # Test case for grpc
  |-grpc.sh                   # A bash script for executing test case, Interact using text UI
  |-init-python.sh            # A bash script for installing python and other dependent packages
  |-temporal.sh               # A bash script can start or close a local temporal service.
|-src                         # transform v2 source code
  |-abandoned                 # abandoned
  |-config                    # transform config variant and other const variant
  |-grpcserver                # grpc interfaces are implemented here
  |-monitor                   # Monitoring package, a standalone command line tool for monitoring the health of worker processes
  |-proto                     # Submodule for git@gitlab.zixel.cn:z-jumeaux-engine/services/framework/grpc.git
  |-services                  # Protoc generated files for go
  |-web                       # web interfaces are implemented here
  |-worker                    # Each folder is an independent worker for different requirements
    |-job-scheduler           # Job scheduling module, used to decide how to process jobs.
    |-res-scheduler           # Resource scheduling module, used to decide when and how to start a new server
    |-services                # Protoc generated files for python
    |-transformer             # Worker for convert files by hoops command line
    |-zcad                    # Worker for ZCAD 
    |-install.sh              # Install
|-framework                   # Submodule for zixel golang micro service framework
```

> .vscode/tasks.json  
can be called with hotkey ctrl+shift+b in vscode default keymap setting. currently configured three task, generate grpc protocol files for golang, generate grpc protocol files for golang and build transform2 to folder /bin 

## How to use

### How to build a debug environment?
1. use "ktctl connect --namespace dev" connect to dev environment
2. run "./temporal.sh start" to open local temporal service
3. press F5 to launch transform2 using debuger.
4. use "go run src/worker/xxx/main.go" to start worker
5. run "./grpctest.sh" do grpc interface test.

### How to add new test case?
1. open cases.json
2. add new object, for the parameter field, the better way is using base64 encode.
3. edit grpctest.sh 
    - add new MenuItem
    - add new Option in case
    - add grpc command under the new option

