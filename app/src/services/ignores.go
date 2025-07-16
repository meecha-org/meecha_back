package services

import (
	"errors"
	"new-meecha/models"
	"slices"
)

//除外ポイントの作成、更新
func UpdateIgnores(myid string,args []models.IgnoresArgs) error {
	for _, arg := range args {
		if !slices.Contains(ValidationList, arg.Size) {
			return errors.New("InvalidDistance")
		}
    }
	
	//myidに関連する除外ポイントを削除
	err := models.RemoveIgnores(myid)
	if err != nil {
		return err
	}

	if args != nil {
		//送られてきたポイントを登録する
		err = models.SaveIgnores(myid,args)
		if err != nil {
			return err
		}
	}

	return nil
}


//除外ポイントの作成、更新
func GetIgnores(myid string)([]models.Ignores,error){
	//myidに関連する除外ポイントを取得
	ignores,err := models.GetIgnores(myid)
	if err != nil {
		return []models.Ignores{},err
	}

	return ignores,nil
}

