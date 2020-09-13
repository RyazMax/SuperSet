package auth

import (
	"testing"

	"../models"
	"../repos/session"
	"../repos/user"
)

var auth Auth
var host = "127.0.0.1"
var port = 6666

var testUser = models.User{}

func TestAuth(t *testing.T) {
	repos := []Auth{
		&SimpleAuth{},
	}

	ur := user.TarantoolRepo{}
	sr := session.TarantoolRepo{}
	err := ur.Init(host, port)
	if err != nil {
		t.Error("Cant init repo", err)
	}
	err = sr.Init(host, port)
	if err != nil {
		t.Error("Cant init repo", err)
	}

	for _, repo := range repos {
		err = repo.Init(&ur, &sr)
		if err != nil {
			t.Errorf("Init failed on %T with error %v", repo, err)
		}

		sess, err := repo.NewUser(&testUser)
		if err != nil {
			t.Errorf("NewUser failed on %T with error %v", repo, err)
		}
		if sess == nil {
			t.Errorf("NewUser failed on %T, got nil session", repo)
		}

		newSess, err := repo.CheckSession(sess.ID)
		if err != nil {
			t.Errorf("CheckSession after NewUser failed on %T, with error %v", repo, err)
		}
		if !sess.IsEqual(newSess) {
			t.Errorf("CheckSession after NewUser failed on %T, expected %v, got %v", repo, sess, newSess)
		}

		err = repo.Logout(sess.ID)
		if err != nil {
			t.Errorf("Logout failed on %T, with error %v", repo, err)
		}

		newSess, err = repo.CheckSession(sess.ID)
		if err != nil {
			t.Errorf("CheckSession after Logout failed on %T with error %v", repo, err)
		}
		if newSess != nil {
			t.Errorf("CheckSession after Logout failed on %T, expected nil", repo)
		}

		sess, err = repo.Login(testUser.Login, testUser.PasswordShadow)
		if err != nil {
			t.Errorf("Login failed on %T, with error %v", repo, err)
		}
		if sess == nil {
			t.Errorf("Login failed on %T, expected non nil session", repo)
		}

		newSess, err = repo.CheckSession(sess.ID)
		if err != nil {
			t.Errorf("CheckSession after Login failed on %T, with error %v", repo, err)
		}
		if !sess.IsEqual(newSess) {
			t.Errorf("CheckSession after Login failed on %T, expected %v, got %v", repo, sess, newSess)
		}

		err = repo.DeleteUser(testUser.Login)
		if err != nil {
			t.Errorf("DeleteUser failed on %T, with error %v", repo, err)
		}

		/*newSess, err = repo.CheckSession(sess.ID)
		if err != nil {
			t.Errorf("CheckSession after DeleteUser failed on %T with error %v", repo, err)
		}
		if newSess != nil {
			t.Errorf("CheckSession after DeleteUser failed on %T, expected nil", repo)
		}*/
		// TODO: Also can check that other session of that user was deleted

		newSess, err = repo.Login(testUser.Login, testUser.PasswordShadow)
		if err != nil {
			t.Errorf("Login after DeleteUser failed on %T, with error %v", repo, err)
		}
		if newSess != nil {
			t.Errorf("Login after DeleteUser failed on %T, expected nil session", repo)
		}

	}
}
