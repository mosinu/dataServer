package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

//DefaultNamespace defalut namespace
var DefaultNamespace Namespace

//NamespaceType type of namespace
type NamespaceType uint8

//Namespace types
const (
	UserNamespaceType NamespaceType = iota
	CustomNamespaceType
)

//Namespace a namespace for files
type Namespace struct {
	gorm.Model
	Name   string `gorm:"not null"`
	UserID uint   `gorm:"column:creator;index"`
	User   *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
}

//GetNamespaceFromString return namespace from string
func GetNamespaceFromString(ns string) *Namespace {
	if len(ns) == 0 || strings.ToLower(ns) == "default" {
		return &DefaultNamespace
	}

	return &Namespace{
		Name: ns,
	}
}

//FindNamespace find namespace in DB
func FindNamespace(db *gorm.DB, ns string, user *User) *Namespace {
	namespace := GetNamespaceFromString(ns)
	if namespace.ID != 0 {
		return namespace
	}

	db.Where(&namespace).Preload("User").Find(&namespace)

	if (namespace == nil || namespace.ID == 0) && !strings.HasPrefix(ns, user.Username+"_") {
		namespace.Name = UserNamespaceName(namespace.Name, user)
		db.Where(&namespace).Preload("User").Find(&namespace)
	}

	return namespace
}

//IsOwnedBy returns true if namespace is users
func (namespace *Namespace) IsOwnedBy(user *User) bool {
	if user == nil || namespace == nil {
		return false
	}
	return namespace.UserID == user.ID
}

//FindUserNamespaces get all namespaces for user
func FindUserNamespaces(db *gorm.DB, user *User) ([]Namespace, error) {
	var namespaces []Namespace

	err := db.Model(&Namespace{}).Where("creator = ?", user.ID).Find(&namespaces).Error
	if err != nil {
		return []Namespace{}, err
	}

	return namespaces, nil
}

//UserNamespaceName get name of usernamespace
func UserNamespaceName(namespaceName string, user *User) string {
	if strings.HasPrefix(namespaceName, user.Username+"_") {
		return namespaceName
	}
	return user.Username + "_" + namespaceName
}

//IsValid return true if namespace is valid
func (namespace *Namespace) IsValid() bool {
	return (namespace != nil && namespace.ID > 0)
}
