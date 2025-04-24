package server

import ( "net" 
		"bufio"  
		"strings" 
		"realDB/internal/db" 
		// "realDB/cmd/main" 
		"github.com/sirupsen/logrus"

)

func HandleConnection(conn net.Conn){
	defer conn.Close() // defer lines are executed after everything else is executed in the surrounding fucntion
	reader := bufio.NewReader(conn)  // buffer reader is required as we are reading from the network , so io system call is optimized for buffers

	for {  
		line , err := reader.ReadString('\n') 
		if err != nil {
			db.RemoveConnFromWatchers(conn) 
			logrus.Info("Client disconnected")
			return
		}

		args := strings.Fields(strings.TrimSpace(line))
		if len(args) == 0 {
			continue
		}

		cmd := strings.ToUpper(args[0]) 

		switch cmd {
		// case <-  main.ShutdownChan:
		// 	conn.Write([]byte("SERVER_SHUTDOWN: Real-DB is shutting down.\n"))
		// 	conn.Close()
		// 	return
		case "SET":
			if len(args) != 3 {
				conn.Write([]byte("Usage : SET k v \n")) // []byte means we are writing to the connection in small byte level rather than the higher string level	
				continue
			}
			db.HandleSet(conn , args[1]  , args[2])
		
		case "GET":
			if len(args) != 2 {
				conn.Write([]byte("Usage : GET key \n")) 
				continue
			}
			db.HandleGet(conn, args[1])
		case "WATCH" :
			if len(args) != 2 {
				conn.Write([]byte("Usage : WATCH key \n"))
				continue
			}
			db.HandleWatch(conn, args[1])  
		case "DEL" : 
			if len(args) != 2 {
				conn.Write([]byte("Usage : DEL key \n")) 
				continue
			} 
			db.HandleDelete(conn , args[1])
		default : 
			conn.Write([]byte("Unknow command! \n"))
			conn.Write([]byte("real-db> "))
		}
	}
}	 

  