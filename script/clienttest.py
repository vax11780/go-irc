#!/usr/bin/env python

import os
import pexpect
import socket
import sys
import time

def waitServer(procIn, waitData):
    try:
        index = 1
        while (index != 0):
            index = procIn.expect([waitData, '>you'], timeout=5)
        
    except pexpect.EOF:
        print "FAILURE: Server exited too soon"
        sys.exit(1)
    except pexpect.TIMEOUT:
        print "FAILURE: Server failed to send expected data"
        sys.exit(2)
    except:
        print "FAILURE: Unexpected failure"
        sys.exit(3)

myDir = os.path.dirname(sys.argv[0])
while myDir.find('script') != -1:
    myDir = os.path.dirname(myDir);
   
myDir = os.path.join(myDir, "bin")
serverCmd = os.path.join(myDir, "server")
clientCmd = os.path.join(myDir, "client")

if len(sys.argv) == 1:
    clientCmd += r' -s 127.0.0.1 -p 9999'
    print 'Start server with: ' + serverCmd
    server = pexpect.spawn(serverCmd)
    time.sleep(1)
elif len(sys.argv) == 2:
    ipaddr = socket.gethostbyname(sys.argv[1]) 
    clientCmd += (' -s ' + ipaddr + ' -p 9999')
else:
    print "FAILURE: Too many paramaters - 'python clienttest.py ipaddr'"
    sys.exit(6)

print 'Testing with: ' + clientCmd

fout = file('client1.txt','w')
proc = pexpect.spawn(clientCmd)
proc.logfile = fout
waitServer(proc, 'Please give you name:')
proc.sendline('Client1')
waitServer(proc, 'you>')

fout2 = file('client2.txt','w')
proc2 = pexpect.spawn(clientCmd)
proc2.logfile = fout2
waitServer(proc2, 'Please give you name:')
proc2.sendline('Client2')
waitServer(proc2, 'you>')

proc.sendline('My Message1')
proc2.sendline('My Message2')

waitServer(proc, 'you>')
waitServer(proc2, 'you>')
proc.sendline('/quit') # for myclient to quit
proc2.sendline('/quit') # for myclient to quit

proc.expect([pexpect.EOF])
proc2.expect([pexpect.EOF])

fout.close()
fout2.close()

clientData = ""
clientData = open('client1.txt').read()
if clientData.find(r'> My Message2') == -1:
    print "FAILURE: Client2 message not received"
    sys.exit(4)

clientData = ""
clientData = open('client2.txt').read()
if clientData.find(r'> My Message1') == -1:
    print "FAILURE: Client1 message not received"
    sys.exit(4)

os.remove('client1.txt')
os.remove('client2.txt')

sys.exit(0)

