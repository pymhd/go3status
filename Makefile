all: 
	go build -o go3status ./

install: uninstall
	cp go3status /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/go3status

clean:
	rm -vrf go3status
