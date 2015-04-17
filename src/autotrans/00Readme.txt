# Auto transition sample

## How to build

(0) Install go language. 
    Tested in linux/amd64. Go version 1.4. 
    I hope Mac/Windows works, too.

(1) Build

  go build *.go

## How to run

(0) Find LiveWedge's IP address by iPad
(1) Run the command as below

  ./autotrans IP_address_of_LiveWedge

  example)
  ./autotrans 172.16.130.244

(2) Open WebUI by any web browser

  http://localhost:8080/

Note:
* The file name of A still picuture to upload is fixed to 'a.jpg'.
* Settings are saved in 'autotrans.json'.
