package main

import (
        "testing"
)

type FakeDbMap struct {
}

func (fake FakeDbMap) Insert(list ...interface{}) error {
  // // left in if it's needed later
  // v, ok := list[0].(*User)
	// if !ok {
	// 		log.Fatalf("not of type User!")
	// }
  //
  // log.Println(v.Password)

  return nil
}

func TestSave(t *testing.T) {
  DbMap := FakeDbMap{}
  u1 := User{Password: "tester"}
  u1.save(DbMap)

  u2 := User{Password: "tester"}
  u2.save(DbMap)

  if u1.Password == u2.Password {
    t.Fatalf("Expected the hashes to be different")
  }
}

func TestCheckPassword(t *testing.T) {
  DbMap := FakeDbMap{}
  u1 := User{Password: "tester"}
  u1.save(DbMap)

  if !u1.checkPassword("tester") {
    t.Fatalf("Passwords should match")
  }

  if u1.checkPassword("asdgth") {
    t.Fatalf("Passwords should not match")
  }
}

func TestJSON(t *testing.T) {
  u1 := User{Username: "ggap",
    Password: "tester",
    Email: "asdf@gasdg.com",
    Salt: []byte{0x01,0x02,0x03}} // Note: Json should not contain the salt

  json := u1.JSON()

  if json != "{\"Id\":0,\"email\":\"asdf@gasdg.com\",\"username\":\"ggap\",\"password\":\"tester\"}" {
    t.Fatalf("Json did not match")
  }
}
