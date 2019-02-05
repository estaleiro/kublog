# KubLog 
Tool to monitor IP changes in containers orchestrated with Kubernetes

## Usage
```
./kublog -config=kublog.toml
```

## Configuration
**Kublog**

|Parameter|Example|Description|
| ----------- | ----------- |-----------|
|kubeconfig|"/root/.kube/config"|kubeconfig file location |
|log_filename|"/var/log/kublog.log"|log location|
|period|1|how frequently data is collected (seconds)|

 **ElasticSearch Output**

|Parameter|Example|Description|
| ----------- | ----------- |-----------|
|hosts|["http://127.0.0.1:9200"]|ElasticSearch hosts|
|indexname|"kublogindex"|Index name where information will be indexed|
|timeout|"2s"|Timeout to connect to ElasticSearch (seconds)|
|enable_sniffer|false|true/false to enable sniffer|
|health_check_interval|"0s"|Interval for health check (seconds)|


## Contributing
Go version: 1.11.5

1. Get from github

```
mkdir myprojetct/src

cd myprojetct/src

git init

Fork the project to your user

git remote add origin https://github.com/mygithubuser/kublog.git

git config --global user.email "email@email.com"

git config --global user.name "mygithubuser"

git pull origin master
```
