# Oauth_

### Run API

Before you run the API need a .env file, check the .env-sample for the variables that you need and
create yours or change the name for use this one

* Run the API with the commands below
```shell
    make up
```
* This one do a rebuild docker images
```shell
    make up_build
```

* Run with out docker
```shell
    make up_dbs
```
###### if you use vscode check the lauch-sample.json to add in .vscode folder and run the API with out the API in docker last but not least, if you don't use vscode you need to build the project with the environment variables in the launch-sample.json
----
### Docs
* See the documentation ones the API is up in the following URL

[http://localhost:8080/api/OAuth/v1/public/docs/index.html](http://localhost:8080/api/OAuth/v1/public/docs/index.html)