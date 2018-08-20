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
12. 远程mysql测试；【在另外一台主机；mysqldump and increment 测试；】


操作：

0. bind-address=0.0.0.0

1. CREATE USER root IDENTIFIED BY 'root';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'root'@'%';
-- GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' ;
FLUSH PRIVILEGES;CREATE USER root IDENTIFIED BY 'root';
                 GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'root'@'%';
                 -- GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' ;
                 FLUSH PRIVILEGES;
grant all privileges on *.* to root@'%' identified by '123456';
                 FLUSH PRIVILEGES;

2. 数据库db ;table要对应好！







简单架构图

  mysql (mask slave ) ->  mq (broker) ->  异构数据处理（目的地）




# todo 

1. 区分数据表；各个执行；（并发）
2. 将数据库直接执行到目的端；
3. DDL处理；
4.由于binlog 需要同步所有日志；会严重耗费网络带宽；（即使设置了表，它是基于所有过来的日志在过滤
可以考虑只分析指定的database->table的数据；严重减小网络带宽的消耗。）




