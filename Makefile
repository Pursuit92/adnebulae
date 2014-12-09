all: srcpackage

srcpackage:
	mkdir -p /tmp/adnebulae
	GOPATH=/tmp/adnebulae go get github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOOS=windows go install github.com/Pursuit92/adnebulae/an-cli
	find /tmp/adnebulae -name .git -exec rm -rf {} \+
	rm -rf /tmp/adnebulae/pkg
	cd /tmp && zip -mr adnebulae.zip adnebulae
