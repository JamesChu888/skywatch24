說明: 

1. 開發環境 MAC(intel) + VSCODE + go version go1.23.4 darwin/amd64

2. 目錄結構

   - skywatch24/  
     - mypack/encode.go  
     - mypack/decode.go  
     - mypack/mypack.go
     - mypack/mypack_test.go  
   - main.go  
   - go.mod  
   - go.sum  
   - README.md  
   - skywatch24  
   - movies.json  // json sample 預設測試檔案  
   - movies2.json // json array  
   - test_int  // 只有一個int, json.Unmarshal會視為合法, 且將數字都轉為float64  
   - test_str32 // 不合法json, 要修改main.go才能執行, main.go中有註解  
   - test_binary // 不合法json, 要修改main.go才能執行, main.go中有註解  

3. 編譯方式,如果是Mac(Intel)可略過這步: go build 

    * 因為想對照msgpack的結果, 所以, 在程式中分別用自己寫的mypack跟msgpack編解碼, 然後印出看結果, 對照一下, 所以要安裝msgpack(請使用go get)
    
4. 執行方式: 

  在Mac(intel)環境下, 可以直接開terminal, 切到skywatch24目錄下, 直接執行 

  ./skywatch24 完整檔名 // 沒打檔名 = ./skywatch24 movies.json 

  ex:

  ./skywatch24 movies.json  
  ./skywatch24 movies2.json  
  ./skywatch24 test_int  
  ./skywatch24 test_str32  // 要改一下main.go  
  ./skywatch24 test_binary // 要改一下main.go  

  或是  

  go run main.go

5. Unit test: go test ./...

6. 改善空間
  * msgpack遇到map好像encode順序不固定, 可能要花時間看一下對方的實作研究一下  
  * 使用goroutine加速 array, map的encode/decode  
  * 使用新版泛型T  
  
