FROM hashicorp/terraform:0.8.4

RUN apk add --update nodejs bash
RUN npm install -g n
RUN n 6.3.1

# These steps ensure that npm install is only run when package.json changes
RUN mkdir -p /usr/src/exo-deploy
COPY ./package.json /usr/src/exo-deploy
WORKDIR /usr/src/exo-deploy
RUN npm install --production

# Copy source code
COPY . /usr/src/exo-deploy

# Create directory for generating terraform scripts
RUN mkdir -p /usr/src/terraform

ENTRYPOINT ["./bin/start-deploy"]
