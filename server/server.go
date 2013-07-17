package main

import (
    "fmt";
    "net";
    "log";
    "os";
    "container/list";
    "bytes";
    "flag";
)

// flag for debuging info. or a simple log
var debug = flag.Bool("d", false, "set the debug modus( print informations )")

type ClientChat struct {
    Name string;        // name of user
    IN chan string;     // input channel for to send to user
    OUT chan string;    // input channel from user to all
    Con net.Conn;      // connection of client
    Quit chan bool;     // quit channel for all goroutines
    ListChain *list.List;    // reference to list
}

// read from connection and return true if ok
func (c ClientChat) Read(buf []byte) bool{
    _, err := c.Con.Read(buf);
    if err!=nil {
        c.Close();
        return false;
    }
    Log("Read():  ", " bytes");
    return true;
}

// close the connection and send quit to sender
func (c *ClientChat) Close() {
    c.Quit<-true;
    c.Con.Close();
    c.deleteFromList();
}

// compare two clients: name and network connection
func (c *ClientChat) Equal(cl *ClientChat) bool {
    if bytes.EqualFold([]byte(c.Name), []byte(cl.Name)) {
        if c.Con == cl.Con {
            return true;
        }
    }
    return false;
}

// delete the client from list
func (c *ClientChat) deleteFromList() {
    for e := c.ListChain.Front(); e != nil; e = e.Next() {
        client := e.Value.(ClientChat);
        if c.Equal(&client) {
            Log("deleteFromList(): ", c.Name);
            c.ListChain.Remove(e);
        }
    }
}

// func Log(v ...): logging. give log information if debug is true
func Log(v ...string) {
    if *debug == true {
        ret := fmt.Sprint(v);
        log.Printf("SERVER: %s", ret);
    }
}

// func test(): testing for error
func test(err error, mesg string) {
    if err!=nil {
        log.Printf("SERVER: ERROR: ", mesg);
         os.Exit(-1);
    } else {
        Log("Ok: ", mesg);
    }
}

// handlingINOUT(): handle inputs from client, and send it to all other client via channels.
//type ClientChat
func handlingINOUT(IN <-chan string, lst *list.List) {
    for {
        Log("handlingINOUT(): wait for input");
        input := <-IN;  // input, get from client
        // send to all client back
        Log("handlingINOUT(): handling input: ", input);
        for value := lst.Front(); value != nil; value = value.Next() {
            client := value.Value.(ClientChat)
            Log("handlingINOUT(): send to client: ", client.Name);
            client.IN<- input;
        }  
    }
}



// go routine spun up by clientHandling to manage incoming
// client traffic and terminate upon receipt of /quit command from client
// prints out client string data and thens forwards that data
// to the handlingINOUT routine via a channel notifier 
func clientreceiver(client *ClientChat) {
    buf := make([]byte, 2048);

    Log("clientreceiver(): start for: ", client.Name);
    for client.Read(buf) {
        
        if bytes.EqualFold(buf, []byte("/quit")) {
            client.Close();
            break;
        }
        Log("clientreceiver(): received from ",client.Name, " (", string(buf), ")");
        send := client.Name+"> "+string(buf);
        client.OUT<- send;
        for i:=0; i<2048;i++ {
            buf[i]=0x00;
        }
    }    

    client.OUT <- client.Name+" has left chat";
    Log("clientreceiver(): stop for: ", client.Name);
}

// clientsender(): get the data from handlingINOUT via channel (or quit signal from
// clientreceiver) and send it via network
func clientsender(client *ClientChat) {
    Log("clientsender(): start for: ", client.Name);
    for {
        Log("clientsender(): wait for input to send");
        select {
            case buf := <- client.IN:
                Log("clientsender(): send to \"", client.Name, "\": ", string(buf));
                client.Con.Write([]byte(buf));
            case <-client.Quit:
                Log("clientsender(): client want to quit");
                client.Con.Close();
                break;
        }
    }
    Log("clientsender(): stop for: ", client.Name);
}

// clientHandling(): get the username and create the clientsturct
// start the clientsender/receiver, add client to list.
func clientHandling(con net.Conn, ch chan string, lst *list.List) {
    buf := make([]byte, 1024);
    con.Read(buf);
    name := string(buf);
    newclient := &ClientChat{name, make(chan string), ch, con, make(chan bool), lst};

    Log("clientHandling(): for ", name);
    go clientsender(newclient);
    go clientreceiver(newclient);
    lst.PushBack(*newclient);
    ch<- name+" has joinet the chat";
    ch<- "l33t users only";
}

func main() {
    flag.Parse();
    Log("main(): start");

    // create the list of clients
    clientlist := list.New();
    in := make(chan string);
    Log("main(): start handlingINOUT()");
    go handlingINOUT(in, clientlist);
    
    // create the connection
    netlisten, err := net.Listen("tcp", "0.0.0.0:9999");
    test(err, "main Listen");
    defer netlisten.Close();

    for {
        // wait for clients
        Log("main(): wait for client ...");
        conn, err := netlisten.Accept();
        test(err, "main: Accept for client");
        go clientHandling(conn, in, clientlist);
    }
}
