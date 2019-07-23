#!/bin/bash

##Default setting
MANU_NAME="mazida"
DEALER_A_NAME="wheel"
DEALER_B_NAME="light"
DEALER_A_CHANNEL="default"
DEALER_B_CHANNEL="default"
CC_NAME="obcs-cardealer"
CC_VERSION="v0"
REST_ADDR="https://0197755BED494AA496F12E715FB64150.blockchain.ocp.oraclecloud.com:443/restproxy1"
USER=""
PASSWORD=""


## Get dealer network config from user
# echo "We need the following details to initialize the ledger."
# echo  "Enter name of Vehicle Manufacturer and press [ENTER]: "
# read -e MANU_NAME
# echo -n "Enter name of Car Dealer A and press [ENTER]: "
# read -e DEALER_A_NAME
# echo -n "Enter name of Car Dealer B and press [ENTER]: "
# read -e DEALER_B_NAME
# echo -n "Enter A channel name and press [ENTER]: "
# read -e DEALER_A_CHANNEL
# echo -n "Enter B channel name and press [ENTER]: "
# read -e DEALER_B_CHANNEL
# echo -n "Enter Chaincode Name and press [ENTER]: "
# read -e CC_NAME
# echo -n "Enter Chaincode version and press [ENTER]: "
# read -e CC_VERSION
# echo -n "Enter address of your BCS REST proxy and press [ENTER]: "
# read -e REST_ADDR
# echo -n "Enter port of your BCS REST proxy and press [ENTER]: "
# read -e REST_PORT
if [[ "$USER" == "" ]]; then
	#statements
	echo -n "Enter username of your BCS REST proxy and press [ENTER]: "
	read -e USER
fi
if [[ "$PASSWORD" == "" ]]; then
	#statements
	echo -n "Enter password of your BCS REST proxy and press [ENTER]: "
    read -e PASSWORD
fi

echo

starttime=$(date +%s)

#check the ressponse
checkRep(){
    res=${@:1}
    if [[  $res = *Failure* ]]; then
    	if [[ $res = *exists* ]]; then
    		echo "Things already exists. Continue...."
    	else
    		echo "Option Failure, please check the response message and try again."
    		exit
    	fi
    fi

    if [[ $res = *Success* ]]; then
    	echo "Option Success, Continue..."
    fi
}

echo "Set Initial State of the Ledger"

echo "Create 5 parts for Sam Dealership"
echo
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["abg1234", "panama-parts", "1502688979", "airbag 2020", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE
echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["abg1235", "panama-parts", "1502688979", "airbag 4050", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["ser1236", "panama-parts", "1502688979", "seatbelt 10020", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["win1237", "panama-parts", "1502688979", "windshield auto201", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["bra1238", "bobs-bits", "1502688979", "brakepad 4200", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

echo "Transfer parts to Sam from DAuto"
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferVehiclePart","args":["bra1238", "'"$MANU_NAME"'", "'"$DEALER_A_NAME"'"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferVehiclePart","args":["win1237", "'"$MANU_NAME"'", "'"$DEALER_A_NAME"'"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferVehiclePart","args":["ser1236", "'"$MANU_NAME"'", "'"$DEALER_A_NAME"'"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

echo "Create 2 vehicles for Sam owned by DAuto"
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehicle","args":["dtrt10001", "'"$MANU_NAME"'", "a coupe", "1502688979", "abg1235", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehicle","args":["dtrt10002", "'"$MANU_NAME"'", "big pickup", "1502688979", "abg1234", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

echo "Transfer the airbags to the vehicles"
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferPartToVehicle","args":["abg1235", "dtrt10001"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferPartToVehicle","args":["abg1234","dtrt10002"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

echo "Transfer the vehicles to Sam"
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferVehicle","args":["dtrt10001", "'"$MANU_NAME"'", "'"$DEALER_A_NAME"'"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_A_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"transferVehicle","args":["dtrt10002", "'"$MANU_NAME"'", "'"$DEALER_A_NAME"'"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

echo "Create 5 parts for Jude Dealership"
RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_B_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["abg1239", "bobs-bits", "1502688979", "airbag w020", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_B_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["abg1240", "bobs-bits", "1502688979", "airbag w030", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_B_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["whl1241", "wizzard-auto", "1502688979", "wheel 28374", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_B_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["win1242", "panama-parts", "1502688979", "windshield auto201", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo

RESPONSE=$(curl -H "Content-type:application/json" -X POST -u $USER:$PASSWORD \
-d '{"channel":"'"$DEALER_B_CHANNEL"'","chaincode":"'"$CC_NAME"'","method":"initVehiclePart","args":["sen1243", "wizzard-auto", "1502688979", "sensor p228", "'"$MANU_NAME"'", "false", "1502688979"],"chaincodeVer":"'"$CC_VERSION"'"}' \
$REST_ADDR/bcsgw/rest/v1/transaction/invocation);
echo $RESPONSE
checkRep $RESPONSE

echo
echo



echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
