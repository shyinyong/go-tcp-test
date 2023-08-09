package utils

import (
	"testing"
)

func TestToInt(t *testing.T) {
	// int
	var input int = 50
	var expected int = 50
	result := ToInt(input)
	if expected != result {
		t.Errorf("For input %v, expected %d, but got %d", input, expected, result)
	}

	// int64
	var input2 int64 = 50000
	var expected2 int = 50000
	result2 := ToInt(input2)
	if expected2 != result2 {
		t.Errorf("For input %v, expected %d, but got %d", input2, expected2, result2)
	}

	// float
	var input3 float64 = 3.14
	var expected3 int = 3
	result3 := ToInt(input3)
	if expected3 != result3 {
		t.Errorf("For input %v, expected %d, but got %d", input3, expected3, result3)
	}

}

func TestToInt64(t *testing.T) {
	// int
	var input int = 50
	var expected int64 = 50
	result := ToInt64(input)
	if expected != result {
		t.Errorf("For input %v, expected %v, but got %v", input, expected, result)
	}

	// int64
	var input2 int64 = 50000
	var expected2 int64 = 50000
	result2 := ToInt64(input2)
	if expected2 != result2 {
		t.Errorf("For input %v, expected %v, but got %v", input2, expected2, result2)
	}

	// float
	var input3 float64 = 3.14
	var expected3 int64 = 3
	result3 := ToInt64(input3)
	if expected3 != result3 {
		t.Errorf("For input %v, expected %v, but got %v", input3, expected3, result3)
	}

	// string
	var input4 string = "100108454613"
	var expected4 int64 = 100108454613
	result4 := ToInt64(input4)
	if expected4 != result4 {
		t.Errorf("For input %v, expected %v, but got %v", input4, expected4, result4)
	}
}

func TestToFloat(t *testing.T) {
	// int
	var input int = 50
	var expected float64 = 50
	result := ToFloat(input)
	if expected != result {
		t.Errorf("For input %v, expected %v, but got %v", input, expected, result)
	}

	// int64
	var input2 int64 = 50000
	var expected2 float64 = 50000
	result2 := ToFloat(input2)
	if expected2 != result2 {
		t.Errorf("For input %v, expected %v, but got %v", input2, expected2, result2)
	}

	// float
	var input3 float64 = 3.14
	var expected3 float64 = 3.14
	result3 := ToFloat(input3)
	if expected3 != result3 {
		t.Errorf("For input %v, expected %v, but got %v", input3, expected3, result3)
	}

	// string
	var input4 string = "100108454613"
	var expected4 float64 = 100108454613
	result4 := ToFloat(input4)
	if expected4 != result4 {
		t.Errorf("For input %v, expected %v, but got %v", input4, expected4, result4)
	}
}

func TestToString(t *testing.T) {
	// int
	var input int = 50
	var expected string = "50"
	result := ToString(input)
	if expected != result {
		t.Errorf("For input %v, expected %s, but got %s", input, expected, result)
	}

	// int64
	var input2 int64 = 50000
	expected2 := "50000"
	result2 := ToString(input2)
	if expected2 != result2 {
		t.Errorf("For input %v, expected %v, but got %v", input2, expected2, result2)
	}

	// float
	var input3 float64 = 3.14
	var expected3 string = "3.14"
	result3 := ToString(input3)
	if expected3 != result3 {
		t.Errorf("For input %v, expected %v, but got %v", input3, expected3, result3)
	}

	// string
	input4 := "apple"
	expected4 := "apple"
	result4 := ToString(input4)
	if expected4 != result4 {
		t.Errorf("For input %v, expected %v, but got %v", input4, expected4, result4)
	}
}
