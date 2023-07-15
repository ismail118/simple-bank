package mail

import "testing"

func Test_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	subject := "A test email"
	content := `
	<h1>hello world</h1>
	<p>This is test email</p>
`
	to := []string{"ismailalfiyasin6@gmail.com"}
	attachFiles := []string{"../README.md"}

	err := SenderEmailTest.SendEmail(subject, content, to, nil, nil, attachFiles)
	if err != nil {
		t.Fatalf("failed send email: %s", err)
	}
}
