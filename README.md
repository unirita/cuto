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

**options**

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

**options**

|Option       |Description                 |
|-------------|----------------------------|
|-v           |Show version information    |
|-c FilePath  |Set file path of servant.ini|

### Show

Show command is a viewer for Jobnet execution result.

    show -c /path/to/master.ini [options]

**options**

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

Copyright. (C) 2015 UNIRITA Inc,
