package domain

import "gorm.io/gorm"

// "from":{"id":205478553,"is_bot":false,"first_name":"Guionardo","last_name":"Furlan","username":"guionardo","language_code":"en"}

type(
	User struct{
		gorm.Model
		ID int
		FirstName string
		LastName string
		UserName string
		LanguageCode string
		Associados []*MappaAssociado `gorm:"many2many:user_associados;"`
	}
)
