#4bytes calculator

calculate 4bytes value

Usage:  
cal4bytes [flags]  

Flags:  
-h, --help                help for cal4bytes  
--regular string      regular, example: buy_%d(bytes32,uint256)    
--target string       target 4bytes value, example: 0x00000000  
--thread uint         thread number (default 5)  
--thread.range uint   number of computations per thread (default 1000000000)  




./cal4bytes --regular "buy_%d(bytes32,uint256)" --target "0x00000000" --thread 5 --thread.range 1000000000