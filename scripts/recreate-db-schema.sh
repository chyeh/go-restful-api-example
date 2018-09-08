#!/bin/bash

psql $1 -f drop-schema.sql
psql $1 -f init-schema.sql