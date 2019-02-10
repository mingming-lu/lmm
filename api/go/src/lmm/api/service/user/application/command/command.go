package command

// AssignRole command
type AssignRole struct {
	OperatorUser string
	TargetUser   string
	TargetRole   string
}

// ChangePassword command
type ChangePassword struct {
	User        string
	NewPassword string
}
