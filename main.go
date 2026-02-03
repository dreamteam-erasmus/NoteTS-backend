package main

import (
	"noteTS-backend/tools"
	"webtools"
	"webtools/httptools"
)

func main() {
	tools.INITLogger()
	tools.LoadUsersFromJSON()
	tools.SaveUsersToJSON()
	sv := httptools.NewWebSocketServer("0.0.0.0:8080", readFunc, nil, "../noteTS/src", false, true, true)
	sv.Start()
}

var tokensToUsers webtools.SafeMap[string, *tools.User] = webtools.MakeSafeMap[string, *tools.User]()

func checkValidInstance(params map[string]string) bool {
	user := tokensToUsers.Get(params["token"])
	if user == nil {
		return false
	}
	return user.Login == params["login"]
}

func readFunc(conn *httptools.WebSocketServerConn, data []byte, status uint8, isBinary bool) {
	_ = isBinary
	if status != webtools.ReadDataStatus {
		return
	}

	//Read data of websocket
	command, params := httptools.CreateParametersFromURL(string(data))
	switch command {
	case "login":
		{
			//Try to login the user
			foundUser := false
			tools.Logger.Log(1, "Trying to login: "+params["login"])
			for _, v := range tools.Users {
				if v.Login == params["login"] {
					foundUser = true
					isOk, err := v.Password.CheckPassword(params["password"])
					if err != nil {
						//Error checking password
						tools.Logger.Log(3, "Error checking password: "+err.Error())
						conn.Send([]byte(httptools.CreateURLFromParameters("error", nil)))
						return
					}
					if !isOk {
						//Error checking password
						tools.Logger.Log(3, "Invalid password for user: "+params["password"])
						conn.Send([]byte(httptools.CreateURLFromParameters("error", nil)))
						return
					}

					//Password OK, send auth token
					tools.Logger.Log(1, "User logged in: "+params["login"])
					token := webtools.GenerateRandomID()
					tokensToUsers.Set(token, &v)
					conn.Send([]byte(httptools.CreateURLFromParameters("login", map[string]string{"token": token})))
					break
				}
			}

			//User not found
			if !foundUser {
				tools.Logger.Log(3, "Invalid user: "+params["login"])
				conn.Send([]byte(httptools.CreateURLFromParameters("error", nil)))
				return
			}
			break
		}
	case "checkLogin":
		{
			//Cbeck if token is valid
			tools.Logger.Log(0, "Checking instance for: "+params["login"])
			conn.Send([]byte(httptools.CreateURLFromParameters(webtools.FormatByBool(checkValidInstance(params), "instanceOK", "invalidInstance"), nil)))
			break
		}
	}
}
