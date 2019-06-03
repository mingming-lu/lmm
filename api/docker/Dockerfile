FROM golang:1.12.5

RUN export CLOUD_SDK_REPO="cloud-sdk-stretch" && \
    echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    apt-get update -y && apt-get install -y \
    google-cloud-sdk \
    google-cloud-sdk-app-engine-go \
    google-cloud-sdk-app-engine-python \
    google-cloud-sdk-datastore-emulator
