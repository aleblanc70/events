package users

import "time"

type User struct {
	Id         string    `json:"id" pg:"type=uuid,primary"`
	CreatedAt  time.Time `json:"createdAt" pg:"type=timestamp,default=clock_timestamp()"`
	UpdatedAt  time.Time `json:"updatedAt" pg:"type=timestamp,default=clock_timestamp()"`
	Email      string    `json:"email" pg:"type=varchar(255),username,unique"`
	Password   string    `json:"password" pg:"type=varchar(128),password"`
	FirstName  string    `json:"firstName" pg:"type=varchar(255)"`
	MiddleName *string   `json:"middleName" pg:"type=varchar(255),null"`
	LastName   string    `json:"lastName" pg:"type=varchar(255)"`
	Sex        *string   `json:"sex" pg:"type=char(1),null"`
}
