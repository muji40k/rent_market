
FROM byjg/nginx-extras

RUN mkdir -p /etc/nginx/sites-enabled/

COPY nginx/main.conf /etc/nginx/nginx.conf
COPY nginx/balance_nginx.conf /etc/nginx/sites-enabled/web_backend.conf

RUN mkdir -p /static/
RUN mkdir -p /static/docs
RUN mkdir -p /render/md

COPY readme.md /static/docs/
COPY loadbalance.md /static/docs/
COPY res/ /static/docs/res/
COPY openapi.json /static/docs/
COPY assets/index.html /static/
COPY assets/md-renderer.html /render/md/

