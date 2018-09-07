#!/bin/bash

psql postgres://hellofresh:hellofresh@localhost:5432/hellofresh -f drop-schema.sql
psql postgres://hellofresh:hellofresh@localhost:5432/hellofresh -f init-schema.sql