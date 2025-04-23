package main 

import ( 
	"bufio" 
	"fmt"  
	"net" 
	"strings" 
	"sync"  	 
) 

var (
	store 	 = make(map[string]string) // key value store 
	watchers = make(map[string][]net.Conn) 
	mu		 	sync.RWMutex // synchornous mutex allows only one goroutine to lock and unlock at once
)

func main(){
	ln , err := net.Listen("tcp" , ":6369") 
	if err != nil {
		panic(err)
	} 
	fmt.Println("Server started on port 6369") 

	for {
		conn , err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error : ", err) 
			continue
		}
		go handleConnection(conn)      // This infinite loop keeps looking for connections and wherever they come a seperate go routine is 
										// created to handle the request ( seperate go routine for every client connected )
	}


} 

func handleConnection(conn net.Conn){
	defer conn.Close() // defer lines are executed after everything else is executed in the surrounding fucntion
	reader := bufio.NewReader(conn)  // buffer reader is required as we are reading from the network , so io system call is optimized for buffers

	for {
		line , err := reader.ReadString('\n') 
		if err != nil {
			removeConnFromWatchers(conn) 
			fmt.Println("Client disconnected")
			return
		}

		args := strings.Fields(strings.TrimSpace(line))
		if len(args) == 0 {
			continue
		}

		cmd := strings.ToUpper(args[0]) 

		switch cmd {
		case "SET":
			if len(args) != 3 {
				conn.Write([]byte("Usage : SET k v \n")) // []byte means we are writing to the connection in small byte level rather than the higher string level	
				continue
			}
			handleSet(conn , args[1]  , args[2])
		
		case "GET":
			if len(args) != 2 {
				conn.Write([]byte("Usage : GET key \n")) 
				continue
			}
			handleGet(conn, args[1])
		case "WATCH" :
			if len(args) != 2 {
				conn.Write([]byte("Usage : WATCH key \n"))
				continue
			}
			handleWatch(conn, args[1]) 
		default : 
			conn.Write([]byte("Unknow command! \n"))
		}
	}
}	 

func handleSet(conn net.Conn , key string , val string ) {
	mu.Lock() 
	store[key] = val 
	conns := watchers[key]
	mu.Unlock()  

	for _ , watcher := range conns {
		if watcher != conn { // the setting function is sending out exept to itself !
			watcher.Write([]byte(fmt.Sprintf("WATCH : key '%s' was update to '%s' \n" , key ,val ))) 
		}
	}
} 

func handleGet( conn net.Conn , key string ) {
	mu.RLock()
	val , ok := store[key] 
	mu.RUnlock()

	if ok {
		conn.Write([]byte(fmt.Sprintf("%s \n", val))) 
	} else  {
		conn.Write([]byte("nil \n"))
	}
} 

func handleWatch(conn net.Conn , key string) {
	mu.Lock()
	watchers[key] = append(watchers[key] , conn) 
	mu.Unlock()	

	conn.Write([]byte("WATCHING " + key + "\n" )) 
} 


func removeConnFromWatchers(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for key , conns := range watchers {
		newConns := make([]net.Conn , 0 , len(conns)) 
		for _ , c := range conns {
			if c != conn {
				newConns = append(newConns, c) 
			} 
		} 
		if len(newConns) == 0 {
			delete(watchers, key) 
		} else { 
			watchers[key] = newConns
		}
	}
}