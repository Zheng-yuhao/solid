package main

import (
	"fmt"
	"solid/internal/repository/mysql"
	"solid/paste"
	"solid/user"
)

func main() {
	// SRP
	userRepo := mysql.NewUserRepo()
	userSvc := user.NewUserService(userRepo)
	user, err := userSvc.GetByID("1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("User: %+v\n", user)

	/*
		OCP
		・md5.goとsha256.goを確認すること
		　・registerで各アルゴリズムを登録している
		　・ShortLinkGeneratorFactoryでアルゴリズムを選択している
	*/
	generator := paste.ShortLinkGeneratorFactory("md5", "http://example.com")
	fmt.Println("short-link:", generator.GenerateShortLink())
}
