# rm ./auth-module || true

go get -u "github.com/lib/pq"
go get -u "github.com/satori/go.uuid"
go get -u "github.com/natefinch/lumberjack"
go get -u "github.com/PH9/go-lib"

CGO_ENABLED=0 GOOS=linux go build -o mrtr-auth ./... || exit 0

export APP_AUTH_PATH="wasith/auth-module:golang"

docker build --no-cache -t $APP_AUTH_PATH .