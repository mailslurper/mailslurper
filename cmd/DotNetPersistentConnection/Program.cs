using System;
using System.Net.Mail;

namespace DotNetPersistentConnection {
	class Program {
		static void Main(string[] args) {
			Console.WriteLine("MailSlurper persistent connection test from .NET");

			SmtpClient client = new SmtpClient("localhost", 2500);

			client.Send(new MailMessage(
				"test@test.com",
				"test@test.com",
				"Persistent Connection Test .NET",
				"This is a test message. This tests how persistent connections work from a .NET client"
			));
		}
	}
}
