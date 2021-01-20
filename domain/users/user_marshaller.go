package users

import "encoding/json"

type PublicUser struct {
	Id int64 `json:"id"`
	//FirstName  string `json:"first_name"`
	//LastName   string `json:"last_name"`
	//Email      string `json:"email"`
	DateCreate string `json:"date_create"`
	Status     string `json:"status"`
	//Password   string `json:"password"`
}

type PrivateUser struct {
	Id         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	DateCreate string `json:"date_create"`
	Status     string `json:"status"`
	//Password   string `json:"password"`
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:         user.Id,
			DateCreate: user.DateCreate,
			Status:     user.Status,
		}
	}
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}
