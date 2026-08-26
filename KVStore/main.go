package main

import (
	"fmt"
	"time"
	"errors"

	"Store/store"
)

func main() {
	kv := store.New(500*time.Millisecond)
	defer kv.Close()
	kv.Set("user_1","harish",0)
	kv.Set("temp_token","KVZ89@#$O",1*time.Second)
	val,err := kv.Get("user_1")
	if err == nil {
		fmt.Printf("Fetched User 1 %v\n", val)
	}
	time.Sleep(2000 * time.Millisecond)
	val,err = kv.Get("temp_token")
	if err == nil {
		fmt.Printf("Fetched temp_token %v\n", val)
	}
	fmt.Println("waiting 1.2 seconds for temp_token to expire")
	
	_,err = kv.Get("temp_token")
	if errors.Is(err, store.ErrKeyExpired )  || errors.Is(err, store.ErrKeyNotFound) {
		fmt.Println("temp_token has expired or does not exist")
	}
	

}

