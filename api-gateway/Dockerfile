FROM nginx:latest

EXPOSE 80

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY index.html /var/www/html/

COPY 400.html  401.html  500.html /var/www/html/errors/

CMD ["nginx", "-g", "daemon off;"]
