# Predix Insights Command Line Interface

[![Build status](https://travis-ci.org/scottmcclary1/predix-insights-cli.svg?master)](https://travis-ci.org/scottmcclary1)

## Remove Configuration (optional)
```
$ rm -rf ~/.pi
```

## Build PI CLI
```
$ ./build.sh
```

## Source Bash Profile
```
$ source ~/.bash_profile
```

## Show PI CLI
```
$ pi -h
```

## Login & Configure
```
$ pi configure --interactive
```

## Create Flow Template
```
$ pi flow-template create --desc "PI CLI Lunch and Learn" --flowTemplateName "pi-cli-lunch-learn" --flowType "SPARK_JAVA" --templateFileName "spark-examples.zip" --templateFilePath "/Users/scottmcclary/Desktop/PredixInsightsExamples/spark-examples/spark-examples.zip" --flowTemplateVersion 1.0.0 -i
```
## Update Flow Template Spark Arguments
```
$ pi flow-template update-spark-args -i --sparkArgs "{\"sparkArguments\": {\"applicationArgs\":[\"100\"],\"className\":\"org.apache.spark.examples.SparkPi\"}}"
```

## List Flow Template
```
$ pi flow-template list
```

## Save Flow Template Tags
```
$ pi flow-template save-tags --tags "[\"type:dev\", \"size:large\"]" -i
```

## List Flow Template Tags
```
$ pi flow-template list-tags
```

## Post Flow
```
$ pi flow create --flowName "pi-cli-lunch-learn-flow" -i
```

## Update Flow Spark Arguments
```
$ pi flow update-spark-args --sparkArgs "{\"sparkArguments\": {\"applicationArgs\":[\"250\"],\"className\":\"org.apache.spark.examples.SparkPi\"}}" -i
```

## Save Flow Tags
```
$ pi flow save-tags --tags "[\"type:prod\", \"size:small\"]" -i
```

## List Flow Tags
```
$ pi flow list-tags
```

## Get Flow
```
$ pi flow list
```

## Add Flow Configuration File(s)
```
$ pi flow add-config-file --configFileDetails "[{\"FileName\": \"scott.json\", \"FileLocation\": \"/Users/scottmcclary/Desktop/scott.json\"}]"
```

## List Flow Configuration File(s)
```
$ pi flow list-config-files -i
```

## Launch Flow
```
$ pi flow launch -i
```

## Get ContainerID
```
$ pi instance list-containers -i
```

## Get Container Logs

### stdout
```
$ pi instance list-container-logs --containerLogSink 1 --tail -i
```

### stderr
```
$ pi instance list-container-logs --containerLogSink 0 --tail -i
```

## Delete Flow
```
$ pi flow delete -i
```

## Delete Flow Template
```
$ pi flow-template delete -i
```
