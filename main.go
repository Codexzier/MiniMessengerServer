// ====================================================================
// Projekt:       mini messenger server
// Author:        Johannes P. Langner
// Description:   A simple service to send and get messages.
// ====================================================================

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

// UserItem = hold the user informaiton
type UserItem struct {
	ID       int64
	Username string
	IsOnline bool
}

// MessageItem = hold message and information about sender and receiver
type MessageItem struct {
	ID       int64
	Text     string
	UserID   int64
	ToUserID int64
	FromMe   bool
}

// DeviceItem = hold the device value setup
type DeviceItem struct {
	ID    int64
	Value int64
	Text  string
}

// ResponseError = Error message
type ResponseError struct {
	Success bool
	Content string
}

// ResponseGetUsers = result of request with an array of user
type ResponseGetUsers struct {
	Success bool
	Content []UserItem
}

// ResponseGetMessages = result of request with an array of message
type ResponseGetMessages struct {
	Success bool
	Content []MessageItem
}

// ResponseGetMessage = result of request with an array of message
type ResponseGetMessage struct {
	Success bool
	Content string
}

// ResponseAddUser = result of request of added user
type ResponseAddUser struct {
	Success bool
	Content UserItem
}

// ResponseSendMessages = result of request of added message
type ResponseSendMessages struct {
	Success bool
	Content MessageItem
	Value   int64
}

// ResponseDevice = only for flat structure
type ResponseDevice struct {
	Success bool
	ID      int64
	Content string
	Value   int64
	Text    string
}

// ResponseDevices = response all devices
type ResponseDevices struct {
	Success bool
	Content []DeviceItem
}

var userItems []UserItem = []UserItem{
	{ID: 1, Username: "Admin", IsOnline: false},
}

var messengerItems []MessageItem = []MessageItem{
	{ID: 1, Text: "Test Message", UserID: 1, ToUserID: 1},
}

var deviceItems []DeviceItem = []DeviceItem{
	{ID: 0, Value: 0},
}

func main() {

	var portNumber int = 5000

	fmt.Println("start webservice")
	fmt.Println(fmt.Sprintf("open the address %s:%v in your favorite browser", getHostAdress(), portNumber))

	startServer(portNumber)
}

func getHostAdress() string {
	networks, err := net.Interfaces()

	if err != nil {
		return "<error, can get host adress>"
	}

	var ipaddress string

	for _, network := range networks {

		addresses, err := network.Addrs()

		if err != nil {
			continue
		}

		if network.Flags == 0 {
			fmt.Println(fmt.Sprintf("Skip NETWORK: %s", network.Name))
			continue
		}

		fmt.Println(fmt.Sprintf("NETWORK Name: %s", network.Name))
		fmt.Println(fmt.Sprintf("NETWORK Flags: %s", network.Flags.String()))

		for _, address := range addresses {

			var ip net.IP

			switch netType := address.(type) {
			case *net.IPNet:

				if netType.IP.IsUnspecified() {
					fmt.Println("- Skip: Is unspecified")
					continue
				}

				if netType.IP.IsLoopback() {
					fmt.Println("- Skip: Is loopback")
					continue
				}

				ip = netType.IP
			case *net.IPAddr:
				ip = netType.IP
			}

			ipaddress = ip.String()
			fmt.Println(fmt.Sprintf("- IP address: %s", ip.String()))
		}
	}

	return ipaddress
}

// ====================================================================
// execute this method to start the webservice.
func startServer(port int) {

	http.HandleFunc("/", webserviceHandler)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

// ====================================================================
// handler to host the webside or request the data to json result
func webserviceHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Path //[1:]

	fmt.Println(command)

	// prevent the two time call
	// TODO: How can it do better?
	if command == "/favicon.ico" {
		return
	}

	if command == "/" {
		fmt.Fprintf(w, getWebsite())
		return
	}

	// get parameter id number from url request
	userID := r.URL.Query().Get("id")
	toUserID := r.URL.Query().Get("touserid")
	username := r.URL.Query().Get("username")
	messageText := r.URL.Query().Get("messagetext")
	valueStr := r.URL.Query().Get("value")
	textStr := r.URL.Query().Get("text")
	fmt.Println(fmt.Sprintf("received: %s %s %s %s %s %s", userID, toUserID, username, messageText, valueStr, textStr))

	jsonResult := getJSONnResult(command, userID, toUserID, username, messageText, valueStr, textStr)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	fmt.Fprintf(w, jsonResult)
}

// ====================================================================
// getWebsite = get website with information about possible requests
func getWebsite() string {
	webside, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	return string(webside)
}

// ====================================================================
// get webside login
// --------------------------------------------------------------------
// PARAMETERS
// command = parameter from the url request
func getJSONnResult(command string, idStr string, toUserID string, username string, messageText string, valueStr string, textStr string) string {

	if command == "" {
		return "error"
	}

	fmt.Println(command)

	var result []byte
	var err error

	switch command {
	case "/getAllUsers":
		{
			users := getOnlineUser(idStr)
			response := ResponseGetUsers{Success: true, Content: users}
			result, err = json.Marshal(&response)
			break
		}
	case "/getMessages":
		{
			messages := getMessages(idStr, toUserID)
			fmt.Println("- NORMAL Get MESSAGES")
			response := ResponseGetMessages{Success: true, Content: messages}
			result, err = json.Marshal(&response)
			break
		}
	case "/addUser":
		{
			userItem := addUser(username)
			response := ResponseGetUsers{Success: true, Content: []UserItem{userItem}}
			result, err = json.Marshal(&response)
			break
		}
	case "/sendMessage":
		{
			message := sendMessage(idStr, toUserID, messageText)
			response := ResponseSendMessages{Success: true, Content: message}
			result, err = json.Marshal(&response)
			break
		}
	case "/deviceGetAll":
		{
			fmt.Println(fmt.Sprintf("Devices %v", len(deviceItems)))
			response := ResponseDevices{Success: true, Content: deviceItems}
			result, err = json.Marshal(&response)
			break
		}
	case "/deviceSendCommand":
		{
			message, id, valueRe := deviceSendCommand(idStr, valueStr, textStr)
			response := ResponseDevice{Success: true, ID: id, Content: message, Value: valueRe, Text: textStr}
			result, err = json.Marshal(&response)
			break
		}
	case "/deviceGetValue":
		{
			message, id, value := deviceGetValue(idStr)
			response := ResponseDevice{Success: true, ID: id, Content: message, Value: value}
			result, err = json.Marshal(&response)
			break
		}
	case "/deviceGetText":
		{
			message, id, value := deviceGetText(idStr)
			response := ResponseDevice{Success: true, ID: id, Content: message, Text: value}
			result, err = json.Marshal(&response)
			break
		}
	case "/deviceGet":
		{
			message, id, value, text := deviceGet(idStr)
			response := ResponseDevice{Success: true, ID: id, Content: message, Value: value, Text: text}
			result, err = json.Marshal(&response)
			break
		}
	default:
		{
			errorMessage := fmt.Sprintf("no case for this command: %s", command)
			fmt.Println(errorMessage)
			response := ResponseError{Success: false, Content: errorMessage}
			result, err = json.Marshal(&response)
			break
		}
	}

	if err != nil {
		fmt.Println(err)
		return "err"
	}
	return string(result)
}

// deviceSendCommand = set the value for device
func deviceSendCommand(deviceIDStr string, valueStr string, textStr string) (string, int64, int64) {

	deviceID := parseValidNumber(deviceIDStr)
	value := parseValidNumber(valueStr)

	for index := 0; index < len(deviceItems); index++ {
		if deviceID == deviceItems[index].ID {
			deviceItems[index].Value = value
			deviceItems[index].Text = textStr
			return "Device found", deviceID, value
		}
	}

	return "no device", deviceID, value
}

func deviceGetValue(idStr string) (string, int64, int64) {

	deviceID := parseValidNumber(idStr)

	for index := 0; index < len(deviceItems); index++ {
		if deviceID == deviceItems[index].ID {
			return "OK", deviceID, deviceItems[index].Value
		}
	}

	var deviceItem = DeviceItem{ID: deviceID, Value: 0}
	deviceItems = append(deviceItems, deviceItem)

	return "missing", deviceID, 0
}

func deviceGetText(idStr string) (string, int64, string) {

	deviceID := parseValidNumber(idStr)

	for index := 0; index < len(deviceItems); index++ {
		if deviceID == deviceItems[index].ID {
			return "OK", deviceID, deviceItems[index].Text
		}
	}

	var deviceItem = DeviceItem{ID: deviceID, Value: 0}
	deviceItems = append(deviceItems, deviceItem)

	return "missing", deviceID, "--"
}

func deviceGet(idStr string) (string, int64, int64, string) {

	deviceID := parseValidNumber(idStr)

	for index := 0; index < len(deviceItems); index++ {
		if deviceID == deviceItems[index].ID {
			return "OK", deviceID, deviceItems[index].Value, deviceItems[index].Text
		}
	}

	var deviceItem = DeviceItem{ID: deviceID, Value: 0}
	deviceItems = append(deviceItems, deviceItem)

	return "missing", deviceID, 0, "--"
}

func parseValidNumber(numberStr string) int64 {

	if numberStr == "" {
		return 0
	}

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return number
}

// getOnlineUser = get all user with status and without self
func getOnlineUser(userIDStr string) []UserItem {
	users := []UserItem{}

	userID := parseValidNumber(userIDStr)

	for index := 0; index < len(userItems); index++ {
		if userItems[index].ID != userID {
			users = append(users, userItems[index])
		}
	}

	return users
}

// getMessages = get all receive message from other user
func getMessages(userIDStr string, toUserIDStr string) []MessageItem {

	temp := []MessageItem{}

	userID := parseValidNumber(userIDStr)
	toUserID := parseValidNumber(toUserIDStr)

	// get all message for this calling user
	for index := 0; index < len(messengerItems); index++ {
		if messengerItems[index].ToUserID == userID {
			temp = append(temp, messengerItems[index])
		}
	}

	resultItems := []MessageItem{}

	// get all message
	for i := 0; i < len(temp); i++ {
		if temp[i].UserID == toUserID {
			resultItems = append(resultItems, temp[i])
		}
	}

	// TODO: special condition
	if toUserID == 1 {

		if len(resultItems) == 0 {
			sendMessage(toUserIDStr, userIDStr, "Hallo ich bin der Admin Benutzer")

			return getMessages(userIDStr, toUserIDStr)
		}

	}

	// TODO: special condition
	// if wemos1, return only the last message
	if userID == 3 {
		if len(resultItems) == 0 {
			return []MessageItem{{ID: 0, Text: "NO DATA"}}
		}

		return []MessageItem{resultItems[len(resultItems)-1]}
	}

	return resultItems
}

// addUser = add new user and return that user data
func addUser(username string) UserItem {

	// check exist user and exit, if exist
	for _, item := range userItems {
		if item.Username == username {
			return item
		}
	}

	var user = UserItem{ID: createUserID(), Username: username, IsOnline: false}
	userItems = append(userItems, user)

	return user
}

// createUserID get the next user id
func createUserID() int64 {
	var maxID int64

	for index := 0; index < len(userItems); index++ {
		if userItems[index].ID > maxID {
			maxID = userItems[index].ID
		}
	}

	return maxID + 1
}

// sendMessage add new message
func sendMessage(userIDStr string, toUserIDStr string, messageText string) MessageItem {
	userID := parseValidNumber(userIDStr)
	toUserID := parseValidNumber(toUserIDStr)

	// message for self
	messageItemSelf := MessageItem{ID: createMessageID(), UserID: toUserID, ToUserID: userID, Text: messageText, FromMe: true}
	messengerItems = append(messengerItems, messageItemSelf)

	messageItem := MessageItem{ID: createMessageID(), UserID: userID, ToUserID: toUserID, Text: messageText, FromMe: false}
	messengerItems = append(messengerItems, messageItem)

	return messageItem
}

func createMessageID() int64 {
	var maxID int64

	for index := 0; index < len(messengerItems); index++ {
		if messengerItems[index].ID > maxID {
			maxID = messengerItems[index].ID
		}
	}

	return maxID + 1
}
