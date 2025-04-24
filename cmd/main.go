package main 

import ( 
	"fmt"  
	"net" 
	"realDB/internal/server"
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
		conn.Write([]byte( 
		`
.-----------------------------------------------------------.
|                                                           |
|     ██████╗ ███████╗ █████╗ ██╗      ██████╗ ██████╗      |
|     ██╔══██╗██╔════╝██╔══██╗██║      ██╔══██╗██╔══██╗     |
|     ██████╔╝█████╗  ███████║██║█████╗██║  ██║██████╔╝     |
|     ██╔══██╗██╔══╝  ██╔══██║██║╚════╝██║  ██║██╔══██╗     |
|     ██║  ██║███████╗██║  ██║███████╗ ██████╔╝██████╔╝     |
|     ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝      |
|                                                           |
|                  Connected to real-db>                    |
'-----------------------------------------------------------'` + "\n")) 
		conn.Write([]byte("real-db> "))

		go server.HandleConnection(conn)      // This infinite loop keeps looking for connections and wherever they come a seperate go routine is 
										// created to handle the request ( seperate go routine for every client connected )
	}
  
}
 
