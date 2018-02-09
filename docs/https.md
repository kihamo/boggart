Первоначальная установка
========================
1. Ставим софт
```
sudo apt-get install certbot python-certbot-nginx
```
2. Регистрируемся
```
certbot register --email me@example.com
```
3. Создаем общий конфиг для чекера
```
# cat /etc/nginx/acme 
location /.well-known {
    root /var/www/html;
}
```
4. Конфиг для тулзы
```
# cat /etc/letsencrypt/cli.ini 
authenticator = webroot
webroot-path = /var/www/html
post-hook = service nginx reload
text = True
```

Подключение домена
==================
1. Подключаем чекер в конфиге nginx /etc/nginx/sites-enabled/* до location
```
include acme;
```
2. Запускаем мастер
```
certbot --authenticator webroot --installer nginx
```