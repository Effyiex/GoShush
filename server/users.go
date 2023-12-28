package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const USERS_STORAGE = "users.json"

type User struct {
	UUID        uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Session     uuid.UUID `json:"session"`
	Permissions []string  `json:"permissions"`
}

type ClientUser struct {
	UUID        uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Session     uuid.UUID `json:"session"`
	Permissions []string  `json:"permissions"`
}

func (user *User) ToClient() ClientUser {
	return ClientUser{
		UUID:        user.UUID,
		Name:        user.Name,
		Session:     user.Session,
		Permissions: user.Permissions,
	}
}

type UserManager struct {
	Users        []User   `json:"users"`
	RegisterKeys []string `json:"register_keys"`
	Server       *Server  `json:"-"`
}

func BootUserManager(dataFolder string) UserManager {
	um := UserManager{}
	um.Load(dataFolder)
	return um
}

func (um *UserManager) GenerateRegister(key string) {
	um.RegisterKeys = append(um.RegisterKeys, key)
}

func (um *UserManager) Load(dataFolder string) {
	data, err := os.ReadFile(dataFolder + "/" + USERS_STORAGE)
	if err != nil {
		um.Save(dataFolder)
		return
	}
	err = json.Unmarshal(data, um)
	if err != nil {
		panic(err)
	}
}

func (um *UserManager) Save(dataFolder string) {
	data, _ := json.MarshalIndent(um, "", "  ")
	os.WriteFile(dataFolder+"/"+USERS_STORAGE, data, 0644)
}

func (um *UserManager) StripRegisterKey(key string) {
	for i, k := range um.RegisterKeys {
		if k == key {
			um.RegisterKeys = append(um.RegisterKeys[:i], um.RegisterKeys[i+1:]...)
			return
		}
	}
}

func (um *UserManager) Login(username string, password string) *User {
	for i, user := range um.Users {
		if user.Name == username && user.Password == password {
			um.Users[i].Session = uuid.New()
			return &um.Users[i]
		}
	}
	return nil
}

func (um *UserManager) Register(username string, password string) bool {
	if um.Login(username, password) != nil {
		return false
	}
	um.Users = append(um.Users, User{
		UUID:        uuid.New(),
		Name:        username,
		Password:    password,
		Session:     uuid.New(),
		Permissions: []string{},
	})
	return true
}

func (um *UserManager) RegisterOn(srv Server) {
	srv.HandleFunc("/api/auth", um.RunAuthentification()).Methods("POST")
	srv.HandleFunc("/api/rename", um.Rename()).Methods("POST")
}

type AuthPacket struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	RegisterKey string `json:"key"`
}

func (um *UserManager) RunAuthentification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var packet AuthPacket
		err := json.NewDecoder(r.Body).Decode(&packet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if packet.RegisterKey != "" {
			for _, key := range um.RegisterKeys {
				if key == packet.RegisterKey {
					if um.Register(packet.Username, packet.Password) {
						um.StripRegisterKey(key)
					}
					break
				}
			}
		}

		user := um.Login(packet.Username, packet.Password)
		if user != nil {
			json.NewEncoder(w).Encode(user.ToClient())
		} else {
			w.Write([]byte("null"))
		}

	}
}

type RenamePacket struct {
	User    uuid.UUID `json:"user"`
	Session uuid.UUID `json:"session"`
	NewName string    `json:"new_name"`
}

func (um *UserManager) Rename() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var packet RenamePacket
		err := json.NewDecoder(r.Body).Decode(&packet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if um.IsValidSession(packet.User, packet.Session) {
			var userIndex = -1
			for i, u := range um.Users {
				if u.UUID == packet.User {
					userIndex = i
				} else if u.Name == packet.NewName {
					w.Write([]byte("name already taken"))
					return
				}
			}
			um.Users[userIndex].Name = packet.NewName
			w.Write([]byte("success"))
		} else {
			w.Write([]byte("invalid session"))
		}

	}
}

func (um *UserManager) IsValidSession(user uuid.UUID, session uuid.UUID) bool {
	for _, u := range um.Users {
		if u.UUID == user && u.Session == session {
			return true
		}
	}
	return false
}

func (um *UserManager) GrantPermission(user uuid.UUID, permission string) bool {
	for i, u := range um.Users {
		if u.UUID == user {
			um.Users[i].Permissions = append(um.Users[i].Permissions, permission)
			return true
		}
	}
	return false
}

func (um *UserManager) RevokePermission(user uuid.UUID, permission string) bool {
	for i, u := range um.Users {
		if u.UUID == user {
			for j, p := range um.Users[i].Permissions {
				if p == permission {
					um.Users[i].Permissions = append(um.Users[i].Permissions[:j], um.Users[i].Permissions[j+1:]...)
					return true
				}
			}
		}
	}
	return false
}

func (um *UserManager) HasPermission(user uuid.UUID, permission string) bool {
	for _, u := range um.Users {
		if u.UUID == user {
			for _, p := range u.Permissions {
				if p == permission {
					return true
				}
			}
		}
	}
	return false
}

func (um *UserManager) GetUser(user uuid.UUID) *User {
	for _, u := range um.Users {
		if u.UUID == user {
			return &u
		}
	}
	return nil
}
