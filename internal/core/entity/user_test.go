package entity

import "testing"

func TestUser_SetPassword_Hashes(t *testing.T) {
	u := &User{}
	if err := u.SetPassword("pw123"); err != nil {
		t.Fatal(err)
	}
	if u.Password == "" || u.Password == "pw123" {
		t.Fatalf("password not hashed: %q", u.Password)
	}
}

func TestUser_CheckPassword(t *testing.T) {
	u := &User{}
	if err := u.SetPassword("pw123"); err != nil {
		t.Fatal(err)
	}
	if !u.CheckPassword("pw123") {
		t.Fatal("correct password should verify")
	}
	if u.CheckPassword("wrong") {
		t.Fatal("wrong password must not verify")
	}
}

func TestUser_SetPassword_TooLongRejected(t *testing.T) {
	// bcrypt caps input at 72 bytes; GenerateFromPassword returns an error above that.
	u := &User{}
	long := make([]byte, 73)
	for i := range long {
		long[i] = 'a'
	}
	if err := u.SetPassword(string(long)); err == nil {
		t.Fatal("want error for >72-byte password")
	}
}
