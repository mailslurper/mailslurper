openssl genrsa -out mailslurper-key.pem 1024
openssl req -new -key mailslurper-key.pem -out mailslurper.csr
openssl x509 -req -in mailslurper.csr -signkey mailslurper-key.pem -out mailslurper-cert.pem
