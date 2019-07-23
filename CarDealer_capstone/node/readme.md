## NODE JS SDK Sample

Here is a sample application that utilizes the Hyperledger Fabric NODE JS SDK to 

* Connect to the Oracle Autonomous Blockchain Cloud Service (OABCS) network using a set of config files
* Connect to a channel
* Install chaincode written in the "go" programming language
* Instantiate chaincode on a set of specific peers on a specific channel
* Invoke chaincode

It demonstrates how you could utilize the **__fabric-client__** Node.js SDK APIs.

The "network.yaml" file located in the parent directory mirrors your existing Oracle Autonomous Blockchain Cloud Service environment. Namely it describes

* A client
* Channels
* An organization
* Orderers
* Peers 
* Rest Proxies

It also describes where the security certificates with which to connect with your environment are located.

### Step 1: Install prerequisites

* **Node.js** v 6x

### Step 2: Initialize the sample application

We need to use the "npm" package manager to initialize the application. 
To do this, for Linux and MacOS,  run the following command in your terminal from the current directory: `sh npm_bcs_client.sh [-g]`; for Windows, run this: `call npm_bcs_client_win.bat [-g]`.

Note:
The sample provides a script for 'grpc-node' module, and you can rebuild this module from the official source code yourself. Instructions about how to rebuild the 'grpc-node' module are available in the OABCS Console 'Developer Tools' pages.

### Step 3: Modify configuration files

In the current directory "car-demo.js", change the `CHANNEL_NAME` to the channel you wish to utilize to run the sample. The default channel is provided as 'default'.
Or, you can add a param for the command when you run the sample: `node car-demo.js [channelName]`. 

Notice, if you want to run the sample on a new channel which is not included in the `network.yaml`, you should download a new `network.yaml` config file from OABCS.  

### Step 4: Run the sample application

To run the application, execute the following node command: `node car-demo.js [channelName]`.

"All Done"
