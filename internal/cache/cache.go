package cache

import (
	"container/list"
	"fmt"
	"sync"
	"github.com/sirupsen/logrus"
) 

type entry struct {
	Key string 
	Val string
} 

type LRUCache struct {
	Capacity 	int64  // this is the max limit in bytes of the datastore
	currentSize int64  // this variable keeps track of the current size of datastore in bytes
	ll 		 *list.List  // this is a pointer to a list used for implementing LRU cache
	cache 	 map[string]*list.Element  // the actual cache is a map datastore that takes in key:string and gives out a pointer to an element in the linked list ( the linked list contains elements of type entry struct )
	mu 		 sync.Mutex  // mutex to ensure locking and prevent conncurrent access.
} 

func NewLRUCache(Capacity int64) (* LRUCache) { // func to initialize an new LRU cache
	return &LRUCache{
		Capacity: Capacity, 
		ll 		: list.New() , 
		cache 	: make(map[string]*list.Element) ,
	}
} 


func( c *LRUCache)  Get(key string) (string , bool) { 
	c.mu.Lock()	
	defer	c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)  
		return ele.Value.(*entry).Val , true  // ele is a pointer , ele.Value gives the actual element in that pointer , (*entry) is the type of the node in the list and .Val gives the actual value in the list.
	} 
	
	return "" , false
} 
 
func(c *LRUCache) Set(key string , value string) {
	c.mu.Lock() 
	defer c.mu.Unlock()   
	newsize := int64(len(key) + len(value))

	if ele , ok := c.cache[key] ; ok {
		c.ll.MoveToFront(ele)   
		oldpair := ele.Value.(*entry)
		c.currentSize -= int64(len(oldpair.Key) + len(oldpair.Val)) 
		oldpair.Val = value 
		c.currentSize += newsize
		logrus.Infof("The DB memory is : %d" , c.currentSize)
 		return 
	} 

	ele := c.ll.PushFront(&entry{ Key : key , Val : value})  
	c.cache[key] = ele  
	c.currentSize += newsize
	logrus.Infof("The DB memory is : %d" , c.currentSize)
 
	if c.currentSize > c.Capacity{ // we loop until the DB memory comes below the capacity.
	for  c.currentSize > c.Capacity { 
		logrus.Info("DB memory limit exceeded! LRU Eviction in progress.") 
		fmt.Print("DB memory limit exceeded! LRU Eviction in progress. \n")
		lru := c.ll.Back()
		if lru == nil {
			break
		}
		evicted := lru.Value.(*entry)
		c.ll.Remove(lru)
		delete(c.cache, evicted.Key)
		c.currentSize -= int64(len(evicted.Key) + len(evicted.Val))
	} 
	logrus.Infof("LRU Eviction process completed! DB memory : %d" , c.currentSize) 
	fmt.Printf("LRU Eviction process completed! DB memory : %d \n" , c.currentSize) 
	}
}  

func (c *LRUCache) Del (key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok { 
		entry := ele.Value.(*entry) 
		c.currentSize -= int64(len(entry.Key) + len(entry.Val))
		c.ll.Remove(ele) 
		delete(c.cache , key)
	} 
	logrus.Infof("The DB memory is : %d" , c.currentSize)
}