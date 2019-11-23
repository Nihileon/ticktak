package api

import "github.com/nihileon/ticktak/log"

var (
    adminList  []string
    SecretPass string
)

func Init(admins []string, secret string) error {
    adminList = admins
    log.GetLogger().Info("API init, adminList: %v", admins)
    SecretPass = secret
    log.GetLogger().Info("API init, SecretPass: %s", SecretPass)
    //initCleaner()
    return nil
}

func isAdmin(user string) bool {
    for _, admin := range adminList {
        if user == admin {
            return true
        }
    }
    return false
}
