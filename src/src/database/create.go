package database

import "src/rstring"

func UniqueUCode(str string) bool {
	if _, ok := UsedUC[str]; ok {
		return true
	}
	return false
}

func CreateUCode() string {
	var str string
	for str == "" || UniqueUCode(str) {
		str = rstring.GetRandString(6)
	}

	UsedUC[str] = true
	return str
}

func CreateDCode() string {
	return rstring.GetRandString(6)
}
