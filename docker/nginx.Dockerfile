FROM nginx:latest
LABEL maintainer="Hebert <hebert.t.santos@gmail.com>"
COPY /docker/config/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80 443
ENTRYPOINT ["nginx"]
CMD ["-g", "daemon off;"]