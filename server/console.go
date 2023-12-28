package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func (srv *Server) CommandLineListener() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		args := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		if len(args) == 0 {
			fmt.Println("Insufficient Arguments.")
			continue
		}
	SwitchA0:
		switch args[0] {

		case "help", "/help", "/?", "?", "-help":
			fmt.Println(" > HelpPage")
			fmt.Println("1.: \"users ...\"")
			fmt.Println("2.: \"walkwebroot\"")
			fmt.Println("3.: \"chats ...\"")
			fmt.Println("4.: \"<stop / exit / quit>\"")
			break SwitchA0

		case "walkwebroot":
			fmt.Println("Walking web root...")
			srv.WalkWebFilesRoot()
			fmt.Println("Finished walk!")
			break SwitchA0

		case "stop", "exit", "quit":
			srv.Dispose(0)
			break SwitchA0

		case "chats":
			if len(args) == 1 {
				fmt.Println("Insufficient Arguments.")
				break SwitchA0
			}
		SwitchA1_a:
			switch args[1] {

			case "help", "/help", "/?", "?", "-help":
				fmt.Println(" > HelpPage")
				fmt.Println("1.: \"" + args[0] + " save\"")
				break SwitchA1_a

			case "save":
				fmt.Println("Saving Chats...")
				srv.Chats.Save(srv.Config.DataFolder)
				fmt.Println("Chats saved.")
				break SwitchA1_a

			default:
				fmt.Println("Unknown command: \"" + strings.Join(args[0:1], "\"."))

			}
			break SwitchA0

		case "users":
			if len(args) == 1 {
				fmt.Println("Insufficient Arguments.")
				break SwitchA0
			}
		SwitchA1:
			switch args[1] {

			case "help", "/help", "/?", "?", "-help":
				fmt.Println(" > HelpPage")
				fmt.Println("1.: \"" + args[0] + " by-name <username>\"")
				fmt.Println("2.: \"" + args[0] + " gen-register <key>\"")
				fmt.Println("3.: \"" + args[0] + " permissions ...\"")
				fmt.Println("4.: \"" + args[0] + " save\"")
				break

			case "by-name":
				if len(args) == 2 {
					fmt.Println("Insufficient Arguments.")
					break SwitchA1
				}
				for _, user := range srv.Users.Users {
					if user.Name == args[2] {
						client_user_buffer, _ := json.MarshalIndent(user.ToClient(), "", "  ")
						fmt.Println(string(client_user_buffer))
						break SwitchA1
					}
				}
				fmt.Println("User \"" + args[2] + "\" not found.")
				break SwitchA1

			case "gen-register":
				if len(args) == 2 {
					fmt.Println("Insufficient Arguments.")
					break SwitchA1
				}
				srv.Users.GenerateRegister(args[2])
				fmt.Println("New Register-Key: \"" + args[2] + "\".")
				break SwitchA1

			case "save":
				fmt.Println("Saving users...")
				srv.Users.Save(srv.Config.DataFolder)
				fmt.Println("Users saved.")
				break SwitchA1

			case "permissions":
				if len(args) == 2 {
					fmt.Println("Insufficient Arguments.")
					break SwitchA1
				}
			SwitchA2:
				switch args[2] {

				case "?":
				case "help":
				case "/help":
				case "-help":
					fmt.Println("Commands of \"" + args[1] + "\".")
					fmt.Println("... " + args[1] + " grant <permission> <uuid>")
					fmt.Println("... " + args[1] + " revoke <permission> <uuid>")
					break SwitchA2

				case "grant":
					uuid, err := uuid.Parse(args[3])
					if err != nil {
						fmt.Println("Invalid UUID \"" + args[3] + "\".")
						break SwitchA2
					}
					if srv.Users.GrantPermission(uuid, args[4]) {
						fmt.Println("Permission \"" + args[4] + "\" granted to \"" + args[3] + "\".")
					} else {
						fmt.Println("User \"" + args[3] + "\" not found.")
					}
					break SwitchA2

				case "revoke":
					uuid, err := uuid.Parse(args[3])
					if err != nil {
						fmt.Println("Invalid UUID \"" + args[3] + "\".")
						break SwitchA2
					}
					if srv.Users.RevokePermission(uuid, args[4]) {
						fmt.Println("Permission \"" + args[4] + "\" revoked from \"" + args[3] + "\".")
					} else {
						fmt.Println("User \"" + args[3] + "\" not found.")
					}
					break SwitchA2

				default:
					fmt.Println("Unknown command: \"" + strings.Join(args[0:2], "\"."))

				}
				break SwitchA1

			default:
				fmt.Println("Unknown command: \"" + strings.Join(args[0:1], "\"."))

			}
			break SwitchA0

		default:
			fmt.Println("Unknown command: \""+args[0], "\".")

		}
	}
}
