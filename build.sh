go build -ldflags " \
-X 'main.Version=$(git describe --tags --always)' \
-X 'main.GitCommit=$(git rev-parse --short HEAD)' \
-X 'main.BuildTime=$(date -u '+%Y-%m-%d %H:%M:%S')' \
" -o aid