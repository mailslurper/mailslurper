import socket
import sys

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
serverAddress = ("localhost", 2500)

sock.connect(serverAddress)

try:

	# Say hi
	print "Saying hi"
	sock.sendall("EHLO localhost")

	data = sock.recv(2048)
	print data

	# From
	print "Sending from"
	sock.sendall("MAIL FROM: test@test.com")

	data = sock.recv(2048)
	print data

	# To
	print "Sending to"
	sock.sendall("RCPT TO: adam@test.com")

	data = sock.recv(2048)
	print data

	# Send  a malformed DATA
	print "sending data"
	sock.sendall("DATA\r\nFrom:test@test.com\r\nTo: adam@test.com\r\n\r\n\r\n.\r\n")

	data = sock.recv(2048)
	print data

finally:
	print "closing"
	sock.close()
