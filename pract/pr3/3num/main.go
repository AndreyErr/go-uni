package main

import (
	"fmt"
	"math/rand"

	"github.com/reactivex/rxgo/v2"
)

type UserFriend struct {
	userId int
	friendId int
}

func main(){
	var userFriendArr []UserFriend
	for i :=0; i <= 10; i++{
		userFriend:= UserFriend{
			userId: rand.Intn(10),
			friendId: rand.Intn(10),
		}
		userFriendArr = append(userFriendArr, userFriend)
	}

	getFriends := func(userId2 int) rxgo.Observable {
		return rxgo.
			Just(userFriendArr)().
			Filter(func(i interface{})bool {
				userFriendId := i.(UserFriend)
				return userFriendId.userId == userId2
			})
	}

	for _, uf := range userFriendArr {
        fmt.Printf("User ID: %d, Friend ID: %d\n", uf.userId, uf.friendId)
    }

	randomUserId := make([]int, 2)
	for i:=0; i < 2; i++{
		randomUserId[i] = rand.Intn(10)
	}
	fmt.Printf("\n Рандомные UserId: %d \n\n", randomUserId)

	observable := rxgo.Just(randomUserId)().
	FlatMap(func(i rxgo.Item) rxgo.Observable {
		v := i.V.(int)
		return(getFriends(v))
	})

	for v:= range observable.Observe(){
		fmt.Printf("Нахождение информации в массиве UserFriend о рандомных UserId %d \n",v.V)
	}
}

