FROM python:3.7.2-alpine3.9

RUN apk add --no-cache jpeg-dev zlib-dev
RUN apk add --no-cache postgresql-dev
RUN apk add --no-cache libmemcached-dev zlib-dev
RUN apk add --no-cache --virtual .build-deps build-base linux-headers

RUN pip3 install --upgrade pip

RUN apk add --no-cache --virtual .build-deps build-base linux-headers \
    && pip3 install pip --upgrade \
    && pip3 install Image \
    && apk del .build-deps

CMD ["sh"]
