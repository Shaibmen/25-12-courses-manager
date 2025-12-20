FROM node:20-bullseye  AS build

ARG VITE_API_URL_CORE
ARG VITE_API_URL_AUTH
ARG VITE_API_URL_DOC

ENV VITE_API_URL_CORE=$VITE_API_URL_CORE
ENV VITE_API_URL_AUTH=$VITE_API_URL_AUTH
ENV VITE_API_URL_DOC=$VITE_API_URL_DOC


WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build 

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80