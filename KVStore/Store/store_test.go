package store

import (
	"testing"
	"fmt"
	"time"
	"sync"
	"errors"
)

func TestStore_BasicOperation(t *testing.T) {
	s := New(0)
	defer s.Close()
	tests := []struct{
		name string
		setup func()
		key string
		expectedValue any
		expectedErr error
	}{
		{
			name:"Get existing string key",
			setup: func() {
				s.Set("greeting","hello",0)

			},
			key : "greeting",
			expectedValue: "hello",
			expectedErr: nil,
		},
		{
			name: "Get existing integer key",
			setup: func() {
				s.Set("count",42,0)

			},
			key: "count",
			expectedValue: 42,
			expectedErr: nil,
		},
		{
			name: "Get non-existing key",
			setup: func() {},
			key : "nonexistent",
			expectedValue: nil,
			expectedErr: ErrKeyNotFound,
		},
		{
			name: "Delete and verify removal",
			setup: func() {
				s.Set("to_delete","temporary",0 )
				_ =  s.Delete("to_delete")
			},
			key: "to_delete",
			expectedValue: nil,
			expectedErr: ErrKeyNotFound,
		} ,
	}
	for _,tc := range tests {
		t.Run(tc.name,func(t *testing.T){
			tc.setup()
			val,err := s.Get(tc.key)
			if !errors.Is(err,tc.expectedErr) {
				t.Fatalf("expected error %v, got %v\n",tc.expectedErr,err )
			}
			if val != tc.expectedValue {
				t.Fatalf("expected value %v but got %v",tc.expectedValue,val)
			}
		})
	}
	
}

func TestStore_TTLEXpirstion(t *testing.T) {
	s := New(50*time.Millisecond)
	defer s.Close()
	s.Set("short_lived","data",60*time.Millisecond)
	time.Sleep(60*time.Millisecond)
	_,err := s.Get("short_lived")
	if !errors.Is(err,ErrKeyExpired) {
		t.Fatalf("expected ErrKeyExpired on lazy eviction , got %v",err)	
	}
	s.Set("auto_sweep","data",30*time.Millisecond)
	time.Sleep(100*time.Millisecond)
	s.mu.RLock()
	_,exists := s.items["auto_sweep"]
	s.mu.RUnlock()
	if exists {
		t.Errorf("expected background cleaner to delete 'auto_sweep' from map, but is still exists")

	}

}

func TestSore_ConcurrentAccess(t *testing.T) {
	s := New(200*time.Millisecond)
	defer s.Close()
	var wg sync.WaitGroup
	numGorountine := 100
	for i:=0; i<numGorountine; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d",id%10)
			s.Set(key,id,50*time.Millisecond)
		}(i)
	}
	for i :=0; i<numGorountine; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d",id%10)
			_,_ = s.Get(key)
		}(i)
	}
	for i:=0; i<numGorountine; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d",id%10)
			_ = s.Delete(key)
		}(i)
	}
	wg.Wait()
}