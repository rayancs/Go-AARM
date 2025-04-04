package test

import (
	"app/configs"
	"app/logger"
	repo "app/repos"
)

// manual test
func Test() {
	lg := logger.New()
	lg.Warning("\n***********************************************\nTesing is going to start,enter 1 to continue,\nanything else will skip the tests\n***********************************************\n\n")
	db := repo.NewMongoDB(configs.GetMongoDBURI(), "dev-shared")
	authRepo := repo.NewMongoUser(db)
	authTest := NewAuthTest(authRepo)
	authTest.insetRepoTest()
	authTest.TestAuthRepo()
	lg.Warning("\nAll Tests Performed\n\n\n")
}
