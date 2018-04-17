# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
# From specifies which container image we want to build upon.
FROM nginx

# Replace default NGINX configuration file with our modified version.
ADD default.conf /etc/nginx/conf.d/default.conf

# <relative-source-file-path> <absolute-target-file-path>
COPY . /usr/share/nginx/html

EXPOSE 443 80