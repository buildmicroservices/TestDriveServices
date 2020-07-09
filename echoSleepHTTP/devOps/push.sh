# FOR NOW... use localhost line below 
#localhost:32000/echosleephttp:latest

#docker tag echosleephttp andros.karux.net:5000/echosleephttp
#docker push andros.karux.net:5000/echosleephttp
#docker tag echosleephttp nassau.karux.net:32000/echosleephttp
#docker push nassau.karux.net:32000/echosleephttp
if [ -z "$1" ]
  then
    echo -e "No registry hostname.\n"
    echo -e "please provide registry hostname arg1 (e.g. andros.karux.net)\n"
    exit 1
fi
if [ -z "$2" ]
  then
    echo -e "No port.\n"
    echo -e "please provide registry port arg1 (e.g. 5000)\n"
    exit 1
fi
REGISTRY=$1
REGISTRYPORT=$2
SOURCEIMAGE="echosleephttp"
docker tag ${SOURCEIMAGE} ${REGISTRY}:${REGISTRYPORT}/${SOURCEIMAGE}
docker push ${REGISTRY}:${REGISTRYPORT}/${SOURCEIMAGE}
