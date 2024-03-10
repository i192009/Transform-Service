#Design goals
The primary design goals for Transform Service Version 2 are as follows:
• Seamless integration with Temporal for efficient task scheduling.
• A comprehensive business process ensures that file conversion tasks always report accurate
task conversion status and progress under any circumstances, without any stuck situations.
• It supports up to 100 W+ concurrent tasks, and the single conversion time of a model file of
about 1G does not exceed 10 minutes, less than 100M does not exceed 5 minutes.
• The number of jobs to be executed in parallel can be determined based on the computing
resources of the worker, and optimization parameters or scripts can be specified separately
for each file.
• Having multiple task scheduling capabilities, the system is only responsible for resource
scheduling and does not care about task parameters and tasks
• A clearer error return distinguishes between file conversion failures caused by insufficient
resources, long queues, incorrect file format, and file conversion errors. Provide waiting time
prediction function, which can obtain the time required for file conversion at any time
(queue time, queue location, conversion time)
• The ability to automatically recover from abnormal situations. When the Transform service
or Worker process encounters an exception, the program can automatically recover from the
exception for a period of time, including but not limited to process crashes, database
connection failures, and file conversion interruptions caused by insufficient resources.
• Has the ability to manage preset optimization parameters and scripts, and can specify which
preset optimization parameters and scripts the enterprise can use.
• Based on the payment situation of the organization, support the use of priority queues and
exclusive resource pools by the organization. VIP organization have the priority to execute
tasks or configure exclusive computing resource pools for the organization, allowing the
organization to start and stop computing resources on its own and charge on time.


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

