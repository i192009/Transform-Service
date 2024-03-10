#!/bin/bash
. $(dirname "${BASH_SOURCE[0]}")/debug.sh

tmp=($(whereis -b kubectl))
trace "${tmp[@]}"
if [[ $? == 0 && ${#tmp[@]} > 1 ]]; then
    kubectl=${tmp[1]}
else
    ReadParameter kubectl k
fi
trace "$kubectl"

# 递归执行目录下的部署脚本，若目录下有.applyfiles文件则按该文件中的文件顺序执行。
function ApplyFiles() {
    trace "k8s apply paths=\($@\)"
    for path in $@; do
        if [ -d $path ]; then
            if [ -e $path/.applyfiles ]; then
                ApplyFiles $(sed "s/^/${path//\//\\\/}\/&/g" $path/.applyfiles)
            else
                ApplyFiles $(find $1 | grep .yaml)
            fi
        else
            eval "echo \"$(FB 22)$(sed 's/"/\\"/g' $path)\"$(FB)"
            read -n1 -p "apply this file?[$path](Y/n)" answer
            echo #空一行
            if [[ "$answer" == "" || "$answer" == "y" || "$answer" == "Y" ]]; then
                eval "echo \"$(sed 's/"/\\"/g' $path)\"" | $kubectl apply -f -
            fi
        fi
    done
}

# 部署目录下的所有服务
function DeployService() {
    debug enable
    for path in $@; do
        if [ -d $path ]; then
            trace "current path=$path"
            # 如果是目录
            if [ -e "$path/.deployfiles" ]; then
                trace "exists .deployfiles $(sed -e "/^[ ]*$/d" -e "s/^/${path//\//\\\/}\/&/g" $path/.deployfiles)"
                # 如果存在.deployfiles文件，则按文件中的顺序执行
                DeployService $(sed -e "/^[ ]*$/d" -e "s/^/${path//\//\\\/}\/&/g" $path/.deployfiles)
            else
                trace "find all deploy files"
                # 否则按目录下的文件顺序执行
                DeployService $(find $path | grep ".deploy$" | sed -e "s/.deploy$//g")
            fi
        else
            # 如果是部署文件，则执行
            trace "deploy file=$path.deploy"
            DeployServiceReset
            . "$path.deploy"

            DeployServiceTemplate
            read -n1 -p "deploy file?[$path](Y/n)" answer
            echo #空一行
            if [[ "$answer" == "" || "$answer" == "y" || "$answer" == "Y" ]]; then
                debug disable
                DeployServiceTemplate > $path-deployment.yaml
                $kubectl apply -f $path-deployment.yaml
                debug resume
                rm $path-deployment.yaml
            fi
            trace "deploy file done."
        fi
    done
    debug disable
}

function DeployServiceReset() {
    trace "reset deploy parameters ..."

    # application setting
    # 应用名
    APP_NAME=
    # 应用日志路径
    APP_LOG_PATH=/logs/$APP_NAME
    # 应用配置路径
    APP_CONF_PATH="/conf/application.yaml"
    # 应用启动的GRPC监听端口
    APP_GRPC_PORT=9400
    # 应用启动的HTTP监听端口
    APP_HTTP_PORT=80
    # 应用启动的其它端口
    APP_PORTS=
    # 健康检查接口，为空则不设置，默认为空值
    APP_LIVENESS_PATH=
    # 服务可用性检测接口，为空则不设置，默认为空值
    APP_READINESS_PATH=
    # 服务启动接口，为空则不设置，默认为空值
    APP_STARTUP_PATH=

    # config
    ## configMap
    # 是否开启configmap配置，默认不开启
    K8S_CONFIG_MAP_ENABLE=false
    # 是否加密configmap
    K8S_CONFIG_MAP_ENCRYPT=$K8S_CONFIG_MAP_DEFAULT_ENCRYPT

    # configMap 映射，key为映射路径，value为configmap的key，若value为空则将configMap的所有key映射到指定路径下
    unset K8S_CONFIG_MAP

    ## storage
    # 是否开启存储卷，默认不开启
    K8S_STORAGE_ENABLE=false
    # 存储卷名称，多个存储卷用空格分隔
    unset K8S_STORAGE_NAMES
    # 存储卷路径，key为存储卷名称，value为存储卷路径，此为必填项
    unset K8S_STORAGE_PATH
    # 存储卷大小，key为存储卷名称，value为存储卷大小，此为必填项
    unset K8S_STORAGE_SIZE
    # 存储卷类型，key为存储卷名称，value为存储卷类型，此为必填项
    unset K8S_STORAGE_MODE

    ## service
    # 是否为服务创建service，默认不开启
    K8S_SERVICE_ENABLE=false
    # 服务类型，支持ClusterIP、NodePort、LoadBalancer，默认为ClusterIP
    K8S_SERVICE_TYPE=ClusterIP
    # 服务开启的GRPC端口
    K8S_SERVICE_GRPC_PORT=9400
    # 服务开启的HTTP端口
    K8S_SERVICE_HTTP_PORT=80
    # 服务开启的NodePort端口，仅当K8S_SERVICE_TYPE为NodePort时生效
    K8S_SERVICE_NODE_PORT=30080
    # 服务开启的其它端口，仅当K8S_SERVICE_EXTRA_PORT有值时生效
    unset K8S_SERVICE_EXTRA_TYPE
    # 开启的服务端口
    unset K8S_SERVICE_EXTRA_PORT
    # 开启的服务端口协议
    unset K8S_SERVICE_EXTRA_PROTOCOL
    # 开启的服务端口对应的NodePort端口，仅当K8S_SERVICE_EXTRA_TYPE为NodePort时生效
    unset K8S_SERVICE_EXTRA_NODEPORT

    ## deployment
    # 是否为服务创建deployment，默认不开启
    K8S_DEPLOYMENT_ENABLE=false
    # 是否开启Kuboard属性
    K8S_DEPLOYMENT_KUBOARD=disable
    # Kuboard 分层
    K8S_DEPLOYMENT_KUBOARD_LAYER=
    # 部署名称
    K8S_DEPLOYMENT_NAME=$APP_NAME
    # CPU限制
    K8S_DEPLOYMENT_CPU_LIMIT=500m
    # 内存限制
    K8S_DEPLOYMENT_MEM_LIMIT=512Mi
    # CPU请求
    K8S_DEPLOYMENT_CPU_REQUEST=250m
    # 内存请求
    K8S_DEPLOYMENT_MEM_REQUEST=128Mi
    # 镜像地址
    K8S_DEPLOYMENT_IMAGE_URL=
    # 镜像版本
    K8S_DEPLOYMENT_IMAGE_TAG=
    # 镜像仓库凭证
    K8S_DEPLOYMENT_IMAGE_SECRET=$K8S_DEPLOYMENT_DEFAULT_IMAGE_SECRET
    # 镜像拉取策略
    K8S_DEPLOYMENT_IMAGE_POLICY=$K8S_DEPLOYMENT_DEFAULT_IMAGE_POLICY
    # 镜像的工作目录
    K8S_DEPLOYMENT_WORKDIR=
    # 镜像的命令
    K8S_DEPLOYMENT_COMMAND=
    # 镜像启动参数
    K8S_DEPLOYMENT_ARGS=
    # 部署的服务账号
    K8S_DEPLOYMENT_SERVICEACCOUNT=config-reader
    # 默认副本数
    K8S_DEPLOYMENT_REPLICA_COUNT=1
    # 配置存储卷，key为映射的路径，value为映射的configmap的key。根据命名规则configmap的名字为$APP_NAME-conf
    unset K8S_DEPLOYMENT_CONFIG_MAPS
    # 配置存储卷，key为"卷名称"，value为"容器内路径:存储卷名称:存储卷路径"，其中存储卷路径可以省略
    unset K8S_DEPLOYMENT_VOLUME_KEYS

    trace "reset deploy complete"
}

function DeployServiceTemplate() {
    trace "generate deploy template K8S_CONFIG_MAP_ENABLE=$K8S_CONFIG_MAP_ENABLE ..."
    if [ "$K8S_CONFIG_MAP_ENABLE" == "true" ]; then
        trace "generate configMap ..."
        echo "---"
        echo "kind: ConfigMap"
        echo "apiVersion: v1"
        echo "metadata:"
        echo "  labels:"
        echo "    name: $APP_NAME"
        echo "  name: $APP_NAME-conf"
        echo "  namespace: $K8S_NAMESPACE"
        echo "data:"
        for CONFIG_KEY in ${!K8S_CONFIG_MAP[@]}; do
        trace "generate configMap key = $CONFIG_KEY ..."
        ConfigMap=${K8S_CONFIG_MAP[$CONFIG_KEY]}
        if [[ ! -f $ConfigMap ]]; then
            # 不是文件则跳过
            trace "${ConfigMap} is not file"
            continue;
        fi
        trace "K8S_CONFIG_MAP_ENCRYPT=$K8S_CONFIG_MAP_ENCRYPT"
        # 如果开启了加密，则对非env文件进行加密
        if [ "$K8S_CONFIG_MAP_ENCRYPT" == "true" ] && [ ${ConfigMap##*.} != "env" ]; then
            # 执行configMap加密
            yaml=${ConfigMap##*/}
            encrypt_config_map $yaml
            ConfigMap=./$ENCRYPT_CONFIG_DIR/$yaml
        fi
        echo "  $CONFIG_KEY: |"
        eval "echo \"`sed -e 's/^/    /g' -e 's/\"/\\\\\\\"/g' $ConfigMap`\""
        done
    fi

    trace "K8S_STORAGE_ENABLE=$K8S_STORAGE_ENABLE"
    if [ "$K8S_STORAGE_ENABLE" == "true" ]; then
    for STORAGE_NAME in ${K8S_STORAGE_NAMES[@]}; do
        ValidateParameters \
            'K8S_STORAGE_NAME[$STORAGE_NAME]'\
            'K8S_STORAGE_PATH[$STORAGE_NAME]'\
            'K8S_STORAGE_SIZE[$STORAGE_NAME]'\
            'K8S_STORAGE_MODE[$STORAGE_NAME]'
        
        if [ $? != 0 ]; then
            return 1
        fi
        echo "---"
        echo "kind: PersistentVolumeClaim"
        echo "apiVersion: v1"
        echo "metadata:"
        echo "  name: $STORAGE_NAME"
        echo "  namespace: $K8S_NAMESPACE"
        echo "  annotations:"
        echo "    nfs.io/storage-path: ${K8S_STORAGE_PATH[$STORAGE_NAME]}"
        echo "spec:"
        echo "  storageClassName: nfs-client"
        echo "  accessModes:"
        echo "    - ${K8S_STORAGE_MODE[$STORAGE_NAME]}"
        echo "  resources:"
        echo "    requests:"
        echo "      storage: ${K8S_STORAGE_SIZE[$STORAGE_NAME]}"
    done
    fi

    trace "K8S_DEPLOYMENT_ENABLE=$K8S_DEPLOYMENT_ENABLE"
    if [ "$K8S_DEPLOYMENT_ENABLE" == "true" ]; then
        ValidateParameters \
            K8S_DEPLOYMENT_NAME\
            K8S_DEPLOYMENT_CPU_LIMIT\
            K8S_DEPLOYMENT_MEM_LIMIT\
            K8S_DEPLOYMENT_CPU_REQUEST\
            K8S_DEPLOYMENT_MEM_REQUEST\
            K8S_DEPLOYMENT_IMAGE_URL\
            K8S_DEPLOYMENT_IMAGE_TAG\
            K8S_DEPLOYMENT_REPLICA_COUNT
        
        if [ $? != 0 ]; then
            return 1
        fi
        echo "---"
        echo "#版本号，可以通过 kubectl api-versions 命令进行获取当前支持的版本"
        echo "apiVersion: apps/v1"
        echo "#Kubernetes资源类型，部署应用类型，为ReplicaSet和Pod的创建提供一种声明式的定义方法"
        echo "#无需手动创建ReplicaSet/Replication Controller和Pod，使用Deployment支持滚动升级和回滚等特性"
        echo "kind: Deployment"
        echo "#该资源元数据"
        echo "metadata:"
        echo "  #该资源名称，决定了Pod的显示名称"
        echo "  name: $K8S_DEPLOYMENT_NAME"
        echo "  #该资源所属的命名空间"
        echo "  namespace: $K8S_NAMESPACE"
        echo "  # 自定义的标签列表,key-value键值对"
        echo "  labels:"
        if [ "$K8S_DEPLOYMENT_KUBOARD" == "enable" ]; then
        echo "    k8s.kuboard.cn/layer: $K8S_DEPLOYMENT_KUBOARD_LAYER"
        echo "    k8s.kuboard.cn/name: $K8S_DEPLOYMENT_NAME"
        fi
        echo "    name: $K8S_DEPLOYMENT_NAME"
        echo "#该资源的详细定义"
        echo "spec:"
        echo "  # 指明哪个pod被管理，这里我们指定了name为{{$K8S_DEPLOYMENT_NAME}}"
        echo "  selector:"
        echo "    matchLabels:"
        echo "      name: $K8S_DEPLOYMENT_NAME"
        echo "  # 指定运行的Pod数量，它会根据selector来进行选择，会将选择到的Pod维持在{{.replicas}}个的数目量"
        echo "  replicas: $K8S_DEPLOYMENT_REPLICA_COUNT"
        echo "  # 扩容时创建Pod的模板"
        echo "  template:"
        echo "    #元数据"
        echo "    metadata:"
        echo "      #标签列表"
        echo "      labels:"
        echo "        name: $K8S_DEPLOYMENT_NAME"
        echo "    spec:"

        if [ -n "$K8S_DEPLOYMENT_SERVICEACCOUNT" ]; then
        echo "      serviceAccount: $K8S_DEPLOYMENT_SERVICEACCOUNT"
        fi

        echo "      # pod中的容器列表"
        echo "      containers:"
        echo "        # 容器名称"
        echo "        - name: $K8S_DEPLOYMENT_NAME"
        echo "          # 容器镜像名称"
        echo "          image: $K8S_DEPLOYMENT_IMAGE_URL:$K8S_DEPLOYMENT_IMAGE_TAG"
        if [ -n "$K8S_DEPLOYMENT_COMMAND" ]; then
        echo "          command: [$K8S_DEPLOYMENT_COMMAND]"
        fi
        if [ -n "$K8S_DEPLOYMENT_ARGS" ]; then
        echo "          args: [$K8S_DEPLOYMENT_ARGS]"
        fi
        if [ -n "$K8S_DEPLOYMENT_WORKDIR" ]; then
        echo "          workingDir: $K8S_DEPLOYMENT_WORKDIR"
        fi
        echo "          # 镜像抓取的策略，其中有Always,IfNotPresent,Never"
        echo "          # Always=>每次都会重新下载镜像"
        echo "          # IfNotPresent=>如果本地又镜像则使用，否则下载镜像"
        echo "          # Never=>仅使用本地镜像"
        echo "          imagePullPolicy: $K8S_DEPLOYMENT_IMAGE_POLICY"
        echo "          # 容器需要暴露的端口号列表,注意凡是监听0.0.0.0地址的端口，即使没有设置containerPort参数，都可以被访问"
        echo "          ports:"

        # 暴露HTTP端口
        if [ -n "$APP_HTTP_PORT" ]; then
        echo "            # 端口的名称"
        echo "            - name: http"
        echo "              # 容器需要监听的端口"
        echo "              containerPort: $APP_HTTP_PORT"
        fi

        # 暴露GRPC端口
        if [ -n "$APP_GRPC_PORT" ]; then
        echo "            - name: grpc"
        echo "              containerPort: $APP_GRPC_PORT"
        fi

        if [ -n "${APP_PORTS[@]}" ]; then
        for APP_PORT in ${APP_PORTS[@]}; do
        echo "            - name: port-$APP_PORT"
        echo "              containerPort: $APP_PORT"
        done
        fi
        echo "          # 资源限制和资源请求的设置"
        echo "          resources:"
        echo "            # 通常我们会把request设置为一个比较小的数值，这个值根据项目情况而定，需求限制，也叫软限制"
        echo "            requests:"
        echo "              cpu: $K8S_DEPLOYMENT_CPU_REQUEST"
        echo "              memory: $K8S_DEPLOYMENT_MEM_REQUEST"
        echo "            #最大限制，也叫硬限制，通常requests和limits要一起配置，避免资源使用不合理出现问题"
        echo "            limits:"
        echo "              cpu: $K8S_DEPLOYMENT_CPU_LIMIT"
        echo "              memory: $K8S_DEPLOYMENT_MEM_LIMIT"
        echo "          env:"
        echo "            - name: PAAS_NAMESPACE  # bootstrap.xml中配置的configmap的namespace"
        echo "              valueFrom:"
        echo "                fieldRef:"
        echo "                  fieldPath: metadata.namespace"
        echo "            - name: CONFIGMAP_NAMESPACE  # bootstrap.xml中配置的configmap的namespace"
        echo "              valueFrom:"
        echo "                fieldRef:"
        echo "                  fieldPath: metadata.namespace"
        echo "          # 容器挂载路径"
        echo "          volumeMounts:"

        trace "K8S_CONFIG_MAP_ENABLE=$K8S_CONFIG_MAP_ENABLE"

        if [ "$K8S_CONFIG_MAP_ENABLE" == "true" ]; then
        for CONFIG_MAP in ${K8S_DEPLOYMENT_CONFIG_MAPS[@]}; do
        echo "            - name: conf"
        local Mpath=`echo $CONFIG_MAP | awk -F ':' '{printf $1}'`
        local Spath=`echo $CONFIG_MAP | awk -F ':' '{printf $2}'`
        echo "              mountPath: $Mpath # 映射路径"
        if [ -n "$Spath" ]; then
        echo "              subPath: $Spath"
        fi
        done
        fi

        if [ "$K8S_STORAGE_ENABLE" == "true" ]; then
        for VOLUME_KEY in ${!K8S_DEPLOYMENT_VOLUME_KEYS[@]}; do
        local Mpath=`echo ${K8S_DEPLOYMENT_VOLUME_KEYS[$VOLUME_KEY]} | awk -F ':' '{printf $1}'`
        if [ -n "$Mpath" ]; then
        echo "            - name: $VOLUME_KEY"
        echo "              mountPath: $Mpath # 在容器中的路径"
        local Spath=`echo ${K8S_DEPLOYMENT_VOLUME_KEYS[$VOLUME_KEY]} | awk -F ':' '{printf $3}'`
        if [ -n "$Spath" ]; then
        echo "              subPath: $Spath # 在存储卷中的路径"
        fi
        fi
        done
        fi
        echo "            - name: localtime"
        echo "              mountPath: /etc/localtime"


        if [ -n "$APP_LIVENESS_PATH" ]; then
        echo "          # 对Pod内容器健康检查的设置"
        echo "          # 存活探针livenessProbe: 用于判断容器是否存活(running)，如果livessProbe探针探测到不健康，则kubelet会杀掉该容器，并根据容器的重启策略做相应的处理，不配的话K8S会认为一直是Success"
        echo "          # 判断是否存活，不存活则重启容器"
        echo "          # 其中有exec、httpGet、tcpSocket这3种配置"
        echo "          # exec: 在容器内部执行一个命令，如果该命令返回码为0，则表明容器健康"
        echo "          # httpGet: 通过容器的IP地址、端口号及路径通过用httpGet方法，如果响应的状态码stateCode大于等于200且小于400则认为容器健康"
        echo "          # tcpSocket: 通过容器的IP地址和端口执行tcpSocket检查，如果能够建立tcpSocket连接，则表明容器健康"
        echo "          livenessProbe:"
        echo "            httpGet:"
        echo "              # 健康检查路径，也就是应用的健康检查路径"
        echo "              path: $APP_LIVENESS_PATH"
        echo "              # 访问的容器的端口名字或者端口号"
        echo "              port: $APP_HTTP_PORT"
        echo "              scheme: HTTP"
        echo "            # 容器启动后开始探测之前需要等多少秒，如应用启动一般30s的话，就设置为 30s，设置启动探针startProbe后可不等待"
        echo "#            initialDelaySeconds: 60"
        echo "            # 探测的超时时间，默认 1s，最小 1s"
        echo "            timeoutSeconds: 3"
        echo "            # 健康检查失败后，最少连续健康检查成功多少次才被认定成功，默认值为1。最小值为1。"
        echo "            successThreshold: 1"
        echo "            # 最少连续多少次失败才视为失败。默认值为3。最小值为1"
        echo "            failureThreshold: 3"
        echo "            # 执行探测的频率（多少秒执行一次）。默认为10秒。最小值为1"
        echo "            periodSeconds: 10"
        echo "          # 就绪探针redinessProbe: 用于判断容器是否启动完成（ready），可以接收请求，如果该探针探测到失败，则Pod状态会被修改。Endpoint Controller将从ServiceEndpoint中删除包含该容器的所在Pod的Endpoint"
        echo "          # kubernetes 判断容器是否启动成功,否可以接受外部流量"
        fi

        if [ -n "$APP_READINESS_PATH" ]; then
        echo "          readinessProbe:"
        echo "            httpGet:"
        echo "              # 健康检查路径，也就是应用的健康检查路径"
        echo "              path: $APP_READINESS_PATH"
        echo "              # 端口号"
        echo "              port: $APP_HTTP_PORT"
        echo "              scheme: HTTP"
        echo "            # 容器启动后开始探测之前需要等多少秒，如应用启动一般30s的话，就设置为 30s，设置启动探针startProbe后可不等待"
        echo "            # initialDelaySeconds: 60"
        echo "            # 探测的超时时间，默认 1s，最小 1s"
        echo "            timeoutSeconds: 3"
        echo "            # 健康检查失败后，最少连续健康检查成功多少次才被认定成功，默认值为1。最小值为1。"
        echo "            successThreshold: 1"
        echo "            # 最少连续多少次失败才视为失败。默认值为3。最小值为1"
        echo "            failureThreshold: 3"
        echo "            # 执行探测的频率（多少秒执行一次）。默认为10秒。最小值为1"
        echo "            periodSeconds: 10"
        echo "          # 启动探针startupProbe：成功一次后，才交由存活探针livenessProbe接管"
        echo "          # 针对慢启动应用livenessProbe延迟启动时间不好确定问题的优化措施"
        fi

        if [ -n "$APP_STARTUP_PATH" ]; then
        echo "          startupProbe: # 最长容器启动15*10+30=180s内liveness和readiness探针不会执行"
        echo "            httpGet:"
        echo "              path: $APP_STARTUP_PATH"
        echo "              port: $APP_HTTP_PORT"
        echo "            # 最少连续多少次失败才视为失败"
        echo "            failureThreshold: 10"
        echo "            # 两次执行的间隔时间"
        echo "            periodSeconds: 10"
        echo "            # 启动延时时间"
        echo "            initialDelaySeconds: 60"
        fi

        echo "      volumes:"
        echo "        - name: conf"
        echo "          configMap:"
        echo "            name: $APP_NAME-conf"
        for VOLUME_KEY in ${!K8S_DEPLOYMENT_VOLUME_KEYS[@]}; do
        local PVC_NAME=`echo ${K8S_DEPLOYMENT_VOLUME_KEYS[$VOLUME_KEY]} | awk -F: '{print $2}'`
        if [ -n "$PVC_NAME" ]; then
        echo "        - name: $VOLUME_KEY"
        echo "          persistentVolumeClaim:"
        echo "            claimName: $PVC_NAME"
        fi
        done
        echo "        - name: localtime"
        echo "          hostPath:"
        echo "            type: File"
        echo "            path: /etc/localtime"

        if [ -n "$K8S_DEPLOYMENT_IMAGE_SECRET" ]; then
        echo "      imagePullSecrets:"
        echo "        - name: $K8S_DEPLOYMENT_IMAGE_SECRET"
        fi
    fi

    trace "K8S_SERVICE_ENABLE=$K8S_SERVICE_ENABLE"
    if [ "$K8S_SERVICE_ENABLE" == "true" ]; then
        echo "---"
        echo "apiVersion: v1"
        echo "kind: Service"
        echo "metadata:"
        echo "  name: $K8S_DEPLOYMENT_NAME-svc"
        echo "  namespace: $K8S_NAMESPACE"
        echo "  labels:"
        echo "    name: $K8S_DEPLOYMENT_NAME"
        echo "spec:"
        echo "  type: $K8S_SERVICE_TYPE"
        echo "  selector:"
        echo "    name: $K8S_DEPLOYMENT_NAME"
        echo "  ports:"
        if [ -n "$K8S_SERVICE_HTTP_PORT" ]; then
        echo "    - name: http"
        echo "      protocol: TCP"
        echo "      port: $K8S_SERVICE_HTTP_PORT"
        echo "      targetPort: $APP_HTTP_PORT"
        fi
        
        if [ -n "$K8S_SERVICE_NODE_PORT" ] && [ "$K8S_SERVICE_TYPE" == "NodePort" ]; then
        echo "      nodePort: $K8S_SERVICE_NODE_PORT"
        fi
        
        if [ -n "$K8S_SERVICE_GRPC_PORT" ]; then
        echo "    - name: grpc"
        echo "      protocol: TCP"
        echo "      port: $K8S_SERVICE_GRPC_PORT"
        echo "      targetPort: $APP_GRPC_PORT"
        fi

        trace "APP_PORTS=${APP_PORTS[@]}"
        trace "K8S_SERVICE_EXTRA_TYPE=${K8S_SERVICE_EXTRA_TYPE[@]}"
        trace "K8S_SERVICE_EXTRA_PORT=${K8S_SERVICE_EXTRA_PORT[@]}"
        trace "K8S_SERVICE_EXTRA_PROTOCOL=${K8S_SERVICE_EXTRA_PROTOCOL[@]}"
        trace "K8S_SERVICE_EXTRA_NODE_PORT=${K8S_SERVICE_EXTRA_NODE_PORT[@]}"

        if [ -n "$APP_PORTS" ]; then
        for APP_PORT in ${APP_PORTS[@]}; do
        echo "---"
        echo "apiVersion: v1"
        echo "kind: Service"
        echo "metadata:"
        echo "  name: $K8S_DEPLOYMENT_NAME-$APP_PORT-svc"
        echo "  namespace: $K8S_NAMESPACE"
        echo "  labels:"
        echo "    name: $K8S_DEPLOYMENT_NAME"
        echo "spec:"
        if [ -n "${K8S_SERVICE_EXTRA_TYPE[$APP_PORT]}" ]; then
        # 设置了类型则使用设置的类型
        echo "  type: ${K8S_SERVICE_EXTRA_TYPE[$APP_PORT]}"
        else
        # 默认为ClusterIP
        echo "  type: ClusterIP"
        fi
        echo "  selector:"
        echo "    name: $K8S_DEPLOYMENT_NAME"
        echo "  ports:"
        if [ -n "${K8S_SERVICE_EXTRA_TYPE[$APP_PORT]}" ]; then
        echo "    - name: port-${APP_PORT}"
        if [ -n "${K8S_SERVICE_EXTRA_PROTOCOL[$APP_PORT]}" ]; then
        # 设置了协议则使用设置的协议
        echo "      protocol: ${K8S_SERVICE_EXTRA_PROTOCOL[$APP_PORT]}"
        else
        # 默认使用TCP协议
        echo "      protocol: TCP"
        fi
        echo "      port: ${K8S_SERVICE_EXTRA_PORT[$APP_PORT]}"
        echo "      targetPort: $APP_PORT"
        if [ -n "${K8S_SERVICE_EXTRA_NODEPORT[$APP_PORT]}" ] && [ "${K8S_SERVICE_EXTRA_TYPE[$APP_PORT]}" == "NodePort" ]; then
        echo "      nodePort: ${K8S_SERVICE_EXTRA_NODEPORT[$APP_PORT]}"
        fi
        fi
        done
        fi
    fi
}

function encrypt_config_map(){
    if [ ! -d "$ENCRYPT_CONFIG_DIR" ]; then
        # 创建加密配置文件夹
        mkdir $ENCRYPT_CONFIG_DIR
    fi

    # 保存原始配置
    origin=./$ENCRYPT_CONFIG_DIR/$1.origin
    eval "echo \"`sed -e 's/\"/\\\\\\\"/g' $ConfigMap`\"" > $origin

    # 请求license接口加密配置
    res=`curl -X POST \
     -F "files=@$origin" \
     -F "sn=$SN"\
     $ENCRYPT_URL`
    
    # 解析json得到加密配置
    encryptConfig=`echo $res | jq .\"$1.origin\"`
    # 去除字符串双引号
    encryptConfig=${encryptConfig//\"/}

    # 保存到指定文件中
    echo $encryptConfig > ./$ENCRYPT_CONFIG_DIR/$1
}
