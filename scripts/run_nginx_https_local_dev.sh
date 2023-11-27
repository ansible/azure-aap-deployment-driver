#!/bin/bash

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

if [ $# -eq 1 ]; then
	CERTIFICATES_LOCATION=$1
else
	CERTIFICATES_LOCATION=${SCRIPTPATH}/../certificates
fi

if [ ! -d ${CERTIFICATES_LOCATION} ]; then
  CERTIFICATES_LOCATION=${SCRIPTPATH}/$1
	if [ ! -d ${CERTIFICATES_LOCATION} ]; then
		CERTIFICATES_LOCATION=${SCRIPTPATH}/../certificates/$1
		if [ ! -d ${CERTIFICATES_LOCATION} ]; then
			echo "Can't find certificates folder ${CERTIFICATES_LOCATION}"
			exit 1
		fi
	fi
fi

CERTIFICATES_LOCATION="$( cd -- "${CERTIFICATES_LOCATION}" >/dev/null 2>&1 ; pwd -P )"
echo "Certificates location: ${CERTIFICATES_LOCATION}"

# check that expected certificate files are present
if [ ! -f ${CERTIFICATES_LOCATION}/chain.pem -o  ! -f ${CERTIFICATES_LOCATION}/fullchain.pem -o ! -f ${CERTIFICATES_LOCATION}/privkey.pem ]; then
	echo "One or more of the expected files are missing in certitifactes location: ${CERTIFICATES_LOCATION}"
	echo "Required:"
	echo "  chain.pem"
	echo "  fullchain.pem"
	echo "  privkey.pem"
fi

mkdir -p ${SCRIPTPATH}/../nginx/logs
chmod a+wx ${SCRIPTPATH}/../nginx/logs

if [ ! -f "nginx/dhparam.pem" ]; then
	echo "Have to generate dhparam.pem file..."
	openssl dhparam -out ${SCRIPTPATH}/../nginx/dhparam.pem 2048
fi

echo "Building the container image for nginx..."
docker build -f Dockerfile.nginx-only.local_dev -t deploymentdrivernginx:latest .

echo "Starting the container for nginx... You can exit it with: ctrl-c"
docker run -ti --name nginx --rm \
  -v ${CERTIFICATES_LOCATION}:/etc/letsencrypt/live/aapinstaller:Z \
  -v $(pwd)/ui/build:/var/www/aapinstaller/public:Z \
	-v $(pwd)/nginx/logs:/var/log/nginx:Z \
	-v $(pwd)/nginx/dhparam.pem:/etc/nginx/dhparam.pem \
	--network host \
	deploymentdrivernginx:latest nginx -g "daemon off;"
