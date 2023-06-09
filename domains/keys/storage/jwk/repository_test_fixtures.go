package jwk_storage

import (
	"crypto/ed25519"
	"crypto/x509"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
)

func keyFromBytes(b []byte) ed25519.PrivateKey {
	key, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		panic(err.Error())
	}

	return key.(ed25519.PrivateKey)
}

var Fixtures = []test_utils.FileFixture{
	{
		Name: "foo-test-1",
		Content: []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIKDMtkBHhK5UMhI4PkEHyvX+hzKwTAnzw0xiKb5pNEeB
-----END PRIVATE KEY-----
`),
		Date: baseTime,
	},
	{
		Name: "foo-test-2",
		Content: []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIBKS2G+b8tSr/SE8xRm7IJXQocf8LrSltkp1OBnK5Ovv
-----END PRIVATE KEY-----
`),
		Date: baseTime.Add(time.Hour),
	},
	{
		Name: "foo-test-3",
		Content: []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIO778kSy0BejGR86oc/5VXijoue6FrmsZ761eaEfo2Uw
-----END PRIVATE KEY-----
`),
		Date: baseTime.Add(30 * time.Minute),
	},
	{
		Name: "bar-test-4",
		Content: []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIAGgCAAqob5UhxlrpnUSY6N5k6P8JZImpMY4wK025uba
-----END PRIVATE KEY-----
`),
		Date: baseTime.Add(15 * time.Minute),
	},
}

var MockedKeys = []ed25519.PrivateKey{
	// Fixtures[0]
	keyFromBytes([]byte{48, 46, 2, 1, 0, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32, 160, 204, 182, 64, 71, 132, 174, 84, 50, 18, 56, 62, 65, 7, 202, 245, 254, 135, 50, 176, 76, 9, 243, 195, 76, 98, 41, 190, 105, 52, 71, 129}),
	// Fixtures[1]
	keyFromBytes([]byte{48, 46, 2, 1, 0, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32, 18, 146, 216, 111, 155, 242, 212, 171, 253, 33, 60, 197, 25, 187, 32, 149, 208, 161, 199, 252, 46, 180, 165, 182, 74, 117, 56, 25, 202, 228, 235, 239}),
	// Fixtures[2]
	keyFromBytes([]byte{48, 46, 2, 1, 0, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32, 238, 251, 242, 68, 178, 208, 23, 163, 25, 31, 58, 161, 207, 249, 85, 120, 163, 162, 231, 186, 22, 185, 172, 103, 190, 181, 121, 161, 31, 163, 101, 48}),
	// Fixtures[3]
	keyFromBytes([]byte{48, 46, 2, 1, 0, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32, 1, 160, 8, 0, 42, 161, 190, 84, 135, 25, 107, 166, 117, 18, 99, 163, 121, 147, 163, 252, 37, 146, 38, 164, 198, 56, 192, 173, 54, 230, 230, 218}),
	// Original
	keyFromBytes([]byte{48, 46, 2, 1, 0, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32, 8, 51, 50, 127, 105, 192, 97, 124, 90, 97, 110, 77, 142, 185, 31, 51, 28, 178, 110, 231, 235, 74, 106, 171, 56, 64, 251, 121, 119, 11, 57, 44}),
}
