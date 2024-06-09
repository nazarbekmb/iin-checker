package iin

import (
	"fmt"
	"strconv"
)

type IINInfo struct {
	Correct     bool   `json:"correct"`
	Sex         string `json:"sex"`
	DateOfBirth string `json:"date_of_birth"`
}

func ValidateIIN(iin string) (IINInfo, error) {
	if len(iin) != 12 {
		return IINInfo{Correct: false}, fmt.Errorf("invalid length")
	}

	var digits [12]int
	for i := 0; i < 12; i++ {
		digit, err := strconv.Atoi(string(iin[i]))
		if err != nil {
			return IINInfo{Correct: false}, fmt.Errorf("invalid character")
		}
		digits[i] = digit
	}
	weight1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	weight2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}

	var sum1, sum2 int
	for i := 0; i < 11; i++ {
		sum1 += digits[i] * weight1[i]
		sum2 += digits[i] * weight2[i]
	}

	var checksum int
	if sum1%11 == 10 {
		checksum = sum2 % 11
	} else {
		checksum = sum1 % 11
	}

	if checksum == 10 || checksum != digits[11] {
		return IINInfo{Correct: false}, fmt.Errorf("invalid checksum")
	}

	yearPrefix := 0
	switch digits[6] {
	case 1, 2:
		yearPrefix = 1800
	case 3, 4:
		yearPrefix = 1900
	case 5, 6:
		yearPrefix = 2000
	default:
		return IINInfo{Correct: false}, fmt.Errorf("invalid year indicator")
	}

	year := yearPrefix + digits[0]*10 + digits[1]
	month := digits[2]*10 + digits[3]
	day := digits[4]*10 + digits[5]

	dateOfBirth := fmt.Sprintf("%02d.%02d.%04d", day, month, year)

	sex := "female"
	if digits[6]%2 != 0 {
		sex = "male"
	}

	return IINInfo{
		Correct:     true,
		Sex:         sex,
		DateOfBirth: dateOfBirth,
	}, nil
}
