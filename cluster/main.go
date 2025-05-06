package main

import (
	"fmt"  
	"realDB/cluster/main-node"
)
 
func main() { 
	cl, err := client.NewClient("localhost:6369") 
	if err != nil {
		fmt.Printf("connection failed : %v", err) 
	} 
	defer cl.Close() 

	// sample code to send commands
	resp ,err := cl.SendCommand("SET tash win ") 
	if err != nil {
		fmt.Print(err)
	}  
	fmt.Println("SET Response was : ", resp)  

	// there will be serveral db's connected with different connections like cl1, cl2 ,cl3 and so on ... 
	// Consisted hashing needs to be implemented to choose which db to execute the command.



 

}  


