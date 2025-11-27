package templates

const AlreadyLoggedInTemplate = `You are already logged in.
Do you want to login with a different account? (Y/n): `

const LoginSuccessTemplate = `
✓ Successfully logged in!
Welcome back, {{.Username}}!
`

const AuthenticatingTemplate = `
Authenticating...`
