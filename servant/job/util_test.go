// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package job

import (
	"testing"
)

func TestParamSplit_A(t *testing.T) {
	m := "A B \"C D E\" F G \"H I J\" K"
	params := paramSplit(m)
	if params[0] != "A" {
		t.Errorf("Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B" {
		t.Errorf("Bが返るべきところ、%vが返りました。", m[1])
	}
	if params[2] != "C D E" {
		t.Errorf("C D Eが返るべきところ、%vが返りました。", m[2])
	}
	if params[3] != "F" {
		t.Errorf("Fが返るべきところ、%vが返りました。", m[3])
	}
	if params[4] != "G" {
		t.Errorf("Gが返るべきところ、%vが返りました。", m[4])
	}
	if params[5] != "H I J" {
		t.Errorf("H I Jが返るべきところ、%vが返りました。", m[5])
	}
	if params[6] != "K" {
		t.Errorf("Kが返るべきところ、%vが返りました。", m[6])
	}
}

func TestParamSplit_B(t *testing.T) {
	m := "\"A B"
	params := paramSplit(m)
	if params[0] != "\"A" {
		t.Errorf("\"Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B" {
		t.Errorf("Bが返るべきところ、%vが返りました。", m[1])
	}
}

func TestParamSplit_C(t *testing.T) {
	m := "A B\""
	params := paramSplit(m)
	if params[0] != "A" {
		t.Errorf("Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B\"" {
		t.Errorf("B\"が返るべきところ、%vが返りました。", m[1])
	}
}

func TestParamSplit_D(t *testing.T) {
	m := "A B\"C D"
	params := paramSplit(m)
	if params[0] != "A" {
		t.Errorf("Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B\"C" {
		t.Errorf("B\"Cが返るべきところ、%vが返りました。", m[1])
	}
	if params[2] != "D" {
		t.Errorf("Dが返るべきところ、%vが返りました。", m[2])
	}
}

func TestParamSplit_E(t *testing.T) {
	m := "A B\"CD \" "
	params := paramSplit(m)
	if params[0] != "A" {
		t.Errorf("Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B\"CD \"" {
		t.Errorf("B\"CD \"が返るべきところ、%vが返りました。", m[1])
	}
}

func TestParamSplit_F(t *testing.T) {
	m := "A \"B C D\" E F\" G\""
	params := paramSplit(m)
	if params[0] != "A" {
		t.Errorf("Aが返るべきところ、%vが返りました。", m[0])
	}
	if params[1] != "B C D" {
		t.Errorf("B C Dが返るべきところ、%vが返りました。", m[1])
	}
	if params[2] != "E" {
		t.Errorf("Eが返るべきところ、%vが返りました。", m[2])
	}
	if params[3] != "F\" G\"" {
		t.Errorf("F\" G\"が返るべきところ、%vが返りました。", m[3])
	}
}

func TestShellFormat_Exist(t *testing.T) {
	s := shellFormat("A B C")
	if s != "\"A B C\"" {
		t.Errorf("Invalid Format - %v", s)
	}
}

func TestShellFormat_NoExist(t *testing.T) {
	s := shellFormat("ABC")
	if s != "ABC" {
		t.Errorf("Invalid Format - %v", s)
	}
}
