package apis

import "testing"

func TestFileStoreTableInsertion(t *testing.T) {
	TestSetupDbConnection(t)
	id1, err1 := addNewFileToStoreDB(23, "doctor", "a/b/c")
	if err1 != nil {
		t.Errorf("Unable to insert into %v", fileStoreTbl)
	}
	//Add duplicate path entry
	id2, err2 := addNewFileToStoreDB(23, "doctor", "a/b/c")
	if err2 != nil {
		t.Errorf("Unable to insert duplicate into %v", fileStoreTbl)
	}
	if id1 != id2 {
		t.Errorf("Bad ID returned %v(id1) != %v(id2)", id1, id2)
	}
}

/** TODO: Write better test case of getLoginData **/
func TestFetchSessionData(t *testing.T) {
	TestSetupDbConnection(t)
	_, _, _, err := getLoginData("123")
	if err != nil {
		t.Errorf("No login data found")
	}
}
