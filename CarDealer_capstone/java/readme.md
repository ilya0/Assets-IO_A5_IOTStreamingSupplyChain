##JAVA SDK SAMPLE

Here is a sample application that utilizes the Hyperledger Fabric JAVA SDK to 

* Connect to the Oracle Autonomous Autonomous Blockchain Cloud Service (OABCS) network using a set of config files
* Connect to a channel
* Install chaincode written in the "go" programming language
* Instantiate chaincode on a set of specific peers on a specific channel
* Invoke chaincode

The "network.yaml" file located in the parent directory mirrors your existing Oracle Blockchain Cloud Service environment. Namely it describes

* A client
* Channels
* An organization
* Orderers
* Peers 
* Rest Proxies

It also describes where the security certificates with which to connect with your environment are located.

###Folder Structure

```
--java
  --fabric-sdk-java                       the src of fabric java sdk, including some testes and examples
  --resource                              the app config and endorsement policy config file, etc
    --chaincode-endorsement-policy.yaml   the endorsement policy for chaincode instantiate and invoke
    --demo.properties                     the chaincode and channel config
  --src                                   the src java code for demo 
  --libs                                  the third-party dependencies for project
  --pom.xml                               the maven project config file
```

###Fabric Java SDK 
The sdk resource url : [Fabric-Java-SDK](https://github.com/hyperledger/fabric-sdk-java)

###Demo Project 
#####Step 1:Install prerequisites
*	JAVA JDK
*	Maven

Please download and install [MAVEN](https://maven.apache.org/download.cgi).

#####Step 2:Initialize the sample application

You need to use the MAVEN build automation tool to initialize the application. To do this run the following command in your terminal in current directory: `sh install_and_run.sh install`

Or, you can do it by 2 step manually:
1. Install dependencies for the project with the command : `mvn install`.
2. Replace the grpc-netty package for the project with this command:`mvn install:install-file -Dfile=./libs/grpc-netty-1.11.0.jar -DgroupId=io.grpc -DartifactId=grpc-netty -Dversion=1.11.0 -Dpackaging=jar`.

Note:
The sample provides a patched java package for 'grpc-netty@1.11.0', and you can rebuild this package from the official source code yourself. Instructions about how to rebuild the 'grpc-netty' package are available in the OABCS Console 'Developer Tools' pages.


#####Step 3: Modify configuration files
 
Navigate to the directory "java/resource demo.properties". Change the `CHANNEL_NAME` to the target channel you want to utilize for the sample. The default config is `default`.

Notice, if you want to run the sample on a new channel which is not included in the `network.yaml`, you should download a new `network.yaml` config file from OABCS.  

#####Step 4: Run the sample application

To run the application, execute the following mvn command: 

`mvn exec:java -Dexec.mainClass="main.CarDemo"` 

Or, execute the command using this script:

`sh install_and_run.sh run`

"All Done"
