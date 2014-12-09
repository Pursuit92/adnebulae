all: srcpackage

srcpackage: masters_documentation.pdf
	mkdir -p /tmp/adnebulae
	GOPATH=/tmp/adnebulae go get github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOARCH=386 go install github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOOS=windows go install github.com/Pursuit92/adnebulae/an-cli
	GOPATH=/tmp/adnebulae GOARCH=386 GOOS=windows go install github.com/Pursuit92/adnebulae/an-cli
	find /tmp/adnebulae -name '.git*' -exec rm -rf {} \+
	rm -rf /tmp/adnebulae/pkg /tmp/adnebulae/bin
	cp masters_documentation.pdf /tmp/JoshChase-AdNebulae.pdf
	cd /tmp && zip -e -P CLOUD -mr JoshChase-AdNebulae-PasswordIsCLOUD.zip adnebulae JoshChase-AdNebulae.pdf

masters_documentation.pdf:
	pdflatex masters_documentation.tex
