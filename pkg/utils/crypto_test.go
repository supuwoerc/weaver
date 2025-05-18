package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type cryptoArgs struct {
	s string
}
type cryptoCase struct {
	name string
	args cryptoArgs
	want string
}

func testCryptoFunc(t *testing.T, tests []cryptoCase, fn func(s string) string) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, fn(tt.args.s), "param string:%v", tt.args.s)
		})
	}
}

func TestMd5(t *testing.T) {
	cases := []cryptoCase{
		{
			name: "MD5(123)",
			args: cryptoArgs{
				s: "123",
			},
			want: "202cb962ac59075b964b07152d234b70",
		},
		{
			name: "MD5(abc)",
			args: cryptoArgs{
				s: "abc",
			},
			want: "900150983cd24fb0d6963f7d28e17f72",
		},
	}
	testCryptoFunc(t, cases, Md5)
}

func TestSha1(t *testing.T) {
	cases := []cryptoCase{
		{
			name: "Sha1(123)",
			args: cryptoArgs{
				s: "123",
			},
			want: "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
		},
		{
			name: "Sha1(abc)",
			args: cryptoArgs{
				s: "abc",
			},
			want: "a9993e364706816aba3e25717850c26c9cd0d89d",
		},
	}
	testCryptoFunc(t, cases, Sha1)
}

func TestSha256(t *testing.T) {
	cases := []cryptoCase{
		{
			name: "Sha256(123)",
			args: cryptoArgs{
				s: "123",
			},
			want: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		},
		{
			name: "Sha256(abc)",
			args: cryptoArgs{
				s: "abc",
			},
			want: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
	}
	testCryptoFunc(t, cases, Sha256)
}

func TestSha512(t *testing.T) {
	cases := []cryptoCase{
		{
			name: "Sha512(123)",
			args: cryptoArgs{
				s: "123",
			},
			want: `3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2`,
		},
		{
			name: "Sha512(abc)",
			args: cryptoArgs{
				s: "abc",
			},
			want: `ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f`,
		},
	}
	testCryptoFunc(t, cases, Sha512)
}
