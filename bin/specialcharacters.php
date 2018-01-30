<?php
use PHPMailer\PHPMailer\PHPMailer;

date_default_timezone_set('Etc/UTC');
require 'PHPMailer/Exception.php';
require 'PHPMailer/OAuth.php';
require 'PHPMailer/SMTP.php';
require 'PHPMailer/PHPMailer.php';

$mail = new PHPMailer;
$mail->isSMTP();
$mail->SMTPDebug = 1;
$mail->CharSet = 'UTF-8';
$mail->Host = 'localhost';
$mail->Port = 2500;
$mail->SMTPAuth = false;
$mail->setFrom('test@example.com', 'First Last');
$mail->addAddress('john@doe.com', 'John Doe');
$mail->Subject = 'Modification activité';
$mail->msgHTML("Ceci est mon contenu accentué éàçè");

if (!$mail->send()) {
    echo 'Mailer Error: ' . $mail->ErrorInfo;
} else {
    echo 'Message sent!';
}
?>
