.PHONY: cudy
cudy:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o jg_linux_arm64 && \
	scp -O main.js jg_linux_arm64 cudy:/tmp/ && \
	ssh cudy 'cd /tmp && ./jg_linux_arm64'
