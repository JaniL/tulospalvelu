FROM docker.io/library/php:7.4-apache

ENV COMPOSER_ALLOW_SUPERUSER=1

RUN apt-get update \
    && apt-get install -y --no-install-recommends git unzip libpq-dev \
    && docker-php-ext-install pdo_pgsql \
    && a2enmod rewrite \
    && rm -rf /var/lib/apt/lists/*

COPY --from=docker.io/library/composer:2 /usr/bin/composer /usr/bin/composer

COPY docker/apache/tulospalvelu.conf /etc/apache2/conf-available/tulospalvelu.conf
RUN a2enconf tulospalvelu

WORKDIR /var/www/html
COPY . /var/www/html

RUN composer install --no-dev --optimize-autoloader --no-interaction
