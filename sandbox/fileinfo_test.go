package main

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
)

func TestGetFileInfo(t *testing.T) {
	aFileInfo := getFileInfo("test.txt")
	if len(aFileInfo) >= 1 {
		// ok
	} else {
		t.Fatalf("failed test %#v", aFileInfo)
	}
}

// func TestExampleFailed(t *testing.T) {
// 	result, err := example("fuga")
// 	if err == nil {
// 		t.Fatal("failed test")
// 	}
// 	if result != 0 {
// 		t.Fatal("failed test")
// 	}
// }
