#!/bin/bash

psql postgres://hellofresh:hellofresh@localhost:5432/hellofresh -f init-user-data.sql