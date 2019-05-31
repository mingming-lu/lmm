package command

// Register command
type Register struct {
	UserName     string
	EmailAddress string
	Password     string
}

// Login command
type Login struct {
	UserName string
	Password string
}

// AssignRole command
type AssignRole struct {
	OperatorUser string
	TargetUser   string
	TargetRole   string
}

// ChangePassword command
type ChangePassword struct {
	User        string
	OldPassword string
	NewPassword string
}
