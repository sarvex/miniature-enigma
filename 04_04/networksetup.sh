SUBCRIPTION=""
#org01
ORDERER_ORG_RESOURCE_GROUP="hlfOrderer"
ORDERER_ORG_NAME="twotreesorderer"
ORDERER_ADMIN_IDENTITY="admin.$ORDERER_ORG_NAME"
 
#org02
PEER_ORG_RESOURCE_GROUP="hlfPeer"
PEER_ORG_NAME="twotreespeer"
PEER_ADMIN_IDENTITY="admin.$PEER_ORG_NAME"
 
CHANNEL_NAME="hlfchannel"
STORAGE_RESOURCE_GROUP="storageRG"
STORAGE_ACCOUNT="fabstor"
STORAGE_LOCATION="westus"
STORAGE_FILE_SHARE="hlfshare"

echo $'\e[1;33m' setup storage$'\e[0m'
az account set --subscription $SUBCRIPTION
az group create -l $STORAGE_LOCATION -n $STORAGE_RESOURCE_GROUP
az storage account create -n $STORAGE_ACCOUNT -g  $STORAGE_RESOURCE_GROUP -l $STORAGE_LOCATION --sku Standard_LRS

echo $'\e[1;33m' create storage share$'\e[0m'
STORAGE_KEY=$(az storage account keys list --resource-group $STORAGE_RESOURCE_GROUP  --account-name $STORAGE_ACCOUNT --query "[0].value" | tr -d '"')
az storage share create  --account-name $STORAGE_ACCOUNT  --account-key $STORAGE_KEY  --name $STORAGE_FILE_SHARE

echo $'\e[1;33m' connection strings$'\e[0m' 
STORAGE_KEY=$(az storage account keys list --resource-group $STORAGE_RESOURCE_GROUP  --account-name $STORAGE_ACCOUNT --query "[0].value" | tr -d '"')
SAS_TOKEN=$(az storage account generate-sas --account-key $STORAGE_KEY --account-name $STORAGE_ACCOUNT --expiry `date -u -d "1 day" '+%Y-%m-%dT%H:%MZ'` --https-only --permissions lruwd --resource-types sco --services f | tr -d '"')
AZURE_FILE_CONNECTION_STRING=https://$STORAGE_ACCOUNT.file.core.windows.net/$STORAGE_FILE_SHARE?$SAS_TOKEN


echo $'\e[1;33m' import orderer profile and msp$'\e[0m'
./azhlf adminProfile import fromAzure -o $ORDERER_ORG_NAME -g $ORDERER_ORG_RESOURCE_GROUP -s $SUBCRIPTION
./azhlf connectionProfile import fromAzure -g $ORDERER_ORG_RESOURCE_GROUP -s $SUBCRIPTION -o $ORDERER_ORG_NAME   
./azhlf msp import fromAzure -g $ORDERER_ORG_RESOURCE_GROUP -s $SUBCRIPTION -o $ORDERER_ORG_NAME

echo $'\e[1;33m' import peer profile and msp$'\e[0m'
./azhlf adminProfile import fromAzure -g $PEER_ORG_RESOURCE_GROUP -s $SUBCRIPTION -o $PEER_ORG_NAME
./azhlf connectionProfile import fromAzure -g $PEER_ORG_RESOURCE_GROUP -s $SUBCRIPTION -o $PEER_ORG_NAME
./azhlf msp import fromAzure -g $PEER_ORG_RESOURCE_GROUP -s $SUBCRIPTION -o $PEER_ORG_NAME


echo $'\e[1;33m' create channel$'\e[0m'
./azhlf channel create -c $CHANNEL_NAME -u $ORDERER_ADMIN_IDENTITY -o $ORDERER_ORG_NAME


echo $'\e[1;33m' consortium management$'\e[0m'
./azhlf msp export toAzureStorage -f $AZURE_FILE_CONNECTION_STRING -o $PEER_ORG_NAME
./azhlf msp import fromAzureStorage -o $PEER_ORG_NAME -f $AZURE_FILE_CONNECTION_STRING

./azhlf channel join -c $CHANNEL_NAME -o $ORDERER_ORG_NAME -u $ORDERER_ADMIN_IDENTITY -p $PEER_ORG_NAME
./azhlf consortium join -o $ORDERER_ORG_NAME -u $ORDERER_ADMIN_IDENTITY -p $PEER_ORG_NAME

./azhlf connectionProfile  export toAzureStorage -o $ORDERER_ORG_NAME -f $AZURE_FILE_CONNECTION_STRING
./azhlf connectionProfile  import fromAzureStorage -o $ORDERER_ORG_NAME -f $AZURE_FILE_CONNECTION_STRING


echo $'\e[1;33m' setup peers$'\e[0m'
./azhlf channel joinPeerNodes -o $PEER_ORG_NAME -u $PEER_ADMIN_IDENTITY -c $CHANNEL_NAME --ordererOrg $ORDERER_ORG_NAME
./azhlf channel setAnchorPeers -c $CHANNEL_NAME -p "peer1" -o $PEER_ORG_NAME -u $PEER_ADMIN_IDENTITY