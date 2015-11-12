# GoCuto

GoCuto is a process controller tool that can execute cross-server processes in the order you like.


## Installation

You can download binary packages on [GoCuto Official Website](http://cuto.unirita.co.jp/cuto/).

Or, you can build GoCoto with go command. (need Go 1.4 or higher)

    go get github.com/unirita/cuto/master
    go get github.com/unirita/cuto/servant
    go get github.com/unirita/cuto/show


## Commands

### Master

Master command executes group of processes (called Jobnet).
Use this command when you want to run Jobnet.

    master -n JobnetName -s -c /path/to/master.ini

**Options**

|Option       |Description                                                                                    |
|-------------|-----------------------------------------------------------------------------------------------|
|-v           |Show version information                                                                       |
|-n JobnetName|Set name of Jobnet                                                                             |
|-s           |Use this option if you want to run Jobnet. If didn't, master command only checks Jobnet syntax.|
|-c FilePath  |Set file path of master.ini                                                                    |

### Servant

Servant command is a resident process which executes processes by request from the Master.
You must run this command every server to use Master command.

    servant -c /path/to/servant.ini

**Options**

|Option       |Description                 |
|-------------|----------------------------|
|-v           |Show version information    |
|-c FilePath  |Set file path of servant.ini|

### Show

Show command is a viewer for Jobnet execution result.

    show -c /path/to/master.ini [options]

**Options**

|Option              |Description                                                                  |
|--------------------|-----------------------------------------------------------------------------|
|-v                  |Show version information                                                     |
|-help               |Show usage                                                                   |
|-c FilePath         |Set file path of master.ini                                                  |
|-jobnet JobnetName  |Narrow result by Jobnet                                                      |
|-nid InstanceID     |Narrow result by Instance ID (unique ID for every execution)                 |
|-from Date, -to Date|Narrow result by range of executed date                                      |
|-status Status      |Narrow result by status (select from "normal", "abnormal", "warn", "running")|
|-format Format      |Select output format from "json" or "csv"                                    |
|-utc                |Set or show date value as UTC timezone, not as local timezone                |


## Configuration

GoCuto uses some configuration files written by [toml format](https://github.com/toml-lang/toml).

### master.ini

master.ini is configuration file for Master command.

**Tables and Keys**

|Table|Key                   |Type   |Description                                                                          |
|-----|----------------------|-------|-------------------------------------------------------------------------------------|
|job  |default_node          |String |Host name of default node which Job (process) is executed on.                        |
|job  |default_port          |Integer|Port number of default node which Job is executed on.                                |
|job  |default_timeout_min   |Integer|Default time limit to wait end of Job execution. (minute)                            |
|job  |connection_timeout_sec|Integer|Time limit to wait connection keep alive signal. (second)                            |
|job  |time_tracking_span_min|Integer|Time span to display elapsed time from execution started time. (minute)              |
|job  |attempt_limit         |Integer|Max retry number of times when Job is not able to start.                             |
|dir  |jobnet_dir            |String |Directory to put Jobnet definition files in.                                         |
|dir  |log_dir               |String |Directory to output Master command log files.                                        |
|dir  |db_dir                |String |Directory to put execution result db file in.                                        |
|log  |output_level          |String |Minimum log level. Select from "trace", "debug", "info", "warn", "error", "critical".|
|log  |max_size_kb           |Integer|Max size of log file. (KByte)                                                        |
|log  |max_generation        |Integer|Max generation for log file rotation.                                                |
|log  |timeout_sec           |Integer|Time limit to wait log output ends.                                                  |

### servant.ini

master.ini is configuration file for Servant command.

**Tables and Keys**

|Table|Key               |Type   |Description                                                                          |
|-----|------------------|-------|-------------------------------------------------------------------------------------|
|sys  |bind_address      |String |Listen host name of servant.                                                         |
|sys  |bind_port         |Integer|Listen port number of servant.                                                       |
|job  |multi_proc        |Integer|Max number of Job execution at same time.                                            |
|job  |heartbeat_span_sec|Integer|Time span to send keep alive signel for master. (second)                             |
|dir  |job_dir           |String |Directory to put files be executed as Job in.                                        |
|dir  |joblog_dir        |String |Directory to output Job log files.                                                   |
|dir  |log_dir           |String |Directory to output Servant command log files.                                       |
|log  |output_level      |String |Minimum log level. Select from "trace", "debug", "info", "warn", "error", "critical".|
|log  |max_size_kb       |Integer|Max size of log file. (KByte)                                                        |
|log  |max_generation    |Integer|Max generation for log file rotation.                                                |
|log  |timeout_sec       |Integer|Time limit to wait log output ends.                                                  |

Copyright. (C) 2015 UNIRITA Inc,
