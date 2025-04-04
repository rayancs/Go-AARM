package test

import (
	"app/logger"
	repo "app/repos"
	"app/types"
)

type AuthTest struct {
	// will contian all authentication levels from repo to apis
	AuthRepo repo.IUserRepo
}

func NewAuthTest(A repo.IUserRepo) *AuthTest {
	return &AuthTest{
		AuthRepo: A,
	}
}
func (a *AuthTest) insetRepoTest() {
	lg := logger.NewTestLogger("Insert Auth Repo Test")
	lg.Info("Insert a user to the database'")
	insDoc := types.UserRequestType{
		Emoji:      "👻",
		Email:      "mianrayan211@gmail.com",
		Name:       "rayan",
		PictureURL: "n/a",
	}
	lg.Info("insert object = ", insDoc)
	res, err := a.AuthRepo.CreateUser(insDoc)
	if err != nil {
		lg.Fail("Insert Failed Due To Error", err)
	}
	if res == "" {
		lg.Fail("Insert failed partially , insert id is empty")
	}
	if res != "" {
		lg.Pass("test passed , insert id = ", res)
	}
	lg.Warning("TEST END")

}
func (a *AuthTest) TestAuthRepo() {
	lg := logger.NewTestLogger("Auth Repo Test")
	lg.Info("test case 'find user , valid email , user does not exsits' ")
	res, err := a.AuthRepo.FindUser("mianrayan211@gmail.com")
	if err != nil {
		lg.Fail("error at testing", err)
	}
	if res == nil {
		lg.Fail("User nil", res)
	}
	if res != nil {
		lg.Pass("User not nil , but error free ", res)
	}
}
