package services

import (
	"context"
	"errors"

	//"fmt"
	"regexp"
	"s7/interfaces"
	"s7/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUseService(usercollection *mongo.Collection, ctx context.Context) interfaces.UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}


func (u *UserServiceImpl) CreateUser(user *models.User) error {
	
	name := user.Name
	regexPattern := "^[a-zA-Z ]+$"
	regexpObject := regexp.MustCompile(regexPattern)
	match := regexpObject.FindString(name)

	
	existingUser, err := u.getUserByEmail(user.Email)
	if err != nil {
		return err 
	}
    if existingUser!=nil{
		return errors.New("user already exist")
	}
	phonecount:=u.countdigits(int(user.PhoneNumber))
	if phonecount<10{
		return errors.New("invalid phone number")
	}
	pincount:=u.countdigits(int(user.Address.Pincode))
	if pincount<6{
		return errors.New("invalid pin")
	}

	if match != "" && user.ConfirmPassword == user.Password && existingUser == nil {
		_, err := u.usercollection.InsertOne(u.ctx, user)
		return err
	}

	return errors.New("invalid input")
}


func (u *UserServiceImpl) getUserByEmail(email string) (*models.User, error) {
	filter := bson.M{"email": email}
	var user models.User

	err := u.usercollection.FindOne(u.ctx, filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
	
		return nil, nil
	} else if err != nil {
		
		return &user, err
	}

	return &user, nil
}
func (u*UserServiceImpl)countdigits(number int)(int64){
	count:=0
	for number >0{
       count++
	   number=number/10 
	}
	return int64(count)
}

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "name", Value: name}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}
func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(u.ctx)
	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}
	return users, nil
}
func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{Key: "name", Value: name}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
