echo Starting Zookeeper service...
set ZOOKEEPER_HOME=F:\\zookeeper\\3.6.4
set ZOOCFG=%ZOOKEEPER_HOME%\\conf\\zoo.cfg
call "%ZOOKEEPER_HOME%\\bin\\zkServer.cmd" start