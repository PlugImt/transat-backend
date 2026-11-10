package photobooth

import "testing"

// The booth app pins the same vector (TransatShareApiTest): both signers must keep agreeing.
func TestSignatureVectorSharedWithTheBooth(t *testing.T) {
	keys, err := ParseDeviceKeys("booth-1:" + testSecret)
	if err != nil {
		t.Fatal(err)
	}
	c := Canonical("post", "/photobooth/booth/sessions/s1/finalize", "1800000000", "nonce-aaaaaaaaaaaaaaaa", []byte(`{"photo_count":3}`))
	if got := Sign(keys["booth-1"], c); got != boothVector {
		t.Fatalf("signature %s, want %s", got, boothVector)
	}
}

const boothVector = "ZSBBT88kJFyWXylOzKAQ63XxQga5e_UfZk8s_dPkuVY"
