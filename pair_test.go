package cutedb

import "testing"

func TestValidate(t *testing.T) {
	key1:="1234567890123456789012345678901"
	value:="jffju1rig13rg3r2yg3g3r1t78r3t37t3r273r27f73r2g73r2ff3rf3r27"+
	"f3rasnjjnasdnjadsjndsanjsdanjdsanjdsjasndsjdjnsdajnsdajnsjndsajn"+
	"jasnadsaasdkhadhyy727t22effawfawhhdahsahgasdhgdsahgsadhgsadhghdghads"+
	"asdjajsddsjajsdajsadjndsjn"
	pair:=newPair(key1,value)

	if pair.validate() == nil {
		t.Errorf("Should throw error")
	}
	if newPair(key1,"ss").validate() == nil {
		t.Errorf("Shoudl throw error as key is longer than 30")
	}
	if newPair("smallKEY",value).validate() == nil {
		t.Errorf("Shoudl throw error as value is longer than 90")
	}
}