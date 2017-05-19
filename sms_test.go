package zang

import (
	"testing"
	"fmt"
	"os"
)

func TestZangSender(t *testing.T) {
	c := NewClient(os.Getenv("ZANG_SID"),os.Getenv("ZANG_TOKEN"))
	response, err :=c.Send(Message{
		From:PhoneNumber("+17184129523&"),
		To:PhoneNumber("5099948638"),
		Body:"Heyo",
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", response)
}