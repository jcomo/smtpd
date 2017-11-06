package smtpd

import (
	"net/smtp"
	"testing"
)

func mustGetClient(t *testing.T) *smtp.Client {
	c, err := smtp.Dial("127.0.0.1:8025")
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestSendMail(t *testing.T) {
	c := mustGetClient(t)

	err := c.Hello("testing")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Mail("jonathan.como@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Rcpt("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Rcpt("random@example.com")
	if err != nil {
		t.Fatal(err)
	}

	w, err := c.Data()
	if err != nil {
		t.Fatal(err)
	}

	msg := `
	  Dear Sir or Madam,

	  This is a message for you from the testing environment.
	  Thank you for your time.

	  Best,
	  Jonathan
	`

	w.Write([]byte(msg))
	w.Close()

	err = c.Quit()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRset(t *testing.T) {
	c := mustGetClient(t)

	c.Hello("TestRset")
	c.Mail("test@example.com")

	err := c.Reset()
	if err != nil {
		t.Fatal(err)
	}

	c.Quit()
}

func TestMailError(t *testing.T) {
	c := mustGetClient(t)

	c.Hello("TestMailError")

	err := c.Mail("test@")
	assertErrorEquals(t, "553 mail: invalid string", err)
}

func TestRcptError(t *testing.T) {
	c := mustGetClient(t)

	c.Hello("TestMailError")
	c.Mail("test@example.com")

	err := c.Rcpt("test@")
	assertErrorEquals(t, "553 mail: invalid string", err)
}

func TestSequenceError(t *testing.T) {
	c := mustGetClient(t)

	_, err := c.Data()
	assertErrorEquals(t, "503 bad command sequence", err)
}

func TestNotImplementedError(t *testing.T) {
	c := mustGetClient(t)

	err := c.Verify("test@example.com")
	assertErrorEquals(t, "502 not implemented", err)
}

func assertErrorEquals(t *testing.T, want string, got error) {
	if got == nil {
		t.Error("error expected but got nil")
		return
	}

	if want != got.Error() {
		t.Errorf("wanted: %s, got: %s\n", want, got.Error())
	}
}
