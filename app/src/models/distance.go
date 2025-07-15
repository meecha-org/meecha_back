package models

import(

)

type Distance struct{
	UserId    string  
	Distance  int64   //通知距離 (メートル)
	CreatedAt int64   //作成時間 (UnixTime)
	UpdatedAt int64   //更新時間 (UnixTime)
}

func getdistance(){
	return
}