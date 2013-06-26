package main

import (
    "fmt";
    "net";
    "log";
    "os";
    "bytes";
    "bufio";
    "time";
    "flag";
)

var running bool;  // global variable if client is running
var debug = flag.Bool("d", false, "enable debug mode ( display debugging information )")
var server = flag.String("s", "", "server name to connect to")
var port = flag.String("p", "" , "port number to connect to")

// func Log(v ...): loging. give log information if debug is true

func Log(v ...string) {
    if *debug == true {
        ret := fmt.Sprint(v);
        log.Printf("CLIENT: %s", ret);
    }
}

// func test(): testing for error
func test(err error, mesg string) {
    if err!=nil {
         log.Printf("CLIENT: ERROR: ", mesg);
         os.Exit(-1);
    } else {
        Log("Ok: ", mesg);
    }
}

// read from connection and return true if ok
func Read(con net.Conn) string{
    var buf [4048]byte;
    _, err := con.Read(buf[0:4048]);
    if err!=nil {
        con.Close();
        running=false;
        return "Error in reading!";
    }
    str := string(buf[0:4048]);
    fmt.Println();
    return string(str);
}

// clientsender(): read from stdin and send it via network
func clientsender(cn net.Conn) {
    reader := bufio.NewReader(os.Stdin);
    for {
        fmt.Print("you> ");
        input, _ := reader.ReadBytes('\n');
        if bytes.EqualFold(input, []byte("/quit\n")) {
            cn.Write([]byte("/quit"));
            running = false;
            break;
        }
        Log("clientsender(): send: ", string(input[0:len(input)-1]));
        cn.Write(input[0:len(input)-1]);
    }
}

// clientreceiver(): wait for input from network and print it out
func clientreceiver(cn net.Conn) {
    for running {
        fmt.Println(Read(cn));
        fmt.Print("you> ");
    }
}

func usage() {
    fmt.Fprintf(os.Stderr, "usage: client 192.168.1.1:9999\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
   flag.Usage = usage;
   flag.Parse();
    
    running = true;
    Log("main(): start ");
    
    destination := fmt.Sprintf("%s:%s", *server,*port); 
    fmt.Println("Connected to: ", destination);

    Log("main(): connecto to ", destination);
    cn, err := net.Dial("tcp", destination);
    test(err, "dialing");
    defer cn.Close();
    Log("main(): connected ");

    // get the user name
    fmt.Print("Please give you name: ");
    reader := bufio.NewReader(os.Stdin);
    name, _ := reader.ReadBytes('\n');

    //cn.Write(strings.Bytes("User: "));
    cn.Write(name[0:len(name)-1]);

    // start receiver and sender
    Log("main(): start receiver");
    go clientreceiver(cn);
    Log("main(): start sender");
    go clientsender(cn);
    
    // wait for quiting (/quit). run until running is true
    for ;running; {
        time.Sleep(1*1e9);
    }
    Log("main(): stoped");
}
