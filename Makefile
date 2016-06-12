default: testlocation
	docker build -f Dockerfile -t testlocation .

testlocation:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testlocation

clean:
	rm testlocation
