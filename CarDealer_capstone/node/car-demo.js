var install = require('./app/install-chaincode.js');
var instantiate = require('./app/instantiate-chaincode.js');
var invoke = require('./app/invoke-transaction.js');
var query = require('./app/query.js');

//Abstract commandline arguments  
var args = process.argv.splice(2);
// Param declaration
var CHANNEL_NAME = 'default';
var CHAINCODE_ID = 'obcs-cardealer';
var CHAINCODE_PATH = 'github.com';
var CHAINCODE_VERSION = 'v0';

//Reload the new args : Channel name
param_check(args);

var installChaincodeRequest = {
    chaincodePath: CHAINCODE_PATH,
    chaincodeId: CHAINCODE_ID,
    chaincodeVersion: CHAINCODE_VERSION
};

var instantiateChaincodeRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID,
    chaincodeVersion: CHAINCODE_VERSION,
    fcn: 'init',
    args: ["a"]
};

var targetPart = "ser".concat(getRandomInt(0,100000).toString());
var targetCar = "mer".concat(getRandomInt(0,1000000).toString());

var createVehiclePartRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID,
    fcn: 'initVehiclePart',
    args: [targetPart, 'tasa', '1502688979', 'airbag 2020', 'aaimler ag / mercedes', 'false', '0'],
};

var createCarRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID, 
    fcn: 'initVehicle',
    args: [targetCar, 'mercedes', 'c class', '1502688979', targetPart, 'mercedes', 'false', '0'],
};
         


// STEP 1
// Install chaincode
try {
    install.installChaincode(installChaincodeRequest.chaincodeId, installChaincodeRequest.chaincodePath,installChaincodeRequest.chaincodeVersion).then((result) => {
                console.log(result)
                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 1/4 : Successfully installed chaincode' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);
                // STEP 2
                // Instantiate chaincode
                return instantiate.instantiateChaincode(instantiateChaincodeRequest.chanName, instantiateChaincodeRequest.chaincodeId, 
                    instantiateChaincodeRequest.chaincodeVersion,instantiateChaincodeRequest.fcn,instantiateChaincodeRequest.args);
        }).then((result2) => {
                console.log(result2)
                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 2/4 : Successfully instantiated chaincode on the channel' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);
                // STEP 3
                // invoke chaincode to create a vehicle part
                return invoke.invokeChaincode(createVehiclePartRequest.chanName, createVehiclePartRequest.chaincodeId,
                    createVehiclePartRequest.fcn, createVehiclePartRequest.args);
        }).then((result3) => {
                console.log(result3)

                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 3/4 : Successfully committed vehicle part to ledger' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);

                // STEP 4
                // invoke chaincode to create a vehicle
                return invoke.invokeChaincode(createCarRequest.chanName, createCarRequest.chaincodeId,
                    createCarRequest.fcn, createCarRequest.args);
        }).then((result4) => {
                console.log(result4)

                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 4/4 : Successfully committed vehicle to ledger' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                console.log("All Steps Completed Sucessfully");
                process.exit();
        });
} catch (e) {
    console.log(
        '\n\n*******************************************************************************' +
        '\n*******************************************************************************' +
        '\n*                                          ' +
        '\n* Error!!!!!' +
        '\n*                                          ' +
        '\n*******************************************************************************' +
        '\n*******************************************************************************\n');
    console.log(e);
    return;
}


function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function param_check(args) {
    if (args.length > 0) {
        if (args[0] != undefined && args[0] != null) {
            CHANNEL_NAME = args[0];
        }
    }
}

function getRandomInt(min, max) {
	min = Math.ceil(min);
	max = Math.floor(max);
	return Math.floor(Math.random() * (max - min)) + min; //The maximum is exclusive and the minimum is inclusive
}
