package hotp

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"math"
	"net/url"
	"strings"
)

type Digits int

const (
	DigitsSix   Digits = 6
	DigitsEight Digits = 8
)

func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}

type Algorithm int

const (
	AlgorithmSHA1 Algorithm = iota
	AlgorithmSHA256
	AlgorithmSHA512
	AlgorithmMD5
)

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA1:
		return "SHA1"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	case AlgorithmMD5:
		return "MD5"
	}
	panic("unreached")
}

func (a Algorithm) Hash() hash.Hash {
	switch a {
	case AlgorithmSHA1:
		return sha1.New()
	case AlgorithmSHA256:
		return sha256.New()
	case AlgorithmSHA512:
		return sha512.New()
	case AlgorithmMD5:
		return md5.New()
	}
	panic("unreached")
}

type Key struct {
	orig string
	url  *url.URL
}

func NewKeyFromURL(orig string) (*Key, error) {
	s := strings.TrimSpace(orig)

	u, err := url.Parse(s)

	if err != nil {
		return nil, err
	}

	return &Key{
		orig: s,
		url:  u,
	}, nil
}

func (k *Key) String() string {
	return k.orig
}

const debug = false

// Validate a HOTP passcode given a counter and secret.
// This is a shortcut for ValidateCustom, with parameters that
// are compataible with Google-Authenticator.
func Validate(passcode string, counter uint64, secret string) bool {
	rv, _ := ValidateCustom(
		passcode,
		counter,
		secret,
		ValidateOpts{
			Digits:    DigitsSix,
			Algorithm: AlgorithmSHA1,
		},
	)
	return rv
}

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	Digits    Digits
	Algorithm Algorithm
}

// GenerateCodeCustom uses a counter and secret value and options struct to
// create a passcode.
func GenerateCodeCustom(secret string, counter uint64, opts ValidateOpts) (passcode string, err error) {
	// As noted in issue #10 and #17 this adds support for TOTP secrets that are
	// missing their padding.
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	// As noted in issue #24 Google has started producing base32 in lower case,
	// but the StdEncoding (and the RFC), expect a dictionary of only upper case letters.
	secret = strings.ToUpper(secret)

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", errors.New("Decoding of secret as base32 failed.")
	}

	buf := make([]byte, 8)

	mac := hmac.New(opts.Algorithm.Hash, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)
	if debug {
		fmt.Printf("counter=%v\n", counter)
		fmt.Printf("buf=%v\n", buf)
	}

	mac.Write(buf)
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits.Length()
	mod := int32(value % int64(math.Pow10(l)))

	if debug {
		fmt.Printf("offset=%v\n", offset)
		fmt.Printf("value=%v\n", value)
		fmt.Printf("mod'ed=%v\n", mod)
	}

	return opts.Digits.Format(mod), nil
}

// ValidateCustom validates an HOTP with customizable options. Most users should
// use Validate().
func ValidateCustom(passcode string, counter uint64, secret string, opts ValidateOpts) (bool, error) {
	passcode = strings.TrimSpace(passcode)

	if len(passcode) != opts.Digits.Length() {
		return false, errors.New("Input length unexpected")
	}

	otpstr, err := GenerateCodeCustom(secret, counter, opts)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(otpstr), []byte(passcode)) == 1 {
		return true, nil
	}

	return false, nil
}

// GenerateOpts provides options for .Generate()
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Size in size of the generated Secret. Defaults to 10 bytes.
	SecretSize uint
	// Secret to store. Defaults to a randomly generated secret of SecretSize.  You should generally leave this empty.
	Secret []byte
	// Digits to request. Defaults to 6.
	Digits Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm Algorithm
	// Reader to use for generating HOTP Key.
	Rand io.Reader
}

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

// Generate creates a new HOTP Key.
//func Generate(opts GenerateOpts) (*Key, error) {
//	// url encode the Issuer/AccountName
//	if opts.Issuer == "" {
//		return nil, errors.New("Issuer must be set")
//	}
//
//	if opts.AccountName == "" {
//		return nil, errors.New("AccountName must be set")
//	}
//
//	if opts.SecretSize == 0 {
//		opts.SecretSize = 10
//	}
//
//	if opts.Digits == 0 {
//		opts.Digits = DigitsSix
//	}
//
//	if opts.Rand == nil {
//		opts.Rand = rand.Reader
//	}
//
//	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
//
//	v := url.Values{}
//	if len(opts.Secret) != 0 {
//		v.Set("secret", b32NoPadding.EncodeToString(opts.Secret))
//	} else {
//		secret := make([]byte, opts.SecretSize)
//		_, err := opts.Rand.Read(secret)
//		if err != nil {
//			return nil, err
//		}
//		v.Set("secret", b32NoPadding.EncodeToString(secret))
//	}
//
//	v.Set("issuer", opts.Issuer)
//	v.Set("algorithm", opts.Algorithm.String())
//	v.Set("digits", opts.Digits.String())
//
//	u := url.URL{
//		Scheme:   "otpauth",
//		Host:     "hotp",
//		Path:     "/" + opts.Issuer + ":" + opts.AccountName,
//		RawQuery: v.Encode(),
//	}
//
//	return NewKeyFromURL(u.String())
//}
