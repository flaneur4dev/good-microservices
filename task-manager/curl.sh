curl -L -X POST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
--data-raw '{
    "id": "8b171",
    "name": "Task 1",
    "description": "A task 1 that need to be executed at the timestamp specified",
    "timestamp": 1645275972000
}'

# curl -L -X GET 'localhost:8080/task'

# curl -L -X GET 'localhost:8080/task/8b171'

# curl -L -X DELETE 'localhost:8080/task/8b171'
