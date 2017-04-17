package cache_test

import (
	"fmt"
	"os"
	"testing"

	//"gitlab.engr.illinois.edu/sp-box/boxsync/auth"
	//"gitlab.engr.illinois.edu/sp-box/boxsync/box"
	"gitlab.engr.illinois.edu/sp-box/boxsync/cache"
)

func TestBasicCache(t *testing.T) {
	os.Mkdir("testing_tmp", 0777)
	defer os.RemoveAll("testing_tmp/")

	os.Create("testing_tmp/test1")
	os.Create("testing_tmp/test2")
	os.Create("testing_tmp/test3")

	os.Mkdir("testing_tmp/dir1", 0777)

	os.Create("testing_tmp/dir1/test1")
	os.Create("testing_tmp/dir1/test2")
	os.Create("testing_tmp/dir1/test3")

	checks := map[string]bool{
		"test1":      false,
		"test2":      false,
		"test3":      false,
		"dir1/test1": false,
		"dir1/test2": false,
		"dir1/test3": false,
	}

	db := cache.InitCache(nil, "testing_tmp")

	rows, err := db.Query("select path from testing_tmp")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			t.Fatal(err)
		}
		if val, ok := checks[path]; ok {
			if val {
				//This should never happen
				fmt.Print("path found twice")
				t.Fail()
			}

			checks[path] = true
		}
	}

	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}

	for path, found := range checks {
		if !found {
			fmt.Printf("path %s not found\n", path)
			t.Fail()
		}
	}
}
