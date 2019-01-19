## File Indexer [![Build Status](https://travis-ci.org/dgplug/dexer.svg?branch=master)](https://travis-ci.org/dgplug/dexer)

This is a simple project which reads files in the directory mentioned in the config file and indexes it. This then can be used to search the files present in the directory.



#### Motivation

This is a way of trying to index the logs file in [dgplug](https://dgplug.org/) so that we can easily find the files in which a word has occured in the file.

#### Dependency

This assumes that you have `Go` installed and setup.

`pip3` and `python3` is what we need for development purpose.



### How to install

> **Note** : Right now we are unable to provide binary builds, so have to build the program the regular way or you can build a [docker](#run-locally-using-docker) container.

1. Install `git`.
2. [Install Go](https://golang.org/doc/install) (need a version which supports Go Modules).
3. Run the following commands :

```bash
$ git clone https://github.com/dgplug/dexer.git
$ cd dexer
$ make all
$ dexer
```

To use the program we need a configuration file which is provided by the `config.json` file. Here is an example of one:
```json
{
    "RootDirectory": "logs",
    "IndexFilename": "irclogs.bleve",
    "Port": ":8000",
    "LogFile": "logfile"
}
```

There are 4 entries :

- **RootDirectory** : Location of the logs which we want to search through.
- **IndexFilename** : The file where `bleve` will store all the information related to indexing.
- **Port** : The port on which the server will be run.
- **LogFile** : The file which will be used to store the logs.

### API

Once everything is working fine install the postman plugin for your browser. And from that plugin you need to hit the endpoint as:

`localhost:<port>/search/american`

Here, `port ` number passed to the program using the configuration file ,`american` is the query word that I passed, make sure to open any file in the logs/ directory and find a word to search. It will look like this:

![Missing screenshot](https://image.ibb.co/cWP5TA/postman-query.png "Postman Screenshot")

### Use the Web Front End

You can also visit `localhost:<port>` to use the web frontend which comes with the program to search.

![Web Front End](https://i.ibb.co/x2FDb0F/Screenshot-20181227-114313.png)

### Run locally using docker

You can run the application using Docker in your local machine. It will use the `Dockerfile` instructions. Make sure you have [Docker](https://www.docker.com/) installed in your machine.

Run the following commands to build and run the docker image.

```bash
$ git clone https://github.com/dgplug/dexer.git
$ cd dexer
$ make logs
$ make docker-build # builds docker image
$ make docker-run # runs the image in new container
```

One has to make sure the `logs` directory has all the file because dexer runs the indexing at the starting and then keeps it, so if a file is not there it would not index it.
