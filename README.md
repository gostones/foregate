# foregate

export PORT=8080

Server:

foregate server --bind ${PORT}


Client:

FG_URL="https://foregate.run.aws-usw02-pr.ice.predix.io/__/tunnel"
FG_URL="http://localhost:${PORT}/__/tunnel"

foregate client --url $FG_URL --hostport localhost:23001