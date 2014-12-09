all: srcpackage

srcpackage: masters_documentation.pdf
	mkdir -p /tmp/adnebulae
	GOPATH=/tmp/adnebulae go get github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOARCH=386 go install github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOOS=windows go install github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOARCH=386 GOOS=windows go install github.com/Pursuit92/adnebulae/an-cli
	find /tmp/adnebulae -name '.git*' -exec rm -rf {} \+
	rm -rf /tmp/adnebulae/pkg
	cp masters_documentation.pdf /tmp/adnebulae/JoshChase-AdNebulae.pdf
	cd /tmp && zip -mr adnebulae.zip adnebulae

masters_documentation.pdf:
	pdflatex masters_documentation.tex
