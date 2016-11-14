FROM node
ARG SERVICE_NAME

# These steps ensure that npm install is only run when package.json changes
RUN mkdir -p /usr/src/$SERVICE_NAME
WORKDIR /usr/src/$SERVICE_NAME

ADD . /usr/src/$SERVICE_NAME
