## Nginx

```
wget -O - http://nginx.org/keys/nginx_signing.key | apt-key add -
apt-add-repository "deb http://nginx.org/packages/ubuntu/ bionic nginx"
apt-get update
apt-get install nginx
```

## Opentracing

#### Nginx

Opentracing moduel
```
cd /usr/lib/nginx/modules
wget -O - https://github.com/opentracing-contrib/nginx-opentracing/releases/download/v0.7.0/linux-amd64-nginx-1.14.0-ngx_http_module.so.tgz | tar zxf -
```

Jaeger module
```
cd /usr/local/lib
wget -O libjaegertracing_plugin.so https://github.com/jaegertracing/jaeger-client-cpp/releases/download/v0.4.2/libjaegertracing_plugin.linux_amd64.so
```

```
cat > /etc/nginx/jaeger-config.json
{
  "service_name": "nginx",
  "diabled": false,
  "reporter": {
    "logSpans": true,
    "localAgentHostPort": "127.0.0.1:6831"
  },
  "sampler": {
    "type": "const",
    "param": "1"
  }
}
```

/etc/nginx/nginx.conf
```
http {
    # opentracing
    opentracing on;
    opentracing_load_tracer /usr/local/lib/libjaegertracing_plugin.so /etc/nginx/jaeger-config.json;
    opentracing_tag http_user_agent $http_user_agent;
```

/etc/nginx/site-enabled/*.confg
```
location ~ {
    opentracing_trace_locations off;
    opentracing_operation_name $host:$server_port;
    opentracing_propagate_context;
}
```

## certbot

```
sudo add-apt-repository ppa:certbot/certbot
sudo apt install python-certbot-nginx
```

```
certbot certonly --nginx --cert-name kihamo.ru -d kihamo.ru -d www.kihamo.ru
```

```
cat /etc/cron.d/certbot 
# /etc/cron.d/certbot: crontab entries for the certbot package
#
# Upstream recommends attempting renewal twice a day
#
# Eventually, this will be an opportunity to validate certificates
# haven't been revoked, etc.  Renewal will only occur if expiration
# is within 30 days.
SHELL=/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

0 */12 * * * root test -x /usr/bin/certbot -a \! -d /run/systemd/system && perl -e 'sleep int(rand(43200))' && certbot -q renew
```