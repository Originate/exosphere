FROM nodesource/node

# These steps ensure that npm install is only run when package.json changes
COPY ./package.json .
RUN curl -o- -L https://yarnpkg.com/install.sh | bash
ENV NODE_ENV=""
RUN $HOME/.yarn/bin/yarn install
COPY . .
