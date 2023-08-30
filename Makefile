.PHONY: all build clean run check cover lint docker help

dateTime=`date +%F_%T`
ARCH="linux-amd64"

all: build

build:
	xgo -targets=linux/amd64 -ldflags="-w -s" -out=./build/telegram-monitor -pkg=cmd/telegram-monitor/main.go .
	upx ./build/telegram-monitor-${ARCH}

	xgo -targets=linux/amd64 -ldflags="-w -s" -out=./build/telegram-scanner -pkg=cmd/telegram-scanner/main.go .
	upx ./build/telegram-scanner-${ARCH}

	tar czf build/telegram-monitor-scanner-${dateTime}.tar.gz \
		build/telegram-monitor-${ARCH} \
		build/telegram-scanner-${ARCH} \
 		template \
 		img \
 		telegram-monitor.yaml.example \
 		telegram-scanner.yaml.example