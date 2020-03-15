package smtp

import (
	"crypto/tls"
	"net"
	"net/smtp"
)

func SendMailSSL(addr, password, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	host, _, _ := net.SplitHostPort(addr)
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()
	if password != "" {
		a := smtp.PlainAuth("", from, password, host)
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
