#!/bin/bash
if [ $# -eq 0 ]; then
        echo "Set image file version latest"
        version="latest"
else
    version="$1"
fi
echo "Set image file version $version"
docker build --platform=linux/amd64 -t registry:$version .
docker push registry:$version