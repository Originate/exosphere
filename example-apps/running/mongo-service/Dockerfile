FROM node
ARG SERVICE_NAME

# These steps ensure that npm install is only run when package.json changes
RUN mkdir -p /usr/src/$SERVICE_NAME
ADD ./package.json /usr/src/$SERVICE_NAME/package.json
WORKDIR /usr/src/$SERVICE_NAME
RUN npm install --production

ADD . /usr/src/$SERVICE_NAME
