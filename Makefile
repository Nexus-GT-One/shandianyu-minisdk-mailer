ACCOUNT=htfgtred908
NAME=shandianyu-minisdk-mailer

shandianyu-minisdk-mailer-beta:
	go mod tidy && go build -ldflags="-s -w" --tags netgo -o shandianyu-minisdk-mailer main.go
	docker build -t ${ACCOUNT}/${NAME}:beta --progress plain .

shandianyu-minisdk-mailer-prod:
	go mod tidy && go build -ldflags="-s -w" --tags netgo -o shandianyu-minisdk-mailer main.go
	docker login -u ${ACCOUNT} -p 'yU(fa|2=z6{4qn+?'
	docker build -t ${ACCOUNT}/${NAME} --progress plain .
	docker push ${ACCOUNT}/${NAME}