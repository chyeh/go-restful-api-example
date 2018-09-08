#!/bin/bash

cd scripts
./set-up-test-environment.sh "integration_test"

HTTP_CODE=$(
curl -sL -w "%{http_code}\\n" \
     -X GET http://localhost/recipes \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] GET /recipes"
else
    echo "[ FAILED ] GET /recipes"
fi

HTTP_CODE=$( \
curl -sL -w "%{http_code}\\n" \
     -X POST http://localhost/recipes \
     -H "Content-Type: application/json" \
     -H "Authorization: aGVsbG9mcmVzaDpoZWxsb2ZyZXNo" \
     -d '{"name":"Fabulous Fried Chicken","prepare_time":30,"is_vegetarian":false}' \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] POST /recipes"
else
    echo "[ FAILED ] POST /recipes"
fi

HTTP_CODE=$(
curl -sL -w "%{http_code}\\n" \
     -X GET http://localhost/recipes/1 \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] GET /recipes/{id}"
else
    echo "[ FAILED ] GET /recipes/{id}"
fi

HTTP_CODE=$(
curl -sL -w "%{http_code}\\n" \
     -X PUT http://localhost/recipes/1 \
     -H "Content-Type: application/json" \
     -H "Authorization: aGVsbG9mcmVzaDpoZWxsb2ZyZXNo" \
     -d '{"name":"好吃炸雞","difficulty":3,"is_vegetarian":false}' \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] PUT /recipes/{id}"
else
    echo "[ FAILED ] PUT /recipes/{id}"
fi

HTTP_CODE=$(
curl -sL -w "%{http_code}\\n" \
     -X POST http://localhost/recipes/1/rating \
     -H "Content-Type: application/json" \
     -d '{"rating":5}' \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] POST /recipes/{id}/rating"
else
    echo "[ FAILED ] POST /recipes/{id}/rating"
fi

HTTP_CODE=$(
curl -sL -w "%{http_code}\\n" \
     -X DELETE http://localhost/recipes/1 \
     -H "Authorization: aGVsbG9mcmVzaDpoZWxsb2ZyZXNo" \
     -o /dev/null --connect-timeout 1
)
if [ $HTTP_CODE -eq 200 ];then
    echo "[ PASSED ] DELETE /recipes/{id}"
else
    echo "[ FAILED ] DELETE /recipes/{id}"
fi

./tear-down-test-environment.sh "integration_test"