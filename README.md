# qingcloud-cni
QingCloud Iaas cni plugin and daemon

## How it works

QingCloud plugin runs a small agent on each host which is responsible for allocating and recollecting nics. And each container
would get a nic when it is started. In this way containers are connected to each other without adding additional network overhead.

## Details 

To work with scheduling framework such as Kubernetes and Mesos, two components should be configured. One is daemon process and the other is cni plugin.
Daemon process manages a table of available nic on host and nic plugin works with scheduling framework to get nic from daemon process
and add ip address and routing rules to nic. Daemon process also runs a few cron tasks to detach unused nic from host.


## Getting started

1. run qingcloud agent

    qing cloud agent is started by 'qingagent start' command. parameters are described as below.
    
    ```commandline
    [martin@MartinLaptop qingcloud-cni]$ ./bin/qingagent start -h
    QingCloud container networking agent is a daemon process which allocates and recollects nics resources.
    
    Usage:
      qingagent start [flags]
    
    Flags:
          --QYAccessFilePath string   QingCloud Access File Path (default "/etc/qingcloud/client.yaml")
          --bindaddr string           bind address of daemon process (default "0.0.0.0:31080")
          --gatewayns string          gateway nic name (default "hostnic")
      -h, --help                      help for start
          --iface string              Default nic which is used by host and will not be deleted (default "eth0")
          --policy string             policy of Selecting which vxnet to create nic from.(FailRotate,RoundRotate,Random) (default "FailRotate")
          --vxnet stringSlice         vxnet id list (default [vxnet-xxxxxxx])
          --zone string               QingCloud zone (default "pek3a")
    
    Global Flags:
          --config string     config file (default is ./QingAgent.yaml)
          --loglevel string   log level(debug,info,warn,error) (default "info")

    ```
    
    one can also write these configuration into config file which resides in the working directory and is named as QingAgent.
    The config file support multiple format.(Such as yaml, xml, json and so on.) For example:
    
    ```yaml
    QYAccessFilePath: '/etc/qingcloud/client.yaml'
    vxnet: 'vxnet-xxxxxxx'
    iface: 'eth0'
    loglevel: 'debug'
    ```

1. configure cni plugin

    one should put qingcni to where the cni plugin reside. (e.g. /opt/cni/bin). and add one config file for qingcloud plugin.
    e.g. /etc/cni/net.d/10-mynet.conf 
    ```json
    {
            "cniVersion": "0.3.1",
            "name": "mynet",
            "type": "qingcni",
            "isGateway": true,
            "args": {
                    "BindAddr":"127.0.0.1:31080"
            },
            "ipMasq": true,
            "ipam": {
               "routes":[{"dst":"kubernetes service cidr","gw":"hostip or 0.0.0.0"}]
            }
    }
    ```
    
3. Special notes for Kubernetes users
   
    Hostnic may not work as expected when it is used with Kubernetes framework due to the constrains in the design of kubernetes. However, we've provided a work around to help users setup kubernetes cluster.
    
    When a new service is defined in kubernetes cluster, it will get a cluster ip. And kube-proxy will maintain a port mapping tables on host machine to redirect service request to corresponding pod. And all of the network payload will be routed to host machine before it is sent to router and the service request will be handled correctly. In this way, kubernetes helps user achieve high availability of service. However, when the pod is attached to network directly(this is what hostnic did), Service ip is not recognied by router and service requests will not be processed.
    
    So we need to find a way to redirect service request to host machine through vpc. Here we implemented a feature to write routing rules defined in network configuration to newly created network interface. And if the host machine doesn't have a nic which is under pod's subnet, you can just set gateway to 0.0.0.0 and network plugin will allocate a new nic which will be used as a gateway, and replace 0.0.0.0 with gateway's ip address automatically.
