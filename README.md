<div align="center">
  <img alt="cgapp logo" src="https://i.imgur.com/QdUJPDU.png" width="100px"/>

# Syncron

Usage Patterns all in one cli based application written in Golang.

[![GoProject](https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go)](https://github.com/RedHatCRE/syncron) [![GoReport](https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none)](https://goreportcard.com/badge/github.com/redhatcre/syncron) ![License](https://img.shields.io/badge/license-apache_2.0-red?style=for-the-badge&logo=none)
</div>
<p align="center">
  <img src="https://i.imgur.com/AtTFOVi.png">
</p>

## 🔧 Installation


- [Download](https://golang.org/dl/) and install **Go**. 
    > Follow the link for instructions
    > 🔔 Please note: version 1.18 or higher required
- Clone this repository.
    
    > `git clone git@github.com:RedHatCRE/syncron.git`

- Navigate to local cloned repository folder and install Syncron with:
    
    > `go install`

- Check that Syncron has been properly installed on your environment by running
    > `syncron -v`
    > `syncron version 1.0.0`

    - Make sure your gopath is added to your path
      >  `export PATH=$PATH:$(go env GOPATH)/bin`



## 📖 Setup

To use Syncron, two important steps must be taken.

- Create **configuration file**. 
    - Important information about config:
        - Yaml based
        - Path: ~/.config/
        - Naming: syncron.yaml
        - A minimal syncron configuration file is as follows:

            ```yaml
            s3:
              endpoint: "<endpoint>"
              region: "<region>"
              bucket: "<name of bucket to pull from>"
            prefix: "<targeted keys starting path>"
            downloadDir: "<path where files will be downloaded>"
            ```

- For downloading files from a AWS buckets, proper credentials must be present on running machine. Two options:
  - Credentials file at $HOME/.aws/credentials
    - The credentials file has the following format:
    
      ```
      [default]
      aws_access_key_id = "XXXXXXX"
      aws_secret_access_key = "XXXXXXX"
      ```
      
  - Environment variables must be set:
    - export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
    - export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY


## ⚙️ Commands & Options


The user is given several flags to pick which time frame to pull data from.


> ##### **Download files**

| Option | Description                                              | Type   | Required? |
|--------|----------------------------------------------------------|--------|-----------|
| `download`   | Download files.| `string` | No        |

```bash
syncron download [option] [--flag] [number]
```

| --flag | Description                                              | Type   | Default | Required? |
|--------|----------------------------------------------------------|--------|---------|-----------|
| `days`   | Download files from the past x days. | `int` | `2` | No        |
| `months`   | Download files from the past x months. | `int` | `0` | No        |
| `years`   | Download files from the past x years. | `int` | `0` | No        |
| `filter`   | Filter files to download | `[]str` | []str{} | No        |

| Option | Description                                              | Type   | Default | Required? |
|--------|----------------------------------------------------------|--------|---------|-----------|
| `sosreports`   | Download sosreports files.| `string` | `sosreports` | Yes        |

> ##### **Read parquet files**

```bash
syncron read-parquet [--flag] [option]
```

| Option | Description                                              | Type   | Required? |
|--------|----------------------------------------------------------|--------|-----------|
| `read-parquet`   | Read local parquet files.| `string` | No        |

| --flag | Description                                              | Type   | Default | Required? |
|--------|----------------------------------------------------------|--------|---------|-----------|
| `file`   | What file to read. | `str` | `-` | Yes        |
| `output`   | Path to place unpacked file. | `str` |  | Yes        |

> ##### **Query data from database**

```bash
syncron queries
```

## ⚠️ License

`Syncron` is free and open-source software licensed under the [Apache 2.0 License](https://github.com/RedHatCRE/syncron/blob/main/LICENSE). 
