<div align="center">
  <img alt="cgapp logo" src="https://seeklogo.com/images/G/go-logo-046185B647-seeklogo.com.png" width="100px"/>

# Syncron

Easily fetch files from **S3 buckets** with a cli based application written in Golang.

[![GoProject](https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go)](https://github.com/RedHatCRE/syncron) [![GoReport](https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none)](https://goreportcard.com/badge/github.com/redhatcre/syncron) ![License](https://img.shields.io/badge/license-apache_2.0-red?style=for-the-badge&logo=none)
</div>

## ⚡️ Quick start

[Download](https://golang.org/dl/) and install **Go**. 

> 🔔 Please note: version 1.18 or higher required

Run:

```bash
go run cmd/adhoc/main.go
```

## 📖 Setup

To use Syncron, two important steps must be taken.

- Create **configuration file**. 
    - Important information about config:
        - Yaml based
        - Path: root of the project
        - Naming: syncron.yaml
        - A minimal syncron configuration file is as follows:


            ```yaml
            bucket: "<name of bucket to pull from>"
            s3:
              endpoint: "<endpoint>"
              region: "<region>"
            prefix: "<targeted keys starting path>"
            download_dir: "<path where files will be downloaded>"
            ```
- Proper credentials must be present on running machine at $HOME/.aws/credentials


## ⚙️ Commands & Options


The user is given several flags to pick which time frame to pull data from.

```bash
go run cmd/adhoc/main.go download [option] [--flag] [number]
```

| Option | Description                                              | Type   | Default | Required? |
|--------|----------------------------------------------------------|--------|---------|-----------|
| `sosreports`   | Download sosreports files.| `string` | `sosreports` | Yes        |

| --flag | Description                                              | Type   | Default | Required? |
|--------|----------------------------------------------------------|--------|---------|-----------|
| `days`   | Download files from the past x days. | `int` | `2` | No        |
| `months`   | Download files from the past x months. | `int` | `0` | No        |
| `years`   | Download files from the past x years. | `int` | `0` | No        |



## ⚠️ License

`Syncron` is free and open-source software licensed under the [Apache 2.0 License](https://github.com/RedHatCRE/syncron/blob/main/LICENSE). 