@echo off

(del shandianyu-minisdk-mailer || echo 1) && go mod tidy && go env -w GOOS=linux && go build --tags netgo -o shandianyu-minisdk-mailer main.go && (del shandianyu-minisdk-mailer || echo 1)