package auth

import (
	"testing"
)

func TestValidateHOTPWithValidOTP(t *testing.T) {
	if !ValidateHOTP("644039", 2, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP refused valid OTP")
	}

	if !ValidateHOTP("586787", 3, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP refused valid OTP")
	}
}

func TestValidateHOTPWithInvalidOTP(t *testing.T) {
	if ValidateHOTP("644039", 3, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP accepted invalid OTP")
	}

	if ValidateHOTP("586787", 2, "DFIOJ3BJPHPCRJBT") {
		t.Errorf("HOTP accepted invalid OTP")
	}
}
