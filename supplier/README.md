# increment-subscribe
increment-subscribe

MySQL 全量，增量订阅（便于以后系统消费数据！）
1. row 模式
2. full binlog row image
3. 启动
[root@oms-test lt-test]# pwd
/root/go/src/lt-test
[root@oms-test lt-test]# go build supplier/monitor.go
4.vendor原始做了调整不要使用glide重新拉取；
5. 增加 服务重启不需要全量同步；只需要增量 ；另外binlog 定时2个小时更新 ；保持万一 服务down掉需要全量的情况
6. 如果rabbitmq 挂掉在重启；这边有重连机制；
7. rabbitmq挂掉钉钉报警；


