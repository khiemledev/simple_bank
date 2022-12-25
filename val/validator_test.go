package val

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	email := "khiemledev@gmail.com"
	err := ValidateEmail(email)
	require.NoError(t, err)

	email = "khiemle @2409"
	err = ValidateEmail(email)
	require.Error(t, err)

	email = ""
	err = ValidateEmail(email)
	require.Error(t, err)
}

func TestValidatePassword(t *testing.T) {
	password := "this_is_my_password123"
	err := ValidatePassword(password)
	require.NoError(t, err)

	password = "sort"
	err = ValidatePassword(password)
	require.Error(t, err)
}

func TestValidateFullName(t *testing.T) {
	fullName := "This is my name"
	err := ValidateFullName(fullName)
	require.NoError(t, err)

	fullName = "s"
	err = ValidateFullName(fullName)
	require.Error(t, err)
}
