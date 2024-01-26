package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

var currentUser string

var (
	client   *mongo.Client
	dbName   = "messanger"
	colName  = "users"
	messages = make(map[string][]Message)
)

type Message struct {
	From    string    `bson:"from"`
	To      string    `bson:"to"`
	Message string    `bson:"message"`
	Time    time.Time `bson:"time"`
}
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func connectDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://siddharthpranay8:bZUsMRyyGT7icDl2@messanger.jo6epxd.mongodb.net/?retryWrites=true&w=majority")

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)

	}

	fmt.Println("                                 ")
	fmt.Println("                                 ")
	fmt.Println("                                 ")
	fmt.Println("      Welcome to Let's Go! 🚀     ")

	fmt.Println("                                 ")

}

func disconnectDB() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Disconnected from MongoDB.")
	} else {
		fmt.Println("Client is nil, cannot disconnect.")
	}
}

func createUser(username, password string) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user := User{Username: username, Password: string(hashedPassword)}

	collection := client.Database(dbName).Collection(colName)

	// Create a unique index on the 'username' field
	index := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		// Ignore the error if the index already exists
		if !strings.Contains(err.Error(), "duplicate key error") {
			log.Fatal(err)
		}
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		// Check for duplicate key error and handle accordingly
		if strings.Contains(err.Error(), "duplicate key error") {
			fmt.Println("       ❗ Username already exists choose a different username                  ")

		} else {
			log.Fatal(err)
		}
		return
	}

	currentUser = username

	fmt.Println("       ✅ Logged In as", currentUser)

}

func loginUser(username, password string) bool {
	collection := client.Database(dbName).Collection(colName)

	// filter := bson.M{"username": username, "password": password}
	// err := collection.FindOne(context.Background(), filter).Err()

	// if err != nil {
	// 	fmt.Println("       ❗ Login failed. Incorrect username or password  ")
	// 	return false
	// }

	// currentUser = username

	// fmt.Println("       ✅ Logged In as", currentUser)
	// return true
	var user User

	filter := bson.M{"username": username}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Println("       ❗ Login failed. Incorrect username or password  ")
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("       ❗ Login failed. Incorrect username or password  ")
		return false
	}

	currentUser = username
	fmt.Println("       ✅ Logged In as", currentUser)
	return true
}

func sendMessage(from, to, message string) {
	if !userExists(from) {
		fmt.Println("       ❗ Sender username does not exist ")
		return
	}

	if !userExists(to) {
		fmt.Println("       ❗ Recipient username does not exist ")

		return
	}

	msg := Message{
		From:    from,
		To:      to,
		Message: message,
		Time:    time.Now(),
	}

	collection := client.Database(dbName).Collection("messages")

	_, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("       ✅ Message sent to", to)

}

func viewMessages(username string) {
	collection := client.Database(dbName).Collection("messages")

	filter := bson.M{"$or": []bson.M{{"from": username}, {"to": username}}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			log.Fatal(err)
		}
		messages[username] = append(messages[username], msg)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("                                 ")
	fmt.Println("                                 ")
	fmt.Println("                                 ")

	fmt.Printf("        Hello, %s 👋\n", username)
	fmt.Println("                                 ")

	for _, msg := range messages[username] {
		printDecoratedMessage(msg)
	}
	messages[username] = nil

}

func printDecoratedMessage(msg Message) {
	// Format the time in a specific way with color
	timeStr := "\033[1;34m" + msg.Time.Format(" Jan 2, 2006 at 3:04pm ") + "\033[0m" // Bold Blue

	// Define colors and formatting
	colorFrom := "\033[1;32m" // Bold Green
	colorReset := "\033[0m"   // Reset to default
	// Determine the display name based on whether the message is from the current user
	var displayName string
	if msg.From == currentUser {
		displayName = "you"
	} else {
		displayName = msg.From
	}

	// Determine the recipient display name
	var recipientDisplayName string
	if msg.To == currentUser {
		recipientDisplayName = "you"
	} else {
		recipientDisplayName = msg.To
	}

	// Print the decorated message with sender and recipient information
	fmt.Printf("\n       ┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("       | %s%s %s ─ %s  %s: %s\n", timeStr, colorFrom, displayName, recipientDisplayName, colorReset, msg.Message)
	fmt.Printf("       └──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘\n\n")

}

func userExists(username string) bool {
	collection := client.Database(dbName).Collection(colName)

	filter := bson.M{"username": username}
	err := collection.FindOne(context.Background(), filter).Err()

	return err == nil // If err is nil, the user exists; otherwise, it does not
}

func readFullSentence() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {

	connectDB()
	defer disconnectDB()

	for {

		fmt.Println("                                 ")
		fmt.Println("                                 ")
		fmt.Println("                                 ")

		if currentUser != "" {
			fmt.Println("       🔍  Choose an option", currentUser)
		} else {
			fmt.Println("       🔍  Choose an option       ")
		}

		fmt.Println("                                 ")
		fmt.Println("                                 ")
		fmt.Println("       1️⃣   Sign Up                ")
		fmt.Println("                                 ")

		fmt.Println("       2️⃣   Send Message           ")
		fmt.Println("                                 ")

		fmt.Println("       3️⃣   View Messages          ")
		fmt.Println("                                 ")

		fmt.Println("       4️⃣   Logout                 ")
		fmt.Println("                                 ")

		fmt.Println("       5️⃣   Exit                   ")
		fmt.Println("                                 ")

		var choice int
		fmt.Println("                                 ")
		fmt.Print("       🗨️  Enter your choice : ")
		fmt.Scan(&choice)
		fmt.Println("                                 ")
		fmt.Println("                                 ")

		switch choice {
		case 1:
			var username, password string
			fmt.Print("       🔖 Username : ")
			fmt.Scan(&username)
			fmt.Println("                                 ")
			fmt.Print("       💀 Password : ")
			fmt.Scan(&password)
			fmt.Println("                                 ")
			fmt.Println("                                 ")

			createUser(username, password)

		case 2:
			if currentUser == "" {

				var username, password string
				fmt.Print("       🔖 Enter your Username : ")
				fmt.Scan(&username)
				fmt.Println("                                 ")

				fmt.Print("       💀 Enter your Password : ")

				fmt.Scan(&password)
				fmt.Println("                                 ")
				fmt.Println("                                 ")

				if loginUser(username, password) {
					currentUser = username
					var to, message string
					fmt.Println("")

					fmt.Println("")

					fmt.Print("       🔖 Enter recipient username : ")
					fmt.Scan(&to)
					fmt.Println("")
					fmt.Print("       ✉️  Enter message : ")
					message = readFullSentence()
					fmt.Println("")
					sendMessage(currentUser, to, message)
				}
			} else {
				var to, message string
				fmt.Println("")

				fmt.Println("")

				fmt.Print("       🔖 Enter recipient username : ")
				fmt.Scan(&to)
				fmt.Println("")
				fmt.Print("       ✉️  Enter message : ")
				message = readFullSentence()
				fmt.Println("")
				sendMessage(currentUser, to, message)
			}

		case 3:
			if currentUser == "" {

				var username, password string
				fmt.Print("       🔖 Enter your Username : ")
				fmt.Scan(&username)
				fmt.Println("                                 ")
				fmt.Print("       💀 Enter your Password : ")
				fmt.Scan(&password)
				fmt.Println("                                 ")

				if loginUser(username, password) {
					currentUser = username
				}

			}

			viewMessages(currentUser)

		case 4:
			currentUser = ""
			fmt.Println("                                 ")
			fmt.Println("                                 ")
			fmt.Println("       ✅ Logged Out ")

		case 5:
			fmt.Println("                                 ")
			fmt.Println("       Lets' Go 🚀")
			fmt.Println("                                 ")
			fmt.Println("                                 ")

			os.Exit(0)

		default:
			fmt.Println("                                 ")
			fmt.Println("                                 ")
			fmt.Print("       ❗ Invalid choice ", currentUser)

		}

		time.Sleep(1 * time.Second)
	}
}
