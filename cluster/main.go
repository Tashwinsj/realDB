package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"realDB/cluster/main-node"
	"realDB/internal/db"
	"realDB/cluster/consistentHashing"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)
  
 

type Node struct {
	Name string `yaml:"name"` 
	IP   string `yaml:"ip"` 
	Port int 	`yaml:"port"` 
}
type Config struct {
	Nodes []Node `yaml:"nodes"` 
} 

 
func main() {  
    ln , err := net.Listen("tcp" , ":6333") 
	if err != nil {
		panic(err)
	} 
	fmt.Println("Server started on port 6369")   
	// logrus.Info("Server Started on port 6369")

	// go HandleShutDown(ln)

	data ,err := os.ReadFile("config.yaml") 
	if err != nil {
		log.Fatalf("failed to read config.yaml: %v", err) 
	} 


	var cfg Config 
	err = yaml.Unmarshal(data , &cfg) 
	if err != nil {
		log.Fatalf("failed to parse YAML: %v", err) 
	}

	connections := make(map[string]*client.Client) 

	for _, node := range cfg.Nodes {
		address := fmt.Sprintf("%s:%d", node.IP, node.Port)
		fmt.Printf("Connecting to %s at %s...\n", node.Name, address) 


		cl, err := client.NewClient(address)
		if err != nil {
			fmt.Printf("connection to %s failed: %v\n", node.Name, err)
			continue
		}

		connections[node.Name] = cl
		fmt.Printf("Connected to %s.\n", node.Name)
	} 




	if cl, ok := connections["node1"]; ok {
		resp, err := cl.SendCommand("SET tash winner")
		if err != nil {
			fmt.Printf("error sending command to node1: %v\n", err)
		} else {
			fmt.Printf("[node1] SET response: %s\n", resp)
		}
	} 

	nodes := []string{"node1", "node2" }
	replicas := 100
	hr := chash.NewHashRing(nodes, replicas)
		
	// var keys []string

	for {
		conn , err := ln.Accept()
		if err != nil {
				select {
					default:
						// Unexpected error — log it
						// logrus.Errorf("Connection error : %v", err)
						continue
				}
		} 
		conn.Write([]byte( 
		`
╔╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗
╠╬╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╬╣
╠╣                                                                ╠╣
╠╣                                                                ╠╣
╠╣ ██████╗ ███████╗ █████╗ ██╗     ██████╗       ██████╗ ██████╗  ╠╣
╠╣ ██╔══██╗██╔════╝██╔══██╗██║     ██╔══██╗      ██╔══██╗██╔══██╗ ╠╣
╠╣ ██████╔╝█████╗  ███████║██║     ██║  ██║█████╗██║  ██║██████╔╝ ╠╣
╠╣ ██╔══██╗██╔══╝  ██╔══██║██║     ██║  ██║╚════╝██║  ██║██╔══██╗ ╠╣
╠╣ ██║  ██║███████╗██║  ██║███████╗██████╔╝      ██████╔╝██████╔╝ ╠╣
╠╣ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═════╝       ╚═════╝ ╚═════╝  ╠╣
╠╣                                                                ╠╣
╠╣  ██████╗██╗     ██╗   ██╗███████╗████████╗███████╗██████╗      ╠╣
╠╣ ██╔════╝██║     ██║   ██║██╔════╝╚══██╔══╝██╔════╝██╔══██╗     ╠╣
╠╣ ██║     ██║     ██║   ██║███████╗   ██║   █████╗  ██████╔╝     ╠╣
╠╣ ██║     ██║     ██║   ██║╚════██║   ██║   ██╔══╝  ██╔══██╗     ╠╣
╠╣ ╚██████╗███████╗╚██████╔╝███████║   ██║   ███████╗██║  ██║     ╠╣
╠╣  ╚═════╝╚══════╝ ╚═════╝ ╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝     ╠╣
╠╣                                                                ╠╣
╠╣                                                                ╠╣
╠╬╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝` + "\n")) 
		conn.Write([]byte("real-db Cluster> "))  

		 

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
			// node := hr.Get(args[1]) 


			if cl, ok := connections[hr.GetNode(args[1])]; ok { 
					var cmdd string 
					if len(args) == 3 {
						cmdd = fmt.Sprintf("%s %s %s", cmd , args[1] , args[2]) 
					} else {
						cmdd = fmt.Sprintf("%s %s ", cmd , args[1] ) 
					}
						resp, err := cl.SendCommand(cmdd) 
						conn.Write([]byte(resp))  
						if err != nil {
							fmt.Printf("error sending command to node: %v\n", err)
						} else {
							fmt.Printf("%s SET response: %s\n", hr.GetNode(args[1]) ,resp)
						}
				} 

			conn.Write([]byte("real-db Cluster> "))
		
			} 
		} 


		// go server.HandleConnection(conn)      // This infinite loop keeps looking for connections and wherever they come a seperate go routine is 
											  // spawned to handle the request ( seperate go routine for every client connected )
	
  

  
 


	 

	// // Open the file
	// file, err := os.Open("output.txt")
	// if err != nil {
	// 	fmt.Println("Error opening file:", err)
	// 	return
	// }
	// defer file.Close()

	// // Read the file line by line
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	// Take substring until first space
	// 	parts := strings.SplitN(line, " ", 2)
	// 	key := parts[0]
	// 	keys = append(keys, key)
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// // Example: assigning keys to nodes
	// for _, key := range keys {
	// 	node := hr.GetNode(key) // assuming hr is defined earlier
	// 	fmt.Printf("Key %q is assigned to node %q\n", key, node)
	// }

	// for _, key := range keys {
	// 	// fmt.Printf("Index %d: Key %s\n", i, key) 

	// 	if cl, ok := connections[hr.GetNode(key)]; ok {
	// 		cmd := fmt.Sprintf("SET %s winn", key)
	// 		resp, err := cl.SendCommand(cmd)
	// 		if err != nil {
	// 			fmt.Printf("error sending command to node: %v\n", err)
	// 		} else {
	// 			fmt.Printf("%s SET response: %s\n", hr.GetNode(key) ,resp)
	// 		}
	// 	} 
	// } 


	// counts := make(map[string]int)

	// for _, key := range keys {
	// 	node := hr.GetNode(key)
	// 	counts[node]++
	// }

	// for node, count := range counts {
	// 	fmt.Printf("Node %s got %d keys\n", node, count)
	// }

	

	// // Close all connections before exit
	// for name, cl := range connections {
	// 	cl.Close()
	// 	fmt.Printf("Closed connection to %s\n", name)
	// }

	// fmt.Println("Done.")

}  


