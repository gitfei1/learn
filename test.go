package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

var ip string

func initFlag() {
	flag.StringVar(&ip, "addr", "10.150.87.237:8901", "ip:port")
	flag.Parse()
}

func main() {
	initFlag()
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, uint32(74))
	binary.Write(pkg, binary.BigEndian, uint32(1))
	binary.Write(pkg, binary.BigEndian, uint32(1403))
	timer := int(time.Now().Month())*100000000 + time.Now().Day()*1000000 + time.Now().Hour()*10000 + time.Now().Minute()*100 + time.Now().Second()
	fmt.Println("xxxxxxxxx", timer)
	binary.Write(pkg, binary.BigEndian, uint32(timer))
	binary.Write(pkg, binary.BigEndian, uint32(1))

	code := "21403"
	lenc := len([]byte(code))
	codeS := make([]byte, 0)
	codeS = append(codeS, []byte(code)...)
	for i := 0; i < 21-lenc; i++ {
		codeS = append(codeS, 0x00)
	}

	name := "swgdjitsm"
	lenn := len([]byte(name))
	codeN := make([]byte, 0)
	codeN = append(codeN, []byte(name)...)
	for i := 0; i < 16-lenn; i++ {
		codeN = append(codeN, 0x00)
	}

	password := "gddw@1234"
	value := base64.StdEncoding.EncodeToString([]byte(password))
	lenp := len([]byte(value))
	codep := make([]byte, 0)
	codep = append(codep, []byte(value)...)
	for i := 0; i < 16-lenp; i++ {
		codep = append(codep, 0x00)
	}

	binary.Write(pkg, binary.BigEndian, codeS)
	binary.Write(pkg, binary.BigEndian, int8(1))
	binary.Write(pkg, binary.BigEndian, codeN)
	binary.Write(pkg, binary.BigEndian, codep)

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("conn err", err)
		return
	}
	n, err := conn.Write(pkg.Bytes())
	if err != nil {
		fmt.Println("write err", err)
		return
	}
	fmt.Println("write:", pkg.Bytes(), n)

	buf := make([]byte, 1024)
	rn, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read err ", err)
		return
	}
	fmt.Println("read:", buf[:rn])

	//发送短信
	destination := "13219059343"
	content := "123"
	pkg2 := new(bytes.Buffer)
	binary.Write(pkg2, binary.BigEndian, uint32(76+len(content)))
	binary.Write(pkg2, binary.BigEndian, uint32(1))
	binary.Write(pkg2, binary.BigEndian, uint32(1403))
	timer = int(time.Now().Month())*100000000 + time.Now().Day()*1000000 + time.Now().Hour()*10000 + time.Now().Minute()*100 + time.Now().Second()
	fmt.Println("xxxxxxxxx", timer)
	binary.Write(pkg2, binary.BigEndian, uint32(timer))
	binary.Write(pkg2, binary.BigEndian, uint32(1))

	MsgCode := "33"
	lenM := len([]byte(MsgCode))
	codeM := make([]byte, 0)
	codeM = append(codeM, []byte(MsgCode)...)
	for i := 0; i < 10-lenM; i++ {
		codeM = append(codeM, 0x00)
	}

	binary.Write(pkg2, binary.BigEndian, codeS)
	binary.Write(pkg2, binary.BigEndian, codeM)
	binary.Write(pkg2, binary.BigEndian, uint32(8901))
	binary.Write(pkg2, binary.BigEndian, int8(0))
	binary.Write(pkg2, binary.BigEndian, int8(1))
	binary.Write(pkg2, binary.BigEndian, []byte(destination))
	binary.Write(pkg2, binary.BigEndian, int8(15))
	binary.Write(pkg2, binary.BigEndian, int8(0))
	binary.Write(pkg2, binary.BigEndian, uint32(3))
	binary.Write(pkg2, binary.BigEndian, []byte(content))

	ns, err := conn.Write(pkg2.Bytes())
	if err != nil {
		fmt.Println("write err", err)
		return
	}
	fmt.Println("write submit:", pkg2.Bytes(), ns)

	bufs := make([]byte, 1024)
	rns, err := conn.Read(bufs)
	if err != nil {
		fmt.Println("read err ", err)
		return
	}
	fmt.Println("read submit:", bufs[:rns])

}
