package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"skywatch24/mypack"

	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	// @todo json如果檔案太大, 整檔讀取會有效能問題

	args := os.Args

	var filename = "movies.json"

	if len(args) > 1 {
		filename = args[1]
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 為了支援任意json, 使用泛型讀取json
	var jsonData interface{}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatalf("json Error: %v, input file is not json file.\n", err)
		// fmt.Printf("json Error: %v, input file is not json file.\n", err)
		// jsonData = data
		// 如果一個檔案裡面只有一個數字(例如test_int), 是會被json.Unmarshal當作是json檔案的, 而且該數字會統一被轉成float64
		// 如果是文字(test_str32)或二進位檔(test_binary), 就不會執行, 除非把上面兩行uncomment, 然後comment掉log.Fatalf那行
	}

	// 編碼
	// 1. 用自己的mypack編碼
	// 2. 用msgpack編碼
	// 3. 印出來看看

	encodedMypack, err := mypack.Marshal(jsonData)
	if err != nil {
		log.Fatalf("mypack Error: %v", err)
	}

	encodedMsgpack, err := msgpack.Marshal(jsonData)
	if err != nil {
		log.Fatalf("msgpack Error: %v", err)
	}

	fmt.Printf("Encoded(mypack): %x\n", encodedMypack)
	fmt.Printf("Encoded(msgpack): %x\n", encodedMsgpack)

	// 解碼
	// 1. 用自己的mypack解碼
	// 2. 用msgpack解碼
	// 3. 印出來看看

	decodedMypack, err := mypack.Unmarshal(encodedMypack)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var decodedMsgpack interface{}
	err = msgpack.Unmarshal(encodedMsgpack, &decodedMsgpack)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Decoded(mypack): %+v\n", decodedMypack)
	fmt.Printf("Decoded(msgpack): %+v\n", decodedMsgpack)

	// @todo
	// 雖然編碼後可以解碼, 但是跟msgpack還是有不一樣的地方
	// 觀察 movies2.json的output, 好像是map資料型態的資料, 印出的順序不固定, 如果測試檔案沒有map的話, output就是一樣的
}
