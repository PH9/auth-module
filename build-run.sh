rm ./auth-module || true

#go get -u "github.com/lib/pq"
#go get -u "github.com/satori/go.uuid"
#go get -u "github.com/natefinch/lumberjack"
#go get -u "github.com/PH9/go-lib"

go build -v -o auth-module
APP_MODE=DEBUG APP_ENVIRONMENT=UAT ./auth-module