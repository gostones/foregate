# foregate

export PORT=8080

Server:

foregate server --bind ${PORT}


Client:

FG_URL="https://foregate.run.aws-usw02-pr.ice.predix.io/"
FG_URL="http://localhost:${PORT}/"

foregate client --url $FG_URL --hostport localhost:23001

##
<!-- 
./foregate client --url https://btcpay.run.aws-usw02-pr.ice.predix.io/ --port 5080 --hostport localhost:23001 --proxy $http_proxy
./foregate client --url https://wordpress.run.aws-usw02-pr.ice.predix.io/ --port 5080 --hostport localhost:80 --proxy $http_proxy
./foregate connect --url https://btcpay.run.aws-usw02-pr.ice.predix.io/ --ports 7500:7500 --proxy $http_proxy
-->
