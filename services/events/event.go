package events

type Event struct {
	Id string `json:"id" pg:"type=uuid,primary"`
	//CreatedAt   time.Time `json:"createdAt" pg:"type=timestamp,default=clock_timestamp()"`
	//UpdatedAt   time.Time `json:"updatedAt" pg:"type=timestamp,default=clock_timestamp()"`
	Title       string  `json:"title" pg:"type=varchar(255),unique"`
	Description *string `json:"description" pg:"type=varchar(4096),null"`
}
