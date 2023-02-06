mv ./chaincode.go ./azhlfTool/samples/chaincode/src/twotreescc/

#twotreespeer
ORGNAME=twotreespeer
USER_IDENTITY="admin.$ORGNAME"  
CC_NAME=readings  
CC_VERSION=1
CC_LANG=golang  
CC_PATH=/home/azure/azhlfTool/samples/chaincode/src/twotreescc   
CHANNEL_NAME=hlfchannel

./azhlfTool/azhlf chaincode install -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -p $CC_PATH -l $CC_LANG -v $CC_VERSION  

./azhlfTool/azhlf chaincode instantiate -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -v $CC_VERSION -c $CHANNEL_NAME -f init

./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f createBatch -a "b01" 1 "farmerA"
./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f addReadings -a "b01" 0 "NONE" "AuditorA"
./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f addReadings -a "b01" 1 "CHEMICALS" "AuditorA"
./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f addReadings -a "b01" 2 "2" "AuditorA"
./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f determineGrade -a "b01" "ProducerA"
./azhlfTool/azhlf chaincode invoke -o $ORGNAME -u $USER_IDENTITY -n $CC_NAME -c $CHANNEL_NAME -f queryBatch -a "b01"


#debugging issues
AKS_CLUSTER_RESOURCE_GROUP="hlfPeerTest"
AKS_CLUSTER_NAME="pt1-HLF-AKS"
AKS_CLUSTER_SUBSCRIPTION="aea4b771-0ca3-477c-bded-ff01b0c45dee"
az aks get-credentials -n $AKS_CLUSTER_NAME -g $AKS_CLUSTER_RESOURCE_GROUP
kubectl get pods -n hlf
kubectl logs peer1-78898bff-dmsvb -n hlf
