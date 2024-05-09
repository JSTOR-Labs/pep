FROM node:11.13.0-alpine

# create destination directory
RUN mkdir -p /usr/src/nuxt-app
WORKDIR /usr/src/nuxt-app

# update and install dependency
RUN apk update && apk upgrade

ADD package.json package-lock.json /usr/src/nuxt-app/
RUN npm ci

# copy the app, note .dockerignore
COPY . /usr/src/nuxt-app/

# set app serving to permissive / assigned
ENV NUXT_HOST=0.0.0.0
# set app port
ENV NUXT_PORT=8080
ENV DEPLOY_ENV JSTOR
# expose 8080 on container
EXPOSE $NUXT_PORT

# build necessary, even if no static files are needed,
# since it builds the server as well
RUN npm run install
RUN npm run build

# start the app
CMD [ "npm", "start" ]
