一套 Linux 平台下的自动化运维工具 agent 端

# 已实现
- apiserver
- cmdb
- 远程命令
- 文件分发

# 未实现
- 基础监控
- ipmi

```bash
# 初始化节点
$ ./apiserver --work-dir=/opt/hi-devops-agent/ init --public-ip=192.168.3.10 --server=http://192.168.3.100:8888

# 以守护进程方式运行 agent
$ ./apiserver start -d

# 停止以守护进程方式运行的 agent
$ ./apiserver stop

# 查看命令帮助
$ ./apiserver --help

# 查看子命令的命令帮助
$ ./apiserver init --help

$ ./apiserver start --help

$ ./apiserver stop --help 
```

