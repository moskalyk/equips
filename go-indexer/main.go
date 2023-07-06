package main

import (
    "os"
    "context"
	"fmt"
    "strings"
	"log"
    "net/http"
    "unicode/utf8"
    "encoding/hex"
    "github.com/gin-gonic/gin"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
    "database/sql"
	"github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type Event struct {
	Owner        string
	TokenAddress string
	TokenID      int64
    Index        int64
    Salt         int64
}

func dropTable(db *sql.DB, tableName string) error {
    query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
    _, err := db.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

func getItemsByParameter(parameter string) ([]Event, error) {
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

    if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT owner, token_address, token_id, salt FROM equips WHERE owner = $1", parameter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Event

	for rows.Next() {
		var item Event
        err := rows.Scan(&item.Owner, &item.TokenAddress, &item.TokenID, &item.Salt, &item.Index)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func getEventLogs() ([]Event, error) {
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

    if err != nil {
        return nil, err
    }
    
    defer db.Close()

    rows, err := db.Query("SELECT owner, token_address, token_id, salt, index FROM equips;")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []Event

    for rows.Next() {
        var item Event
        err := rows.Scan(&item.Owner, &item.TokenAddress, &item.TokenID, &item.Salt, &item.Index)
        if err != nil {
            return nil, err
        }

        fmt.Printf("Owner: %s, TokenAddress: %s, TokenID: %d\n", item.Owner, item.TokenAddress, item.TokenID)
        
        items = append(items, item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return items, nil
}

func handleAddItem(c *gin.Context) {
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

    if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    for i := 0; i < 1500; i++ {
            insertSQL := `
            INSERT INTO equips (owner, token_address, token_id)
            VALUES ($1, $2, $3)
        `
        // Execute the insert statement
        _, err = db.Exec(insertSQL, "0xbabe", "0xdeaf", i)
        if err != nil {
            log.Fatal(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        fmt.Println("Table added to")
	}
    c.JSON(http.StatusOK, gin.H{"status": "Table added to"})
}

func handleInitTable(c *gin.Context) {
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
    // Create a table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS equips (
        id SERIAL PRIMARY KEY,
        owner VARCHAR(42),
        token_address VARCHAR(42),
        token_id INT,
        index INT,
        salt INT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return 
	}

	fmt.Println("Table created successfully")
    c.JSON(http.StatusOK, gin.H{"status": "Table created successfully"})
}

func handleParamLookUp(c *gin.Context){
    address := c.Param("address")
    items, err := getItemsByParameter(address)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"items": items})
}

func handleEventLogs(c *gin.Context) {
    logs, err := getEventLogs()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func handleDropTable(c *gin.Context){
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

    // Drop the table
    err = dropTable(db, "equips")
    if err != nil {
        fmt.Println("Error dropping table:", err)
        return
    }

    fmt.Println("Table dropped successfully!")
}


func removeLeadingZeros(bytes []byte) []byte {
	// Find the index where non-zero values start
	index := 0
	for i, b := range bytes {
		if b != 0 {
			index = i
			break
		}
	}

	// Return the slice from the non-zero index to the end
	return bytes[index:]
}

func bytes32ToStringWithAddress(bytes32 []byte) (string, string, string, string) {
	combinedBytes := removeLeadingZeros(bytes32)

	addrBytes := combinedBytes[:20]
	addressConverted := getAddressFromBytes(addrBytes)

	var strBuilder strings.Builder
	for i := 20; i < len(combinedBytes); i++ {
		byteVal := combinedBytes[i]
		if byteVal == 0 {
			break
		}
		strBuilder.WriteByte(byte(byteVal))
	}

	str := strBuilder.String()
    split := strings.Split(str, ":")

	return addressConverted, split[1], split[2], split[3]
}

func listen() {
    db, err := sql.Open("postgres", "postgres://"+os.Getenv("PG_USER")+":"+os.Getenv("PG_PASSWORD")+"@localhost/equips?sslmode=disable")

    client, err := ethclient.Dial("wss://polygon-mumbai.g.alchemy.com/v2/"+os.Getenv("ALCHEMY_RPC_WSS"))
    if err != nil {
        log.Println("err")
        log.Fatal(err)
    }

    eventName := `Equip`
    contractABI := `
    [
	{
		"inputs": [
			{
				"internalType": "address[]",
				"name": "token_addresses",
				"type": "address[]"
			},
			{
				"internalType": "uint256[]",
				"name": "token_ids",
				"type": "uint256[]"
			},
			{
				"internalType": "bytes32[]",
				"name": "equips",
				"type": "bytes32[]"
			},
			{
				"internalType": "uint256",
				"name": "salt",
				"type": "uint256"
			}
		],
		"name": "equip",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "owner",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "payload",
				"type": "bytes32"
			}
		],
		"name": "Equip",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "addr",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "value",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "salt",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "index",
				"type": "uint256"
			}
		],
		"name": "concatenate",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "addr",
				"type": "address"
			},
			{
				"internalType": "string",
				"name": "str",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "salt",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "index",
				"type": "string"
			}
		],
		"name": "stringToBytes32WithAddress",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	}
]
	`
    parsedABI, err := abi.JSON(strings.NewReader(contractABI))
    contractAddress := common.HexToAddress("0x69d6aC536E56C931f0Dd67d13eE3bA23e5Baaa4c")
    query := ethereum.FilterQuery{
        Addresses: []common.Address{contractAddress},
        Topics:    [][]common.Hash{{parsedABI.Events[eventName].ID}},
    }

    logs := make(chan types.Log)
    sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
        log.Println("err")
        log.Fatal(err)
    }

    for {
        select {
            case err := <-sub.Err():
                log.Fatal(err)
            case eventLog := <-logs:
                
                // decodedStr := string(eventLog.Topics[2].Bytes())
	            // fmt.Println("Decoded String (UTF-8):", decodedStr)
                // log.Println(bytes32ToStringWithAddress(eventLog.Topics[2].Bytes()))
                address, token_id, salt, index := bytes32ToStringWithAddress(eventLog.Topics[2].Bytes())

                insertSQL := `
                    INSERT INTO equips (owner, token_address, token_id, salt, index)
                    VALUES ($1, $2, $3, $4, $5)
                `
                // Execute the insert statement
                _, err = db.Exec(
                    insertSQL, 
                    common.BytesToAddress(eventLog.Topics[1].Bytes()).Hex(),
                    address,
                    token_id,
                    salt,
                    index,
                )

                fmt.Println("Table added to")
        }
    }
}

func getAddressFromBytes(bytes []byte) string {
	// Remove the "0x" prefix from the bytes string
	bytesString := strings.TrimPrefix(hex.EncodeToString(bytes), "0x")
	// Pad the bytes string with leading zeros to ensure it has a length of 40 characters
	paddedBytesString := fmt.Sprintf("%0*s", 40, bytesString)
	// Add the "0x" prefix to the padded bytes string
	address := "0x" + paddedBytesString
	return address
}

func decodeUTF8(bytes []byte) string {
	runes := make([]rune, 0, len(bytes))

	for len(bytes) > 0 {
		r, size := utf8.DecodeRune(bytes)
		runes = append(runes, r)
		bytes = bytes[size:]
	}

	return string(runes)
}

func remove0xPrefix(bytes32 string) string {
	return strings.TrimPrefix(bytes32, "0x")
}

func parseBytes32String(bytesLike string) string {
	trimmed := strings.TrimPrefix(bytesLike, "0x")
	bytes := common.Hex2Bytes(trimmed)

	// Remove trailing null bytes
	for len(bytes) > 0 && bytes[len(bytes)-1] == 0 {
		bytes = bytes[:len(bytes)-1]
	}

	return string(bytes)
}

func main() {

    err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

    go func() {
        listen()
    }()

    r := gin.Default()

    // *~*~*~*~ for testing
    r.GET("/init", handleInitTable)
    r.GET("/add", handleAddItem)
    r.GET("/drop", handleDropTable)

    // *~*~*~*~ for api
    r.GET("/all", handleEventLogs)
    r.GET("/lookup/:address", handleParamLookUp)

    // @^@^@ listening
    r.Run(":7077")
}