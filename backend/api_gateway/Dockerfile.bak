FROM nginx:latest

COPY ./api_gateway.conf /etc/nginx/conf.d/api_gateway.conf
COPY ./domain.ext /etc/nginx/ssl/domain.ext
COPY ./domain.crt /etc/nginx/ssl/domain.crt
COPY ./certificate.key /etc/nginx/ssl/certificate.key

RUN rm /etc/nginx/conf.d/default.conf