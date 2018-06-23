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
8. http://127.0.0.1:5000/log/ 调试日志服务
9.mysql挂掉后自动重连
10.rabbitmq 服务健康检查（报警），防止服务宕机后 rows更新过多；阻塞后内存大幅增大；
11. gtid mode支持；




简单架构图

  mysql (mask slave ) ->  mq (broker) ->  异构数据处理（目的地）

