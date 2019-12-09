package hotp

import (
	"testing"
)

func TestHOTPValidateWithValidOTP(t *testing.T) {
	if !Validate("644039", 2, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP refused valid OTP")
	}

	if !Validate("586787", 3, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP refused valid OTP")
	}
}

func TestHOTPValidateWithInvalidOTP(t *testing.T) {
	if Validate("644039", 3, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP accepted invalid OTP")
	}

	if Validate("586787", 2, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP accepted invalid OTP")
	}
}
