package datastore

import "testing"

func TestNewUserRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewUserRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}

func TestUserRepoCreate(t *testing.T) {
	db, userRepo := userRepoTestDeps()
	defer db.Close()

	user := &User{
		Username: "test",
		Password: "abc123",
	}

	user, err := userRepo.Create(user)
	if err != nil {
		t.Error(err)
	}

	if user.Created.IsZero() {
		t.Error("expected", ".Created to be set in database")
		t.Error("got     ", user.Created)
	}

	if user.ID == 0 {
		t.Error("expected", ".ID to be set in database")
		t.Error("got     ", user.ID)
	}
}

func TestUserRepoUpdate(t *testing.T) {
	db, userRepo := userRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	user, err := userRepo.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}

	user.Username = "updated"

	user, err = userRepo.Update(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepoGet(t *testing.T) {
	db, userRepo := userRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	user, err := userRepo.GetByID(1)
	if err != nil {
		t.Error(err)
	}
	if user.Username != "testuser1" {
		t.Error("expected", ".Username to be 'testuser1'")
		t.Error("got     ", user.Username)
	}

	user, err = userRepo.GetByUsername("testuser1")
	if err != nil {
		t.Error(err)
	}
	if user.ID != 1 {
		t.Error("expected", ".ID to equal 1")
		t.Error("got     ", user.ID)
	}

	_, err = userRepo.GetByUsername("nouser")
	if err == nil {
		t.Error("expect error when non existent")
	}

	users, err := userRepo.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(users) != 2 {
		t.Error("expected", "total of 2 users")
		t.Error("got     ", len(users))
	}
}
