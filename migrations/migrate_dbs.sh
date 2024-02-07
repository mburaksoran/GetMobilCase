docker-compose up -d
docker exec -i getmobilcase-mySql-1 mysql --user=root --password=root  get_mobil_case_db < ./migrations/sql_migrate_order_up.sql
docker exec -i getmobilcase-mySql-1 mysql --user=root --password=root  get_mobil_case_db < ./migrations/sql_migrate_up.sql
docker exec -i mongodb mongosh --authenticationDatabase admin < ./migrations/mongo_migrate_up.js
docker exec -i localstack awslocal sqs --endpoint-url=http://localhost:4566 create-queue --queue-name order_updates --region eu-west-1