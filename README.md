# find-facility

## console

Download phantomjs

$ go get -u github.com/sclevine/agouti
$ go get -u github.com/comail/colog

$ go build

## aws

### IFTTT

https://ifttt.com/maker_webhooks
find_facility
https://maker.ifttt.com/use/b-bocpdYO2gSjRl1iqdhRJ
070-4221-3401

### aws

$ go get -u github.com/aws/aws-lambda-go/lambda

$ GOOS=linux GOARCH=amd64 go build -tags aws -o find
$ zip aws.zip find phantomjs

- ハンドラをfindに変更
- ランタイムをGo 1.xに変更

基本設定
- メモリを384MBに変更
- タイムアウトを15分に変更（最大値）

環境変数
- PATHに/var/taskを追加

- aws.zipをアップロード
