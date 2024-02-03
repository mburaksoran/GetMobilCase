start-services:
	chmod u+r+x ./migrations/migrate.sh
	./migrations/migrate_dbs.sh

create-sqs:
	$(info #prepare sqs_client...)
	aws --endpoint http://localhost:4566 sqs create-queue --queue-name order_updates --region eu-west-1


