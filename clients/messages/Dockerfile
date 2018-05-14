FROM nginx

# Replace default NGINX configuration file with our modified version.
ADD default.conf /etc/nginx/conf.d/default.conf

# <relative-source-file-path> <absolute-target-file-path>
COPY src /usr/share/nginx/html

EXPOSE 443 80