package cmn

import "testing"

func TestIsMoneyAmountValid(t *testing.T) {
	a := MaxBigInt
	tested := IsMoneyAmountValid(a)
	expected := false

	if expected != tested {
		t.Errorf(
			"Wrong validation, we shouldn't allow numbers above limit. Expected: %v, received (%v)\n",
			expected, tested)
	}

	a = -1 * MaxBigInt
	tested = IsMoneyAmountValid(a)
	expected = false

	if expected != tested {
		t.Errorf(
			"Wrong validation, we shouldn't allow numbers below limit. Expected: %v, received (%v)\n",
			expected, tested)
	}

	a = 0
	tested = IsMoneyAmountValid(a)
	expected = true

	if expected != tested {
		t.Errorf(
			"Wrong validation, should allow numbers within limit. Expected: %v, received (%v)\n",
			expected, tested)
	}

	a = 10
	tested = IsMoneyAmountValid(a)
	expected = true

	if expected != tested {
		t.Errorf(
			"Wrong validation, should allow numbers within limit. Expected: %v, received (%v)\n",
			expected, tested)
	}
	a = -10
	tested = IsMoneyAmountValid(a)
	expected = true

	if expected != tested {
		t.Errorf(
			"Wrong validation, should allow numbers within limit. Expected: %v, received (%v)\n",
			expected, tested)
	}

}
