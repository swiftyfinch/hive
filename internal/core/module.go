package core

const Modules_File_Name = "modules.yml"

type Module struct {
	Name         string
	Dependencies []string
}
