FROM nginx:latest

EXPOSE 80

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY index.html /var/www/html/

CMD ["nginx", "-g", "daemon off;"]
